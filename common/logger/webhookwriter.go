package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CtxKey string

func PrettyJSON(val interface{}) (string, error) {
	v, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func PrettyPrint(val interface{}) {
	v, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		fmt.Println("expected valid marshable object", err)
	}
	fmt.Println(string(v))
}

type webhookWriter struct {
	Url     string
	client  *http.Client
	service string
}

func NewWebhookWriter(url string, service string) *webhookWriter {
	return &webhookWriter{
		Url:     url,
		client:  http.DefaultClient,
		service: service,
	}
}

// Write implements the io.Writer interface
func (w *webhookWriter) Write(p []byte) (n int, err error) {
	// Making a copy of the incoming bytes to avoid data race in the goroutine
	dataCopy := make([]byte, len(p))
	copy(dataCopy, p)

	go func() {
		type loggerMsg struct {
			Time    float64 `json:"ts"`
			Message string  `json:"msg"`
		}
		var msg loggerMsg
		err = json.Unmarshal(dataCopy, &msg)
		if err != nil {
			fmt.Printf("error unmarshalling JSON: %v\n", err)
			return
		}

		readableTime := time.Unix(int64(msg.Time), 0).UTC().Format(time.RFC3339)

		data := map[string]interface{}{
			"title": msg.Message,
			// raw message to get the fields too
			"message":   json.RawMessage(dataCopy),
			"timestamp": readableTime,
			"service":   "COBI/v2",
			"severity":  "critical",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("error marshaling JSON: %v\n", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create a new HTTP request
		req, err := http.NewRequestWithContext(ctx, "POST", w.Url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating webhook request: %v\n", err)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")

		resp, err := w.client.Do(req)
		if err != nil {
			fmt.Printf("Error requesting the webhook url: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Webhook request failed with status code %d\n", resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Response body: %s\n", string(body))
		}
	}()

	// Return the byte count and no error as if the write succeeded
	return len(dataCopy), nil
}
