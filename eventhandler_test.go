package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
)

/**
 * loads a cloud event from the passed test json file and initializes a keptn object with it
 */
func initializeTestObjects(eventFileName string) (*keptnv2.Keptn, *cloudevents.Event, error) {
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
	myKeptn, err := keptnv2.NewKeptn(incomingEvent, keptnOptions)

	return myKeptn, incomingEvent, err
}

// Tests HandleActionTriggeredEvent
// TODO: Add your test-code
func TestHandleActionTriggeredEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/action.triggered.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnv2.ActionTriggeredEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleActionTriggeredEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

// Tests HandleDeploymentTriggeredEvent
// TODO: Add your test-code
func TestHandleDeploymentTriggeredEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/evaluation.triggered.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnv2.DeploymentTriggeredEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleDeploymentTriggeredEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

// Tests HandleEvaluationTriggeredEvent
// TODO: Add your test-code
func TestHandleEvaluationTriggeredEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/evaluation.triggered.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnv2.EvaluationTriggeredEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleEvaluationTriggeredEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

// Tests the HandleGetSliTriggeredEvent Handler
// TODO: Add your test-code
func TestHandleGetSliTriggered(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/get-sli.triggered.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnv2.GetSLITriggeredEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleGetSliTriggeredEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}

// Tests the HandleReleaseTriggeredEvent Handler
// TODO: Add your test-code
func TestHandleReleaseTriggeredEvent(t *testing.T) {
	myKeptn, incomingEvent, err := initializeTestObjects("test-events/release.triggered.json")
	if err != nil {
		t.Error(err)
		return
	}

	specificEvent := &keptnv2.ReleaseTriggeredEventData{}
	err = incomingEvent.DataAs(specificEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	err = HandleReleaseTriggeredEvent(myKeptn, *incomingEvent, specificEvent)
	if err != nil {
		t.Errorf("Error: " + err.Error())
	}
}
