package main

import (
	"context"
	datadog2 "equityzen/1password-datadog-events/pkg/datadog"
	"equityzen/1password-datadog-events/pkg/onePassword"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
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
	e := onePassword.NewEventsAPI(onePassword.Token, "https://events.1password.com")
	eCache := cache.New(-1, -1)

	dd := datadog2.NewDataDogAPI()
	cc := onePassword.NewConnectAPI(viper.GetString("connect_token"), "http://127.0.0.1:8081")

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
		var itemUseResp *onePassword.ItemUsageResponse
		var err error
		s.log.Infoln("Retrieving items usage")
		if cursor, found := eCache.Get(ITEMUSAGEC); found {
			s.log.Infoln("Found previous cursor")
			itemUseResp, err = e.ItemUsagesRequest(eCTX, onePassword.CursorRequest{Cursor: cursor.(string)})
			if err != nil {
				s.log.Warnln("Error has occurred retrieving item usage: ", err)
			}
		} else {
			s.log.Infoln("Could not find previous cursor resetting cursor")
			itemUseResp, err = e.ItemUsagesRequest(eCTX, onePassword.CursorResetRequest{Limit: 100, StartTime: &tt})
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
		var signInResp *onePassword.SignInAttemptResponse
		if cursor, found := eCache.Get(SIGNINATTEMPTC); found {
			s.log.Infoln("Found previous cursor")
			signInResp, err = e.SignInAttemptsRequest(eCTX, onePassword.CursorRequest{Cursor: cursor.(string)})
			if err != nil {
				s.log.Warnln("Error has occurred retrieving sign in attempts: ", err)
			}
		} else {
			s.log.Infoln("Could not find previous cursor resetting cursor")
			signInResp, err = e.SignInAttemptsRequest(eCTX, onePassword.CursorResetRequest{Limit: 100, StartTime: &tt})
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

func SendSignInAttemptsToDD(body *onePassword.SignInAttemptResponse, client *datadog2.DDClient, ddc context.Context, l *logrus.Logger) {
	if len(body.Items) == 0 {
		l.Infoln("No login events found in time slot")
	}
	for _, i := range body.Items {
		EventAttributes := map[string]string{
			"Name":        i.TargetUser.Name,
			"Email":       i.TargetUser.Email,
			"Category":    i.Category,
			"Type":        i.Type,
			"Application": i.Client.AppName,
			"OS":          i.Client.OSName,
			//"SignInAttempts": i.Details.Value,
		}
		client.PostLog(ddc, client.CreateLogItem("1Password User Sign in Attempt", EventAttributes))
	}
}

func SendItemsUsageToDD(body *onePassword.ItemUsageResponse, client *datadog2.DDClient, api *onePassword.ConnectAPI, ddc context.Context, l *logrus.Logger) {
	if len(body.Items) == 0 {
		l.Infoln("No item usage events found in time slot")
	}
	for _, i := range body.Items {
		connectItem, err := api.RetrieveItemByTitle(i.ItemUUID, i.VaultUUID)
		if err != nil {
			l.Warnln("Error occured while trying to retrive connect info: ", err)
			continue
		}
		EventAttribute := map[string]string{
			"Name":     i.User.Name,
			"Title":    connectItem.Title,
			"Email":    i.User.Email,
			"Category": string(connectItem.Category),
			"Action":   i.Action,
			"Version":  fmt.Sprintf("%v", i.UsedVersion),
		}

		client.PostLog(ddc, client.CreateLogItem("1Password Item usage event", EventAttribute))
	}
}
