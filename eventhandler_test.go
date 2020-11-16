package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	keptnlib "github.com/keptn/go-utils/pkg/lib"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
)

/**
 * loads a cloud event from the passed test json file and initializes a keptn object with it
 */
func initializeTestObjects(eventFileName string) (*keptnlib.Keptn, *cloudevents.Event, error) {
	// load sample event
	eventFile, err := ioutil.ReadFile(eventFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant load %s: %s", eventFileName, err.Error())
	}

	incomingEvent := &cloudevents.Event{}
	err = json.Unmarshal(eventFile, incomingEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing: %s", err.Error())
	}

	var keptnOptions = keptn.KeptnOpts{}
	keptnOptions.UseLocalFileSystem = true
	myKeptn, err := keptnlib.NewKeptn(incomingEvent, keptnOptions)

	return myKeptn, incomingEvent, err
}

// Handles ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"
func TestHandleConfigureMonitoringEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/configure-monitoring.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.ConfigureMonitoringEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleConfigureMonitoringEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func TestHandleConfigurationChangeEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/configuration-change.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.ConfigurationChangeEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleConfigurationChangeEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func TestHandleDeploymentFinishedEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/deployment-finished.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.DeploymentFinishedEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleDeploymentFinishedEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func TestHandleTestsFinishedEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/tests-finished.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.TestsFinishedEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleTestsFinishedEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func TestHandleStartEvaluationEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/start-evaluation.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.StartEvaluationEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleStartEvaluationEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func TestHandleEvaluationDoneEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/evaluation-done.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.EvaluationDoneEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleEvaluationDoneEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

// Tests the InternalGetSLIEvent Handler
func TestHandleInternalGetSLIEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/get-sli.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.InternalGetSLIEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleInternalGetSLIEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func TestHandleProblemEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/problem.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.ProblemEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleProblemEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

//
// Handles ActionTriggeredEventType = "sh.keptn.event.action.triggered"
// TODO: add in your handler code
//
func TestHandleActionTriggeredEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/action-triggered.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnlib.ActionTriggeredEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleActionTriggeredEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}
