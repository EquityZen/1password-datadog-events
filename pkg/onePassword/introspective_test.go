package onePassword

import (
	"context"
	"encoding/json"
	ddc "equityzen/1password-datadog-events/pkg/datadog"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"os"
	"testing"
	"time"
)

func TestGetIntrospectiveUsage(t *testing.T) {
	c := NewEventsAPI(Token, "https://events.1password.com")

	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: os.Getenv("DD_CLIENT_API_KEY"),
			},
			"appKeyAuth": {
				Key: os.Getenv("DD_CLIENT_APP_KEY"),
			},
		},
	)
	tt := time.Now().Add(-148 * time.Hour)

	resp, _ := c.ItemUsagesRequest(ctx, CursorResetRequest{Limit: 1, StartTime: &tt})

	//resp2, _ := c.SignInAttemptsRequest(ctx, api.CursorResetRequest{Limit: 10, StartTime: &tt})
	//fmt.Println(resp.PrintEvents())
	connectToken := os.Getenv("CONNECT_TOKEN")
	cc := NewConnectAPI(connectToken, "http://127.0.0.1:8081")
	ddevents := ddc.NewDataDogAPI()
	for _, i := range resp.Items {
		connectItem, err := cc.RetrieveItemByTitle(i.ItemUUID, i.VaultUUID)
		if err != nil {
			fmt.Println("Errorrr: ", err)
			t.Fail()
		}
		fmt.Println(i.Action, i.User.Email, connectItem.Category)
		payload := &LogMessage{
			Name:     i.User.Name,
			URL:      connectItem.URLs[0].URL,
			Title:    connectItem.Title,
			Action:   i.Action,
			Email:    i.User.Email,
			Category: string(connectItem.Category),
		}

		b, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("unable to convert json payload")

		}
		logItem := ddevents.CreateLogItem(string(b))
		fmt.Println(logItem)
		ddevents.PostLog(ctx, logItem)
		//fmt.Println(i.User, i.VaultUUID, i.ItemUUID, i.Action)
		//evnetMsg := ddevents.CreateEventRequest("1Password Event triggered", message)
		//ddevents.PostEvent(ctx, evnetMsg)
	}
}

func TestConnectAPI_RetrieveItemByTitle(t *testing.T) {
	connectToken := os.Getenv("CONNECT_TOKEN")
	c := NewConnectAPI(connectToken, "http://127.0.0.1:8081")
	item, err := c.RetrieveItemByTitle("uvntxudapvff7ijashdwecw2gy", "otl2e6sadb53znmvvgght3ldgy")

	if err != nil {
		fmt.Println("Errorrr: ", err)
		t.Fail()
	}
	fmt.Println(item.Fields[1])
}
