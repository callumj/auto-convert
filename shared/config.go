package shared

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
)

var Config Configuration

type Configuration struct {
	DropboxKey    string
	DropboxSecret string
	Maps          map[string]string
	Listen        string
	DbFile        string
}

func LoadConfig(path string) bool {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to find configuration %v (%v)\r\n", path, err)
		return false
	}

	err = yaml.Unmarshal([]byte(fileContents), &Config)
	if err != nil {
		log.Printf("Unable to parse configuration %v\r\n", path)
		return false
	}

	verifyRequirements()

	return true
}

func verifyRequirements() {
	if len(Config.DropboxKey) == 0 {
		log.Fatal("dropbox_key is not set!")
	}

	if len(Config.DropboxSecret) == 0 {
		log.Fatal("dropbox_secret is not set!")
	}
}
