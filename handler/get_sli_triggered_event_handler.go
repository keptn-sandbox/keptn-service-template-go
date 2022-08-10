package handler

import (
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk"
)

type GetSliEventHandler struct {
}

func NewGetSliEventHandler() *GetSliEventHandler {
	return &GetSliEventHandler{}
}

// Execute handles get-sli.triggered events if SLIProvider == keptn-service-template-go
// This function acts as an example showing how to handle get-sli events
// TODO: Adapt handler code to your needs
func (g *GetSliEventHandler) Execute(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
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
		k.Logger().Infof("Error while fetching SLI file: %e", err)
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "error while fetching SLI file: " + err.Error()}
	}

	k.Logger().Infof("SLI config content: %s", sliConfigFileContent)

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
