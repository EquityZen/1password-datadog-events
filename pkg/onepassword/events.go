package onepassword

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"time"
)

var Token = viper.GetString("CONNECT_TOKEN")
var Version string
var DefaultUserAgent = fmt.Sprintf("1Password Events API for DataDog / %s", Version)

type EventsAPI struct {
	client    *http.Client
	AuthToken string
	BaseUrl   string
}

func NewEventsAPI(authToken, url string) *EventsAPI {
	log.Println("New Events API Version:", Version)
	return &EventsAPI{
		client:    &http.Client{},
		AuthToken: authToken,
		BaseUrl:   url,
	}
}

func (e *EventsAPI) request(ctx context.Context, method string, route string, body interface{}) (*http.Response, error) {
	var b io.Reader
	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			err := fmt.Errorf("could not marshal request: %w", err)
			panic(err)
		}
		b = bytes.NewReader(reqBody)
	}
	req, err := http.NewRequestWithContext(ctx, method, e.BaseUrl+route, b)
	if err != nil {
		err := fmt.Errorf("could not create new request: %w", err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", e.AuthToken))
	req.Header.Add("User-Agent", DefaultUserAgent)
	res, err := e.client.Do(req)
	if err != nil {
		err := fmt.Errorf("could not make request: %w", err)
		return nil, err
	}
	return res, nil
}

type CursorRequest struct {
	Cursor string `json:"cursor"`
}
type CursorResetRequest struct {
	Limit     int        `json:"limit"`
	StartTime *time.Time `json:"start_time,omitempty"`
}
