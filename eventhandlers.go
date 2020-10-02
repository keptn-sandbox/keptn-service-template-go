package main

import (
	"log"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

/**
* Here are all the handler functions for the individual event
  See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload

  -> "sh.keptn.event.configuration.change"
  -> "sh.keptn.events.deployment-finished"
  -> "sh.keptn.events.tests-finished"
  -> "sh.keptn.event.start-evaluation"
  -> "sh.keptn.events.evaluation-done"
  -> "sh.keptn.event.problem.open"
	-> "sh.keptn.events.problem"
	-> "sh.keptn.event.action.triggered"
*/

// Handles ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"
func HandleConfigureMonitoringEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigureMonitoringEventData) error {
	log.Printf("Handling Configure Monitoring Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func HandleConfigurationChangeEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigurationChangeEventData) error {
	log.Printf("Handling Configuration Changed Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleDeploymentFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.DeploymentFinishedEventData) error {
	log.Printf("Handling Deployment Finished Event: %s", incomingEvent.Context.GetID())

	// capture start time for tests
	// startTime := time.Now()

	// run tests
	// ToDo: Implement your tests here

	// Send Test Finished Event
	// return myKeptn.SendTestsFinishedEvent(&incomingEvent, "", "", startTime, "pass", nil, "keptn-service-template-go")
	return nil
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func HandleStartEvaluationEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.StartEvaluationEventData) error {
	log.Printf("Handling Start Evaluation Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleEvaluationDoneEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.EvaluationDoneEventData) error {
	log.Printf("Handling Evaluation Done Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles InternalGetSLIEventType = "sh.keptn.internal.event.get-sli"
// TODO: add in your handler code
//
func HandleInternalGetSLIEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.InternalGetSLIEventData) error {
	log.Printf("Handling Internal Get SLI Event: %s", incomingEvent.Context.GetID())

	incomingGetSLIEventData := &keptn.InternalGetSLIEventData{}
	incomingEvent.DataAs(incomingGetSLIEventData)

	// Step 1 - Do we need to do something?
	// Lets make sure we are only processing an event that really belongs to our SLI Provider
	/* if incomingGetSLIEventData.SLIProvider != "keptn-service-template-go" {
		return nil
	}*/

	// Step 2 - prep-work
	// Get any additional input / configuration data, e.g
	// Labels: get the incoming labels for potential config data and use it to pass more labels on result, e.g: links
	// SLI.yaml: if your service uses SLI.yaml to store query definitions for SLIs get that file from Keptn
	/* labels := incomingGetSLIEventData.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	testRunID := labels["testRunId"]*/

	// sliConfigFileContent, err := myKeptn.GetKeptnResource("keptn-service-template-go/sli.yaml")

	// Step 3 - do your work - iterate through the list of requested indicators and return their values
	// Indicators: this is the list of indicators as requested in the SLO.yaml
	// SLIResult: this is the array that will receive the results
	/* indicators := incomingGetSLIEventData.Indicators
	sliResults := []*keptn.SLIResult{}

	for _, indicatorName := range indicators {
		sliResult := &keptn.SLIResult{
			Metric: indicatorName,
			Value:  123.4,
		}
		sliResults = append(sliResults, sliResult)
	}*/

	// Step 4 - add additional context via labels
	// labels["Link to Data Source"] = "https://mydatasource/myquery?testRun=" + testRunID

	// Step 4 - send results back to Keptn
	// return myKeptn.SendInternalGetSLIDoneEvent(incomingGetSLIEventData, sliResults, labels, err, "keptn-service-template-go")

	return nil
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func HandleProblemEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

//
// Handles ActionTriggeredEventType = "sh.keptn.event.action.triggered"
// TODO: add in your handler code
//
func HandleActionTriggeredEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())

	// check if action is supported
	if data.Action.Action == "action-xyz" {
		//myKeptn.SendActionStartedEvent()

		// Implement your remediation action here

		//myKeptn.SendActionFinishedEvent()
	}
	return nil
}
