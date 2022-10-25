package main

import (
	"context"
	datadog2 "equityzen/1password-datadog-events/pkg/datadog"
	"equityzen/1password-datadog-events/pkg/onepassword"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	SIGNINATTEMPTC = "SIACURSOR"
	ITEMUSAGEC     = "USAGEC"
)

func main() {

	s := InitSettings("onepassword")
	e := onepassword.NewEventsAPI(onepassword.Token, "https://events.1password.com")
	eCache := cache.New(-1, -1)

	dd := datadog2.NewDataDogAPI()
	cc := onepassword.NewConnectAPI(viper.GetString("connect_token"), "http://127.0.0.1:8081")

	c := cron.New(
		cron.WithLocation(time.UTC),
	)

	ddCTX := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: viper.GetString("dd_api_key"),
			},
		})

	eCTX := context.Background()
	tt := time.Now().Add(-10 * time.Minute)

	_, err := c.AddFunc(viper.GetString("cron_schedule"), func() {
		var itemUseResp *onepassword.ItemUsageResponse
		var err error
		s.log.Infoln("Retrieving items usage")
		if cursor, found := eCache.Get(ITEMUSAGEC); found {
			s.log.Infoln("Found previous cursor")
			itemUseResp, err = e.ItemUsagesRequest(eCTX, onepassword.CursorRequest{Cursor: cursor.(string)})
			if err != nil {
				s.log.Warnln("Error has occurred retrieving item usage: ", err)
			}
		} else {
			s.log.Infoln("Could not find previous cursor resetting cursor")
			itemUseResp, err = e.ItemUsagesRequest(eCTX, onepassword.CursorResetRequest{Limit: 100, StartTime: &tt})
			if err != nil {
				s.log.Warnln("Error has occured retrieving item usage: ", err)
			}
		}

		s.log.Infoln("Sending ItemsUsage")
		SendItemsUsageToDD(itemUseResp, dd, cc, ddCTX, s.log)
		eCache.Set(ITEMUSAGEC, itemUseResp.Cursor, -1)

	})

	if err != nil {
		s.log.Fatalln("Failed to schedule item usage cron: ", err)
	}

	_, err2 := c.AddFunc(viper.GetString("cron_schedule"), func() {
		var signInResp *onepassword.SignInAttemptResponse
		if cursor, found := eCache.Get(SIGNINATTEMPTC); found {
			s.log.Infoln("Found previous cursor")
			signInResp, err = e.SignInAttemptsRequest(eCTX, onepassword.CursorRequest{Cursor: cursor.(string)})
			if err != nil {
				s.log.Warnln("Error has occurred retrieving sign in attempts: ", err)
			}
		} else {
			s.log.Infoln("Could not find previous cursor resetting cursor")
			signInResp, err = e.SignInAttemptsRequest(eCTX, onepassword.CursorResetRequest{Limit: 100, StartTime: &tt})
			if err != nil {
				s.log.Warnln("Error has occurred retrieving sign in attempts: ", err)
			}
		}

		SendSignInAttemptsToDD(signInResp, dd, ddCTX, s.log)
		eCache.Set(SIGNINATTEMPTC, signInResp.Cursor, -1)
	})

	if err2 != nil {
		s.log.Fatalln("Failed to schedule sign in attempts cron: ", err)
	}

	c.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
