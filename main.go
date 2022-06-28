package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/minderjan/opentransport-client/opentransport"
	"skittle.ch/test/opentransport/train"
)

type MockClient struct {
}

func (m MockClient) Search(context.Context, string, string, time.Time) (*opentransport.ConnectionResult, error) {
	return &opentransport.ConnectionResult{Connections: []opentransport.Connection{
		{
			From: opentransport.Stop{Station: opentransport.Location{Name: "Zürich HB"}},
		},
	}}, nil
}

func main() {
	server := train.TrainServer{
		Client: opentransport.NewClient().Connection,
	}
	fs := http.FileServer(http.Dir("./styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

	http.HandleFunc("/connection", server.ServeConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8090", nil)
	err := http.ListenAndServe(":8090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
