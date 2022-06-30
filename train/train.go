package train

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/minderjan/opentransport-client/opentransport"
)

// Create a new Client
// var Client = opentransport.NewClient().Connection

type OpentransportClient interface {
	Search(context.Context, string, string, time.Time) (*opentransport.ConnectionResult, error)
}

type TrainServer struct {
	Client OpentransportClient
	FS     fs.FS
}

func (t TrainServer) Get(origin, destination string) ([]opentransport.Connection, error) {
	// Query a connection
	connResult, err := t.Client.Search(context.Background(), origin, destination, time.Now())
	return connResult.Connections, err
}

func (t TrainServer) ServeConnection(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	origin := req.FormValue("origin")
	destination := req.FormValue("destination")

	fmt.Print(origin, destination)

	tmpl, err := template.New("connection.html").Funcs(
		template.FuncMap{
			"FormatTime": func(t time.Time) string {
				return t.Format("15:04")
			},
		}).ParseFS(t.FS, "*.html")
	if err != nil {
		http.Error(w, "Unable to Parse Template", http.StatusInternalServerError)
		return
	}

	type Connections struct {
		Items []opentransport.Connection
	}

	items, err := t.Get(origin, destination)
	if err != nil {
		http.Error(w, "Unable to Get Connection", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, Connections{Items: items}); err != nil {
		http.Error(w, "Unable to Execute Template"+err.Error(), http.StatusInternalServerError)
		return
	}
}
