package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// RequestBody is used to store the slack message
type RequestBody struct {
	Text string `json:"text"`
}

// NewMessage is used to create a message body for sending to slack
func NewMessage(message string) *RequestBody {
	return &RequestBody{
		message,
	}
}

// SendMessage sends a message to slack
func SendMessage(request *RequestBody, webhook string) error {
	slackBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
