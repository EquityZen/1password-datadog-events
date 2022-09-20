package cmd

import (
	"fmt"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

func main() {
	initSettings()
	log := initLogging()
	var c = cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(log)),
		cron.WithLocation(time.UTC),
	)

}

// initLogging configure logging
func initLogging() *logrus.Logger {
	if level, err := logrus.ParseLevel(viper.GetString("log_level")); err != nil {
		fmt.Printf("Failed to parse log_level: %v\n", err)
		os.Exit(1)
	} else {
		logrus.SetLevel(level)
	}

	// On debug add file/line
	if logrus.GetLevel() == logrus.DebugLevel {
		filenameHook := filename.NewHook()
		filenameHook.Field = "fileline"
		logrus.AddHook(filenameHook)
	}

	// format
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// output consul
	logrus.SetOutput(os.Stdout)

	return logrus.StandardLogger()
}

func initSettings() {
	viper.AutomaticEnv()

	// defaults
	viper.SetDefault("log_level", logrus.InfoLevel)
	viper.SetDefault("schedule_lookup", "5 * * * *")
}
