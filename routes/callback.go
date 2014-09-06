package routes

import (
	"encoding/json"
	"github.com/callumj/auto-convert/workers"
	"io"
	"log"
	"net/http"
	"strconv"
)

type DeltaUserMessage struct {
	Users []int64 `json:"users"`
}

type DeltaMessage struct {
	Delta DeltaUserMessage `json:"delta"`
}

func HandleCallback(c http.ResponseWriter, req *http.Request) {
	chall := req.FormValue("challenge")
	if len(chall) != 0 {
		log.Printf("Incoming challenge %s\n", chall)
		c.Header().Add("Content-Length", strconv.Itoa(len(chall)))
		io.WriteString(c, chall)
	} else {
		dec := json.NewDecoder(req.Body)
		var decoded DeltaMessage
		err := dec.Decode(&decoded)
		if err != nil {
			log.Print(err)
			return
		} else {
			for _, item := range decoded.Delta.Users {
				workers.DispatchDelta(workers.DeltaRequest{Uid: item})
			}
		}
	}
}
