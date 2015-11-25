// transmission-notify sends a pushbullet notification to my phone when a
// torrent download is completed.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	client = http.Client{Timeout: 10 * time.Second}
)

const (
	// pushbullet
	PBBaseURL     = "https://api.pushbullet.com/v2"
	PBAccessToken = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	PBDeviceID    = ""
)

func main() {
	torrent := os.Getenv("TR_TORRENT_NAME")
	err := push(torrent)
	if err != nil {
		log.Printf("Error while pushing a notification: %v\n", err)
		return
	}
}

func push(txt string) error {
	body, err := newNote(txt)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", PBBaseURL+"/pushes", body)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", PBAccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to push the notification. HTTP Status: %v\n", resp.StatusCode)
	}
	return nil
}

type note struct {
	DeviceID string `json:"device_iden,omitempty"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

func newNote(body string) (io.Reader, error) {
	n := note{
		DeviceID: PBDeviceID,
		Type:     "note",
		Title:    "A torrent is complete!",
		Body:     body,
	}
	b, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
