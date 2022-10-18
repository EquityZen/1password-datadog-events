package main

import (
	"fmt"
	"github.com/EquityZen/ez-go/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Service struct {
	log       *logrus.Logger
	vaultData map[string]interface{}
}

func InitSettings(name string) Service {
	var configFilePath string
	var service Service
	pflag.StringVar(&configFilePath, "config_file_name", fmt.Sprintf("%v.yaml", name), "Name of configuration file")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()

	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigName(name)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	//Default values

	viper.SetDefault("events_api", "https://events.1password.com")
	viper.SetDefault("connect_api", "onepassword-connect.1password")
	viper.SetDefault("vault_address", "https://vault.equityzen.com")
	viper.SetDefault("secret_path", "secret/onepassword")
	viper.SetDefault("cron_schedule", "*/10 * * * *")

	// Read in config
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Printf("Failed to read in %v.yaml, all configs will need to be provided in the ENV", name)
	}
	log := logrus.StandardLogger()
	service.log = log

	if vc, err := vault.InitVault(
		vault.WithAddress(viper.GetString("vault_address")),
		vault.WithRoleID(viper.GetString("role_id")),
		vault.WithSecretID(viper.GetString("secret_id")),
	); err != nil {
		log.Fatalln("Failed to initialize vault connection: ", err)
	} else {
		data, err := vc.LoadSecretsAtPath(viper.GetString("secret_path"))
		if err != nil {
			log.Fatalf("Failed to load vault data at path: %v, %v", viper.GetString("secret_path"), err)
		}
		for propName, propValue := range data {
			viper.Set(propName, propValue)
		}
	}
	viper.AutomaticEnv()

	// defaults
	viper.SetDefault("log_level", logrus.InfoLevel)

	return service
}
