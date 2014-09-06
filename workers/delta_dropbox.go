package workers

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

func GetChangedFiles(acc *shared.Account) {
	var url string
	if len(acc.LastCursor) == 0 {
		url = "https://api.dropbox.com/1/delta"
	} else {
		url = fmt.Sprintf("https://api.dropbox.com/1/delta?cursor=%v", acc.LastCursor)
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

	fetched := processRequest(req)
	if fetched == nil {
		return
	} else {
		commitCursor(fetched.Cursor, acc)
	}
	processFetched(fetched, acc)
	for fetched.HasMore {
		log.Printf("Fetched %d\n", len(fetched.Entries))
		req, err = http.NewRequest("POST", fmt.Sprintf("https://api.dropbox.com/1/delta?cursor=%v", fetched.Cursor), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

		if err != nil {
			log.Print(err)
			return
		}

		next := processRequest(req)
		if next != nil {
			commitCursor(next.Cursor, acc)
			fetched = next
			log.Printf("Fetched %d\n", len(fetched.Entries))
			processFetched(fetched, acc)
		} else {
			break
		}
	}
}

func processFetched(fetched *dropboxDeltaResponse, acc *shared.Account) {
	for _, item := range fetched.Entries {
		switch v := item[0].(type) {
		case string:
			for src, _ := range shared.Config.Maps {
				if strings.HasPrefix(v, src) {
					DispatchFile(FileRequest{Uid: acc.Uid, Path: v, MatchedPath: src})
				}
			}
		}
	}
}

func commitCursor(cursor string, acc *shared.Account) {
	acc.LastCursor = cursor
	shared.UpdateAccount(acc)
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
