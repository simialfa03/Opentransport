package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/minderjan/opentransport-client/opentransport"
	"skittle.ch/test/opentransport/train"
)

//go:embed static
var static embed.FS

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

	fsys := fs.FS(static)
	html, _ := fs.Sub(fsys, "static")

	server := train.TrainServer{
		Client: opentransport.NewClient().Connection,
		FS:     html,
	}

	fs := http.FileServer(http.FS(html))
	http.Handle("/", fs)
	http.HandleFunc("/connection", server.ServeConnection)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "index.html")
	// })

	err := http.ListenAndServe(":8090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
