package main

import (
	"github.com/keptn-service-template-go/handler"
	"github.com/keptn/go-utils/pkg/sdk"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

const getSliTriggeredEvent = "sh.keptn.event.get-sli.triggered"
const actionTriggeredEvent = "sh.keptn.event.action.triggered"
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
			actionTriggeredEvent,
			handler.NewActionTriggeredEventHandler()),
		sdk.WithTaskHandler(
			getSliTriggeredEvent,
			handler.NewGetSliEventHandler()),
		sdk.WithLogger(logrus.New()),
	).Start())
}
