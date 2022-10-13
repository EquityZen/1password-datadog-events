package onePassword

import (
	"context"
	ddc "equityzen/1password-datadog-events/pkg/datadog"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	cache "github.com/patrickmn/go-cache"
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
				Key: "",
			},
		},
	)
	tt := time.Now().Add(-148 * time.Hour)

	resp, _ := c.ItemUsagesRequest(ctx, CursorResetRequest{Limit: 100, StartTime: &tt})

	fmt.Println("Curser: ", resp.CursorResponse.Cursor)

	if len(resp.Items) == 0 {
		fmt.Println("No events found")
		t.Fail()
	}
	//resp2, _ := c.SignInAttemptsRequest(ctx, api.CursorResetRequest{Limit: 10, StartTime: &tt})
	//fmt.Println(resp.PrintEvents())
	connectToken := ``
	cc := NewConnectAPI(connectToken, "http://127.0.0.1:8081")
	ddevents := ddc.NewDataDogAPI()
	cci := cache.New(time.Minute*5, time.Minute*10)
	for x, i := range resp.Items {
		fmt.Println("Iteration: ", x)
		cci.Set(i.ItemUUID, i.VaultUUID, cache.NoExpiration)
		connectItem, err := cc.RetrieveItemByTitle(i.ItemUUID, i.VaultUUID)
		if err != nil {
			fmt.Println("Errorrr: ", err)
			continue
		}

		//if connectItem.URLs != nil {
		//	EventAttribute["URL"] = connectItem.URLs[0].URL
		//}
		EventAttribute := map[string]string{
			"Name":     i.User.Name,
			"Title":    connectItem.Title,
			"Email":    i.User.Email,
			"Category": string(connectItem.Category),
			"Action":   i.Action,
			"Version":  fmt.Sprintf("%v", i.UsedVersion),
		}

		fmt.Println("posting datadog events")
		logItem := ddevents.CreateLogItem("Event triggered for 1Passowrd", EventAttribute)
		fmt.Println(logItem)
		ddevents.PostLog(ctx, logItem)
	}
}

func TestConnectAPI_RetrieveItemByTitle(t *testing.T) {
	connectToken := ``
	c := NewConnectAPI(connectToken, "http://127.0.0.1:8081")
	item, err := c.RetrieveItemByTitle("uvntxudapvff7ijashdwecw2gy", "otl2e6sadb53znmvvgght3ldgy")

	if err != nil {
		fmt.Println("Errorrr: ", err)
		t.Fail()
	}
	fmt.Println(item)
}
