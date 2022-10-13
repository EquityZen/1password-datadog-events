package main

import (
	"context"
	datadog2 "equityzen/1password-datadog-events/pkg/datadog"
	"equityzen/1password-datadog-events/pkg/onepassword"
	"fmt"
	"github.com/sirupsen/logrus"
)

func SendSignInAttemptsToDD(body *onepassword.SignInAttemptResponse, client *datadog2.DDClient, ddc context.Context, l *logrus.Logger) {
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

func SendItemsUsageToDD(body *onepassword.ItemUsageResponse, client *datadog2.DDClient, api *onepassword.ConnectAPI, ddc context.Context, l *logrus.Logger) {
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
