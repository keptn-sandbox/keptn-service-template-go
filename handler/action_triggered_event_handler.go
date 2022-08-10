package handler

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk"
	"time"
)

type ActionTriggeredEventHandler struct {
}

func NewActionTriggeredEventHandler() *ActionTriggeredEventHandler {
	return &ActionTriggeredEventHandler{}
}

// Execute handles action.triggered events
// TODO: Add in your handler code
func (g *ActionTriggeredEventHandler) Execute(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	k.Logger().Infof("Handling Action Triggered Event: %s", event.ID)
	actionTriggeredEvent := &keptnv2.ActionTriggeredEventData{}

	if err := keptnv2.Decode(event.Data, actionTriggeredEvent); err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "failed to decode action.triggered event: " + err.Error()}
	}

	k.Logger().Infof("Action=%s", actionTriggeredEvent.Action.Action)

	// check if action is supported
	if actionTriggeredEvent.Action.Action == "action-xyz" {
		k.Logger().Info("Action remediation triggered")
		// -----------------------------------------------------
		// TODO: Implement your remediation action here
		// -----------------------------------------------------
		time.Sleep(1 * time.Second) // Example: Wait 5 seconds. Maybe the problem fixes itself.

		// Return finished event
		finishedEventData := getActionFinishedEvent(keptnv2.ResultPass, keptnv2.StatusSucceeded, *actionTriggeredEvent, "")

		return finishedEventData, nil
	}

	k.Logger().Infof("Retrieved unknown action %s, skipping...", actionTriggeredEvent.Action.Action)
	return nil, nil
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
