package main

import (
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk"
	"time"
)

type EventHandler struct {
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

func (g *EventHandler) Execute(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	// Check the event type and handle the event accordingly. If you only listen to one event this can be skipped!
	if *event.Type == keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName) {
		finishedEventData, err := handleGetSliTriggeredEvent(k, event)

		if err != nil {
			return nil, err
		}
		return finishedEventData, nil
	} else if *event.Type == keptnv2.GetTriggeredEventType(keptnv2.GetActionTaskName) {
		finishedEventData, err := HandleActionTriggeredEvent(k, event)

		if err != nil {
			return nil, err
		}
		return finishedEventData, nil
	}

	return nil, nil
}

// handleGetSliTriggeredEvent handles get-sli.triggered events if SLIProvider == keptn-service-template-go
// This function acts as an example showing how to handle get-sli events
// TODO: Adapt handler code to your needs
func handleGetSliTriggeredEvent(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	k.Logger().Infof("Handling get-sli.triggered Event: %s", event.ID)

	sliTriggeredEvent := &keptnv2.GetSLITriggeredEventData{}

	if err := keptnv2.Decode(event.Data, sliTriggeredEvent); err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "failed to decode sli.triggered event: " + err.Error()}
	}

	// Check if the event belongs to our SLI Provider
	if sliTriggeredEvent.GetSLI.SLIProvider != "keptn-service-template-go" {
		k.Logger().Infof("Not handling get-sli event as it is meant for %s", sliTriggeredEvent.GetSLI.SLIProvider)
		return nil, nil
	}

	// Get SLI File from keptn-service-template-go subdirectory of the config repo - to add the file use:
	sliFile := "keptn-service-template-go/sli.yaml"
	resourceScope := *api.NewResourceScope().Project(sliTriggeredEvent.Project).Stage(sliTriggeredEvent.Stage).Service(sliTriggeredEvent.Service).Resource(sliFile)
	sliConfigFileContent, err := k.GetResourceHandler().GetResource(resourceScope)

	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "error while fetching SLI file: " + err.Error()}
	}

	k.Logger().Info(sliConfigFileContent)

	// TODO: Implement your functionality here
	indicators := sliTriggeredEvent.GetSLI.Indicators
	var sliResults []*keptnv2.SLIResult

	for _, indicatorName := range indicators {
		sliResult := &keptnv2.SLIResult{
			Metric: indicatorName,
			Value:  123.4, // TODO: Fetch the values from your monitoring tool here
		}
		sliResults = append(sliResults, sliResult)
	}

	finishedEventData := getSliFinishedEvent(keptnv2.ResultPass, keptnv2.StatusSucceeded, *sliTriggeredEvent, "", sliResults)

	return finishedEventData, nil
}

// HandleActionTriggeredEvent handles action.triggered events
// TODO: Add in your handler code
func HandleActionTriggeredEvent(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	actionTriggeredEvent := &keptnv2.ActionTriggeredEventData{}

	if err := keptnv2.Decode(event.Data, actionTriggeredEvent); err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "failed to decode action.triggered event: " + err.Error()}
	}

	k.Logger().Infof("Handling Action Triggered Event: %s", event.ID)
	k.Logger().Infof("Action=%s\n", actionTriggeredEvent.Action.Action)

	// check if action is supported
	if actionTriggeredEvent.Action.Action == "action-xyz" {
		// -----------------------------------------------------
		// TODO: Implement your remediation action here
		// -----------------------------------------------------
		time.Sleep(5 * time.Second) // Example: Wait 5 seconds. Maybe the problem fixes itself.

		// Return finished event
		finishedEventData := getActionFinishedEvent(keptnv2.ResultPass, keptnv2.StatusSucceeded, *actionTriggeredEvent, "")

		return finishedEventData, nil
	} else {
		k.Logger().Infof("Retrieved unknown action %s, skipping...", actionTriggeredEvent.Action.Action)
		return nil, nil
	}
	return nil, nil
}

func getSliFinishedEvent(result keptnv2.ResultType, status keptnv2.StatusType, sliTriggeredEvent keptnv2.GetSLITriggeredEventData, message string, sliResult []*keptnv2.SLIResult) keptnv2.GetSLIFinishedEventData {

	return keptnv2.GetSLIFinishedEventData{
		EventData: keptnv2.EventData{
			Project: sliTriggeredEvent.Project,
			Stage:   sliTriggeredEvent.Stage,
			Service: sliTriggeredEvent.Service,
			Labels:  sliTriggeredEvent.Labels,
			Status:  status,
			Result:  result,
			Message: message,
		},
		GetSLI: keptnv2.GetSLIFinished{
			IndicatorValues: sliResult,
		},
	}
}

func getActionFinishedEvent(result keptnv2.ResultType, status keptnv2.StatusType, actionTriggeredEvent keptnv2.ActionTriggeredEventData, message string) keptnv2.ActionFinishedEventData {

	return keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Project: actionTriggeredEvent.Project,
			Stage:   actionTriggeredEvent.Stage,
			Service: actionTriggeredEvent.Service,
			Labels:  actionTriggeredEvent.Labels,
			Status:  status,
			Result:  result,
			Message: message,
		},
	}
}
