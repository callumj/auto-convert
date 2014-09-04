package lib

import (
	"encoding/json"
	"fmt"
	"github.com/callumj/auto-convert/shared"
	"log"
	"net/http"
	"strings"
)

type dropboxDeltaResponse struct {
	HasMore bool            `json:"has_more"`
	Cursor  string          `json:"cursor"`
	Entries [][]interface{} `json:"entries"`
	Reset   bool            `json:"reset"`
}

func GetChangedFiles(acc shared.Account) {
	req, err := http.NewRequest("POST", "https://api.dropbox.com/1/delta", nil)
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

	fetched := processRequest(req)
	if fetched == nil {
		return
	}
	processFetched(fetched)
	for fetched.HasMore {
		log.Printf("Fetched %d\n", len(fetched.Entries))
		req, err = http.NewRequest("POST", fmt.Sprintf("https://api.dropbox.com/1/delta?cursor=%v", fetched.Cursor), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

		log.Printf("Requesting %v\n", req.URL.String())
		if err != nil {
			log.Print(err)
			return
		}

		next := processRequest(req)
		if next != nil {
			fetched = next
			log.Printf("Fetched %d\n", len(fetched.Entries))
			processFetched(fetched)
		} else {
			break
		}
	}
}

func processFetched(fetched *dropboxDeltaResponse) {
	for _, item := range fetched.Entries {
		switch v := item[0].(type) {
		case string:
			// v is a string here, so e.g. v + " Yeah!" is possible.
			for src, _ := range shared.Config.Maps {
				if strings.Contains(v, src) {
					log.Println(v)
				}
			}
		}
	}
}

func processRequest(req *http.Request) *dropboxDeltaResponse {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	} else {
		defer resp.Body.Close()
		if err == nil {
			dec := json.NewDecoder(resp.Body)
			var decoded dropboxDeltaResponse
			err = dec.Decode(&decoded)
			if err != nil {
				log.Print(err)
				return nil
			}

			return &decoded
		} else {
			log.Printf("Error: %v\n", err)
			return nil
		}
	}
}
