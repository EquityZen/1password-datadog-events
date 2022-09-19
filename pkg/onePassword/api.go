package onePassword

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var Token = "eyJhbGciOiJFUzI1NiIsImtpZCI6Inh5YmJzM3MzZGV3YXNjeXNwbmNxa3lqenV1IiwidHlwIjoiSldUIn0.eyIxcGFzc3dvcmQuY29tL2F1dWlkIjoiUVgzU1EzNEQ2NUJCWE9ISFJUT0RCTTM3Wk0iLCIxcGFzc3dvcmQuY29tL2Z0cyI6WyJzaWduaW5hdHRlbXB0cyIsIml0ZW11c2FnZXMiXSwiYXVkIjpbImV2ZW50cy4xcGFzc3dvcmQuY29tIl0sInN1YiI6Ilk0SEE2T0s1WlZGUEpOS0M1SFRXUVZYRFVBIiwiaWF0IjoxNjYzMTY5NDMxLCJpc3MiOiJjb20uMXBhc3N3b3JkLmI1IiwianRpIjoiZnZobXM0cGl0enA3ZmdlcHl3MmJxYXJzN20ifQ.Ws1PlHxJ1Gvk5VG6ve1Ob43eWNmw800OT-qIjut8wa_1fIb6qnj_bVkbzOdhpFaHeZQI0POvFDZjaWsTjU_FBQ"
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
