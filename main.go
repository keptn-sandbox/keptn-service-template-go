package main

import (
	"github.com/keptn/go-utils/pkg/sdk"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

const eventTypeWildcard = "sh.keptn.event.action.triggered,sh.keptn.event.get-sli.triggered" // Defines the event types the service will listen to
const serviceName = "keptn-service-template-go"
const envVarLogLevel = "LOG_LEVEL"

func main() {
	if os.Getenv(envVarLogLevel) != "" {
		logLevel, err := logrus.ParseLevel(os.Getenv(envVarLogLevel))
		if err != nil {
			logrus.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
			logrus.SetLevel(logrus.InfoLevel)
		} else {
			logrus.SetLevel(logLevel)
		}
	}

	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithTaskHandler(
			eventTypeWildcard,
			NewEventHandler()),
		sdk.WithLogger(logrus.New()),
	).Start())
}
