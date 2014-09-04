package lib

import (
	"fmt"
	"github.com/callumj/auto-convert/shared"
	"io/ioutil"
	"log"
	"net/http"
)

func GetChangedFiles(acc shared.Account) {
	req, err := http.NewRequest("POST", "https://api.dropbox.com/1/delta", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		} else {
			log.Print(string(body))
		}
	}
}
