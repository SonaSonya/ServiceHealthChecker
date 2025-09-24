package healthchecker

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
	StatusErr  Status = "ERROR"
)

const (
	PingEndpoint      = "/ping"
	LastEventEndpoint = "/last-event-time"
)

var httpClient = http.Client{Timeout: 2 * time.Second}

type PingResult struct {
	URL    string
	Status Status
	Event  string
	Err    error
}

func Ping(url string, ch chan<- PingResult) {

	ok, err := getPing(url)
	if err != nil {
		ch <- PingResult{URL: url, Status: StatusErr, Event: "", Err: fmt.Errorf("ping error: %w", err)}
		return
	}

	status := StatusDown
	if ok {
		status = StatusUp
	}

	event, err := getLastEventTime(url)
	if err != nil {
		ch <- PingResult{URL: url, Status: status, Event: "", Err: fmt.Errorf("last event error: %w", err)}
		return
	}

	ch <- PingResult{URL: url, Status: status, Event: event, Err: nil}
}

func getPing(url string) (bool, error) {
	pingURL := url + PingEndpoint
	resp, err := httpClient.Get(pingURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("returned status code %d", resp.StatusCode)
	}

	return true, nil
}

func getLastEventTime(url string) (string, error) {
	eventUrl := url + LastEventEndpoint
	resp, err := httpClient.Get(eventUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("returned status code %d", resp.StatusCode)
	}

	return string(body), nil
}
