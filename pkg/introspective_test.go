package pkg

import (
	"context"
	"equityzen/events-api-splunk/src/api"
	"fmt"
	"testing"
	"time"
)

var storeLocation map[float64]float64

func TestGetIntrospectiveUsage(t *testing.T) {
	c := NewEventsAPI(Token, "https://events.1password.com")
	ctx := context.Background()
	tt := time.Now().Add(-148 * time.Hour)

	resp, _ := c.ItemUsagesRequest(ctx, api.CursorResetRequest{Limit: 100, StartTime: &tt})

	//resp2, _ := c.SignInAttemptsRequest(ctx, api.CursorResetRequest{Limit: 10, StartTime: &tt})
	fmt.Println(resp.PrintEvents())
}
