package NamiCommands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateImplant(arch int, implantName string) string {
	fmt.Println("\n\t[*] Building implant...")
	var goArch string
	if arch == 32 {
		goArch = "386"
	} else if arch == 64 {
		goArch = "amd64"
	}
	err := os.Chdir("implant")
	if err != nil {
		return "\t[-] Implant folder does not exist!\n"
	}

	buildVar := fmt.Sprintf("-tags='%s'", goArch)
	ldflagsVar := "-w -s -H=windowsgui"
	outputVar := filepath.FromSlash(fmt.Sprintf("../output/%s.exe", implantName))
	os.Setenv("GOARCH", goArch)
	os.Setenv("GOOS", "windows")

	cmd := exec.Command("go", "build", "-ldflags", ldflagsVar, buildVar, "-o", outputVar, "implant.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		os.Chdir("..")
		return fmt.Sprintf("\t[-] Error building implant! Error: %s\n", &stderr)
	} else {
		os.Chdir("..")
		return "\t[+] Successful payload build!\n"
	}
}
