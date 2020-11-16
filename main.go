package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	keptnlib "github.com/keptn/go-utils/pkg/lib"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
)

var keptnOptions = keptn.KeptnOpts{}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// Path to which cloudevents are sent
	Path string `envconfig:"RCV_PATH" default:"/"`
	// Whether we are running locally (e.g., for testing) or on production
	Env string `envconfig:"ENV" default:"local"`
	// URL of the Keptn configuration service (this is where we can fetch files from the config repo)
	ConfigurationServiceUrl string `envconfig:"CONFIGURATION_SERVICE" default:""`
	// URL of the Keptn event broker (this is where this service sends cloudevents to)
	EventBrokerUrl string `envconfig:"EVENTBROKER" default:""`
}

/**
 * This method gets called when a new event is received from the Keptn Event Distributor
 * Depending on the Event Type will call the specific event handler functions, e.g: handleDeploymentFinishedEvent
 * See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload
 */
func processKeptnCloudEvent(ctx context.Context, event cloudevents.Event) error {
	myKeptn, err := keptnlib.NewKeptn(&event, keptnOptions)

	log.Printf("gotEvent(%s): %s - %s", event.Type(), myKeptn.KeptnContext, event.Context.GetID())

	if err != nil {
		log.Printf("failed to parse incoming cloudevent: %v", err)
		return err
	}

	// ********************************************
	// Lets test on each possible Event Type and call the respective handler function
	// ********************************************
	if event.Type() == keptnlib.ConfigurationChangeEventType {
		log.Printf("Processing Configuration Change Event")

		configChangeEventData := &keptnlib.ConfigurationChangeEventData{}
		err := event.DataAs(configChangeEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleConfigurationChangeEvent(myKeptn, event, configChangeEventData)
	} else if event.Type() == keptnlib.DeploymentFinishedEventType {
		log.Printf("Processing Deployment Finished Event")

		deployFinishEventData := &keptnlib.DeploymentFinishedEventData{}
		err := event.DataAs(deployFinishEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleDeploymentFinishedEvent(myKeptn, event, deployFinishEventData)
	} else if event.Type() == keptnlib.TestsFinishedEventType {
		log.Printf("Processing Test Finished Event")

		testsFinishedEventData := &keptnlib.TestsFinishedEventData{}
		err := event.DataAs(testsFinishedEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleTestsFinishedEvent(myKeptn, event, testsFinishedEventData)
	} else if event.Type() == keptnlib.StartEvaluationEventType {
		log.Printf("Processing Start Evaluation Event")

		startEvaluationEventData := &keptnlib.StartEvaluationEventData{}
		err := event.DataAs(startEvaluationEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleStartEvaluationEvent(myKeptn, event, startEvaluationEventData)
	} else if event.Type() == keptnlib.EvaluationDoneEventType {
		log.Printf("Processing Evaluation Done Event")

		evaluationDoneEventData := &keptnlib.EvaluationDoneEventData{}
		err := event.DataAs(evaluationDoneEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleEvaluationDoneEvent(myKeptn, event, evaluationDoneEventData)
	} else if event.Type() == keptnlib.ProblemOpenEventType || event.Type() == keptnlib.ProblemEventType {
		// Subscribing to a problem.open or problem event is deprecated since Keptn 0.7 - subscribe to sh.keptn.event.action.triggered
		log.Printf("Subscribing to a problem.open or problem event is not recommended since Keptn 0.7. Please subscribe to event of type: sh.keptn.event.action.triggered")
		log.Printf("Processing Problem Event")

		problemEventData := &keptnlib.ProblemEventData{}
		err := event.DataAs(problemEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleProblemEvent(myKeptn, event, problemEventData)
	} else if event.Type() == keptnlib.ActionTriggeredEventType {
		log.Printf("Processing Action Triggered Event")

		actionTriggeredEventData := &keptnlib.ActionTriggeredEventData{}
		err := event.DataAs(actionTriggeredEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleActionTriggeredEvent(myKeptn, event, actionTriggeredEventData)
	} else if event.Type() == keptnlib.ConfigureMonitoringEventType {
		log.Printf("Processing Configure Monitoring Event")

		configureMonitoringEventData := &keptnlib.ConfigureMonitoringEventData{}
		err := event.DataAs(configureMonitoringEventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return HandleConfigureMonitoringEvent(myKeptn, event, configureMonitoringEventData)
	}

	// Unknown Event -> Throw Error!
	var errorMsg string
	errorMsg = fmt.Sprintf("Unhandled Keptn Cloud Event: %s", event.Type())

	log.Print(errorMsg)
	return errors.New(errorMsg)
}

/**
 * Usage: ./main
 * no args: starts listening for cloudnative events on localhost:port/path
 *
 * Environment Variables
 * env=runlocal   -> will fetch resources from local drive instead of configuration service
 */
func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	os.Exit(_main(os.Args[1:], env))
}

/**
 * Opens up a listener on localhost:port/path and passes incoming requets to gotEvent
 */
func _main(args []string, env envConfig) int {
	ctx := context.Background()

	// configure keptn options
	if env.Env == "local" {
		log.Println("env=local: Running with local filesystem to fetch resources")
		keptnOptions.UseLocalFileSystem = true
	}

	keptnOptions.ConfigurationServiceURL = env.ConfigurationServiceUrl
	keptnOptions.EventBrokerURL = env.EventBrokerUrl

	// configure http server to receive cloudevents
	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	log.Println("Starting keptn-service-template-go...")
	log.Printf("    on Port = %d; Path=%s", env.Port, env.Path)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, processKeptnCloudEvent))

	return 0
}
