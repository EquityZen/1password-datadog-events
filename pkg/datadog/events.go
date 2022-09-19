package datadog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"os"
)

type DDClient struct {
	eventsApi datadogV1.EventsApi
}

func NewDataDogApiClient() *datadog.APIClient {
	config := datadog.NewConfiguration()
	return datadog.NewAPIClient(config)
}

func NewDataDogEventsAPI() *DDClient {
	return &DDClient{*datadogV1.NewEventsApi(NewDataDogApiClient())}
}

func (e *DDClient) CreateEventRequest(title, text string) datadogV1.EventCreateRequest {
	return datadogV1.EventCreateRequest{
		Tags:  []string{"1PassWordEvents"},
		Text:  text,
		Title: title,
	}
}

func (e *DDClient) PostEvent(ctx context.Context, event datadogV1.EventCreateRequest) {
	resp, r, err := e.eventsApi.CreateEvent(ctx, event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EventsApi.CreateEvent`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `EventsApi.CreateEvent`:\n%s\n", responseContent)
}
