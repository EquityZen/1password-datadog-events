package datadog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"testing"
)

func TestDDClient_PostLog(t *testing.T) {
	ddctx := context.WithValue(context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: "",
			},
		},
	)

	payloadMessage := map[string]string{
		"Name":     "TestName",
		"URL":      "TestURL",
		"Title":    "TestTitle",
		"Action":   "TestAction",
		"Email":    "TestEmail",
		"Category": "TestCatagory",
	}
	b, _ := json.Marshal(payloadMessage)
	ddc := NewDataDogAPI()
	fmt.Println(string(b))
	payload := ddc.CreateLogItem("test", payloadMessage)
	ddc.PostLog(ddctx, payload)
}
