package utils

import (
	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"log"
	"time"
)

func SentryInit() {
	var debug bool = viper.GetBool("debug")
	var dsn string = viper.GetString("sentry.dsn")
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		Debug: debug,
	})

	if err != nil {
		log.Fatalf("Sentry.Init: %s\n", err)
	}

	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()

	if debug {
		sentry.CaptureMessage("Dev listener has been started!")
	} else {
		sentry.CaptureMessage("Listener has been started!")
	}
}

func CaptureSentryException(exception error) {
	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()
	sentry.CaptureException(exception)
}