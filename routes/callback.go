package routes

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func HandleCallback(c http.ResponseWriter, req *http.Request) {
	chall := req.FormValue("challenge")
	if len(chall) != 0 {
		c.Header().Add("Content-Length", strconv.Itoa(len(chall)))
		io.WriteString(c, chall)
	} else {
		con, err := ioutil.ReadAll(req.Body)
		if err == nil {
			log.Print(string(con))
		}
	}
}
