package main

import (
	"testing"

	"skittle.ch/test/opentransport/train"
)

func TestGetTrain(t *testing.T) {
	server := train.TrainServer{
		Client: MockClient{},
	}
	conns, err := server.Get("Zürich HB", "Wil")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", conns)
}

//type MockClient struct {
//}

//func (m MockClient) Search(context.Context, string, string, time.Time) (*opentransport.ConnectionResult, error) {
//	return &opentransport.ConnectionResult{Connections: []opentransport.Connection{
//		{
//			From: opentransport.Stop{Station: opentransport.Location{Name: "Zürich HB"}},
//		},
//	}}, nil
//}
