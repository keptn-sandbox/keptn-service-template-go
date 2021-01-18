package main

import (
	"fmt"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

// GenericLogKeptnCloudEventHandler is a generic handler for Keptn Cloud Events that logs the CloudEvent
func GenericLogKeptnCloudEventHandler(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data interface{}) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())
	log.Printf("CloudEvent %T: %v", data, data)

	return nil
}

// OldHandleConfigureMonitoringEvent handles old configure-monitoring events
// TODO: add in your handler code
func OldHandleConfigureMonitoringEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigureMonitoringEventData) error {
	log.Printf("Handling old configure-monitoring Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleConfigureMonitoringTriggeredEvent handles configure-monitoring.triggered events
// TODO: add in your handler code
func HandleConfigureMonitoringTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ConfigureMonitoringTriggeredEventData) error {
	log.Printf("Handling configure-monitoring.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleDeploymentTriggeredEvent handles deployment.triggered events
// TODO: add in your handler code
func HandleDeploymentTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.DeploymentTriggeredEventData) error {
	log.Printf("Handling deployment.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleTestTriggeredEvent handles test.triggered events
// TODO: add in your handler code
func HandleTestTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.TestTriggeredEventData) error {
	log.Printf("Handling test.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleApprovalTriggeredEvent handles approval.triggered events
// TODO: add in your handler code
func HandleApprovalTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ApprovalTriggeredEventData) error {
	log.Printf("Handling approval.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleEvaluationTriggeredEvent handles evaluation.triggered events
// TODO: add in your handler code
func HandleEvaluationTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationTriggeredEventData) error {
	log.Printf("Handling evaluation.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleReleaseTriggeredEvent handles release.triggered events
// TODO: add in your handler code
func HandleReleaseTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ReleaseTriggeredEventData) error {
	log.Printf("Handling release.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleRemediationTriggeredEvent handles remediation.triggered events
// TODO: add in your handler code
func HandleRemediationTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.RemediationTriggeredEventData) error {
	log.Printf("Handling remediation.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// sendGetSliFinishedCloudEvent is a helper function to send a get-sli.finished event
// ToDo: This will need to be refactored once https://github.com/keptn/keptn/issues/2913 is ready
func sendGetSliFinishedCloudEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.GetSLITriggeredEventData,
	status keptnv2.StatusType, result keptnv2.ResultType, message string) error {
	log.Printf("Sending getSli Finished Cloud Event with status=%s and result=%s back to Keptn (%s)", status, result, message)
	getSliFinishedEventData := keptnv2.GetSLIFinishedEventData{}

	getSliFinishedEventData.EventData = data.EventData
	getSliFinishedEventData.Status = status
	getSliFinishedEventData.Result = result
	getSliFinishedEventData.Message = message

	// Convert To CloudEvent
	finishedEvent := cloudevents.NewEvent()
	finishedEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName))
	finishedEvent.SetData(cloudevents.ApplicationJSON, getSliFinishedEventData)

	return SendEvent(myKeptn, finishedEvent, incomingEvent)
}

// HandleGetSliTriggeredEvent handles get-sli.triggered events if SLIProvider == keptn-service-template-go
// This function acts as an example showing how to handle get-sli events by sending .started and .finished events
// TODO: adapt handler code to your needs
func HandleGetSliTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.GetSLITriggeredEventData) error {
	log.Printf("Handling get-sli.triggered Event: %s", incomingEvent.Context.GetID())

	// Step 1 - Do we need to do something?
	// Lets make sure we are only processing an event that really belongs to our SLI Provider
	if data.GetSLI.SLIProvider != "keptn-service-template-go" {
		log.Printf("Not handling get-sli event as it is meant for %s", data.GetSLI.SLIProvider)
		return nil
	}

	// Step 2 - Send out a get-sli.started CloudEvent
	// The get-sli.started cloud-event is new since Keptn 0.8.0 and is required for the task to start
	// ToDo: This will need to be refactored once https://github.com/keptn/keptn/issues/2913 is ready
	getSliStartedData := keptnv2.GetSLIStartedEventData{}

	getSliStartedData.EventData = data.EventData
	getSliStartedData.Status = keptnv2.StatusSucceeded // alternative: keptnv2.StatusErrored
	getSliStartedData.Result = keptnv2.ResultPass      // alternative: keptnv2.ResultFailed

	// Convert To CloudEvent
	startedEvent := cloudevents.NewEvent()
	startedEvent.SetType(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName))
	startedEvent.SetData(cloudevents.ApplicationJSON, getSliStartedData)

	// send action.started event
	SendEvent(myKeptn, startedEvent, incomingEvent)

	// Step 4 - prep-work
	// Get any additional input / configuration data
	// - Labels: get the incoming labels for potential config data and use it to pass more labels on result, e.g: links
	// - SLI.yaml: if your service uses SLI.yaml to store query definitions for SLIs get that file from Keptn
	labels := data.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	testRunID := labels["testRunId"]

	// Step 5 - get SLI Config File
	// Get SLI File from keptn-service-template-go subdirectory of the config repo - to add the file use:
	//   keptn add-resource --project=PROJECT --stage=STAGE --service=SERVICE --resource=my-sli-config.yaml  --resourceUri=keptn-service-template-go/sli.yaml
	sliFile := "keptn-service-template-go/sli.yaml"
	sliConfigFileContent, err := myKeptn.GetKeptnResource(sliFile)

	// FYI you do not need to "fail" if sli.yaml is missing, you can also assume smart defaults like we do
	// in keptn-contrib/dynatrace-service and keptn-contrib/prometheus-service
	if err != nil {
		// failed to fetch sli config file
		errMsg := fmt.Sprintf("Failed to fetch SLI file %s from config repo: %s", sliFile, err.Error())
		log.Println(errMsg)
		// send a get-sli.finished event with status=error and result=failed back to Keptn
		return sendGetSliFinishedCloudEvent(myKeptn, incomingEvent, data, keptnv2.StatusErrored, keptnv2.ResultFailed, errMsg)
	}

	fmt.Println(sliConfigFileContent)

	// Step 6 - do your work - iterate through the list of requested indicators and return their values
	// Indicators: this is the list of indicators as requested in the SLO.yaml
	// SLIResult: this is the array that will receive the results
	indicators := data.GetSLI.Indicators
	sliResults := []*keptnv2.SLIResult{}

	for _, indicatorName := range indicators {
		sliResult := &keptnv2.SLIResult{
			Metric: indicatorName,
			Value:  123.4, // ToDo: Fetch the values from your monitoring tool here
		}
		sliResults = append(sliResults, sliResult)
	}

	// Step 7 - add additional context via labels (e.g., a backlink to the monitoring or CI tool)
	labels["Link to Data Source"] = "https://mydatasource/myquery?testRun=" + testRunID

	// Step 8 - Build get-sli.finished event data
	// ToDo: This will need to be refactored once https://github.com/keptn/keptn/issues/2913 is ready
	getSliFinishedEventData := keptnv2.GetSLIFinishedEventData{
		EventData: keptnv2.EventData{
			Project: data.Project,
			Stage:   data.Stage,
			Service: data.Service,
			Labels:  labels,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
		GetSLI: struct {
			Start           string               `json:"start"`
			End             string               `json:"end"`
			IndicatorValues []*keptnv2.SLIResult `json:"indicatorValues"`
		}{
			IndicatorValues: sliResults,
			Start:           data.GetSLI.Start,
			End:             data.GetSLI.End,
		},
	}

	// Step 9 - Convert To CloudEvent
	finishedEvent := cloudevents.NewEvent()
	finishedEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName))
	finishedEvent.SetData(cloudevents.ApplicationJSON, getSliFinishedEventData)

	// Step 10 - send action.finished CloudEvent back to Keptn
	return SendEvent(myKeptn, finishedEvent, incomingEvent)
}

// HandleProblemEvent handles two problem events:
// - ProblemOpenEventType = "sh.keptn.event.problem.open"
// - ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
func HandleProblemEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

// HandleActionTriggeredEvent handles action.triggered events
// TODO: add in your handler code
func HandleActionTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())
	log.Printf("Action=%s\n", data.Action.Action)

	// check if action is supported
	if data.Action.Action == "action-xyz" {
		// -----------------------------------------------------
		// 1. Send Action.Started Cloud-Event
		// -----------------------------------------------------

		// generate an action.started event
		actionStartedData := keptnv2.ActionStartedEventData{}

		actionStartedData.EventData = data.EventData
		actionStartedData.Status = keptnv2.StatusSucceeded // alternative: keptnv2.StatusErrored
		actionStartedData.Result = keptnv2.ResultPass      // alternative: keptnv2.ResultFailed

		// Convert To CloudEvent
		startedEvent := cloudevents.NewEvent()
		startedEvent.SetType(keptnv2.GetStartedEventType(keptnv2.ActionTaskName))
		startedEvent.SetData(cloudevents.ApplicationJSON, actionStartedData)

		// send action.started event
		SendEvent(myKeptn, startedEvent, incomingEvent)

		// -----------------------------------------------------
		// 2. Implement your remediation action here
		// -----------------------------------------------------
		time.Sleep(5 * time.Second) // Example: Wait 5 seconds. Maybe the problem fixes itself.

		// -----------------------------------------------------
		// 3. Send Action.Finished Cloud-Event
		// -----------------------------------------------------

		// generate an action.finished event
		actionFinishedData := keptnv2.ActionFinishedEventData{}

		actionFinishedData.EventData = data.EventData
		actionFinishedData.Status = keptnv2.StatusSucceeded // alternative: keptnv2.StatusErrored
		actionFinishedData.Result = keptnv2.ResultPass      // alternative: keptnv2.ResultFailed
		actionFinishedData.Message = "Successfully sleeped!"

		// Convert To CloudEvent
		finishedEvent := cloudevents.NewEvent()
		finishedEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.ActionTaskName))
		finishedEvent.SetData(cloudevents.ApplicationJSON, actionFinishedData)

		// send action.finished event
		SendEvent(myKeptn, finishedEvent, incomingEvent)

	} else {
		log.Printf("Retrieved unknown action %s, skipping...", data.Action.Action)
		return nil
	}
	return nil
}
