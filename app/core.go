package app

import (
	"github.com/callumj/auto-convert/routes"
	"github.com/callumj/auto-convert/shared"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func Run(args []string) {
	if len(args) >= 1 {
		shared.LoadConfig(args[1])
	} else {
		os.Exit(-1)
	}

	shared.InitDb()

	listenOn := shared.Config.Listen
	if len(listenOn) == 0 {
		listenOn = ":8080"
	}

	log.Printf("Starting web server on %v", listenOn)

	r := mux.NewRouter()

	r.HandleFunc("/begin", routes.BeginAuthHandler)
	r.HandleFunc("/complete_auth", routes.CompleteAuthHandler)

	http.Handle("/", r)
	http.ListenAndServe(listenOn, nil)
}
