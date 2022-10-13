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
	logApi    datadogV1.LogsApi
}

func NewDataDogApiClient() *datadog.APIClient {
	config := datadog.NewConfiguration()
	return datadog.NewAPIClient(config)
}

func NewDataDogAPI() *DDClient {
	return &DDClient{*datadogV1.NewEventsApi(NewDataDogApiClient()), *datadogV1.NewLogsApi(NewDataDogApiClient())}
}

func (e *DDClient) CreateEventRequest(title, text string) datadogV1.EventCreateRequest {
	return datadogV1.EventCreateRequest{
		Tags:  []string{"1PassWordEvents"},
		Text:  text,
		Title: title,
	}
}

func (e *DDClient) CreateLogItem(payload string, attributes map[string]string) datadogV1.HTTPLogItem {
	return datadogV1.HTTPLogItem{
		Ddsource:             datadog.PtrString("1Password"),
		Ddtags:               datadog.PtrString("env:infra, version:beta"),
		Hostname:             datadog.PtrString("localhost"),
		Message:              payload,
		Service:              datadog.PtrString("1passwordevents"),
		AdditionalProperties: attributes,
	}
}

func (e *DDClient) PostLog(ctx context.Context, log datadogV1.HTTPLogItem) {
	resp, r, err := e.logApi.SubmitLog(ctx, []datadogV1.HTTPLogItem{log}, *datadogV1.NewSubmitLogOptionalParameters().WithContentEncoding(datadogV1.CONTENTENCODING_DEFLATE))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LogsApi.SubmitLog`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `LogsApi.SubmitLog`:\n%s\n", responseContent)

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
