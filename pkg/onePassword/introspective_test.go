package onePassword

import (
	"context"
	dd "equityzen/1password-datadog-events/pkg/datadog"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
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
				Key: "637c7abe12b881a66dbe9f4ff6c1c6e3",
			},
		},
	)
	tt := time.Now().Add(-148 * time.Hour)

	resp, _ := c.ItemUsagesRequest(ctx, CursorResetRequest{Limit: 1, StartTime: &tt})

	//resp2, _ := c.SignInAttemptsRequest(ctx, api.CursorResetRequest{Limit: 10, StartTime: &tt})
	//fmt.Println(resp.PrintEvents())

	ddevents := dd.NewDataDogEventsAPI()
	for _, i := range resp.Items {
		//fmt.Println(i.User, i.VaultUUID, i.ItemUUID, i.Action)
		evnetMsg := ddevents.CreateEventRequest("1Password Event triggered")
		ddevents.PostEvent(ctx, evnetMsg)
	}
}

func TestConnectAPI_RetrieveItemByTitle(t *testing.T) {
	c := NewConnectAPI(Token, "https://equityzen.1password.com/")
	item, err := c.RetrieveItemByTitle("y6cm627qxzatbfcwjqob7yz4kq", "otl2e6sadb53znmvvgght3ldgy")

	if err != nil {
		fmt.Println("Errorrr: ", err)
		t.Fail()
	}
	fmt.Println(item)
}
