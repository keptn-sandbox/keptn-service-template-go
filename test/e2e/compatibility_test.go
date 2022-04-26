package e2e

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const shipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "user_managed"`

const actionTriggeredEvent = `
{
  "type": "sh.keptn.event.action.triggered",
  "specversion": "1.0",
  "source": "test-events",
  "id": "f2b878d3-03c0-4e8f-bc3f-454bc1b3d79b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext": "08735340-6f9e-4b32-97ff-3b6c292bc50i",
  "data": {
    "project": "user-managed",
    "stage": "dev",
    "service": "nginx",
    "labels": {
      "testId": "4711",
      "buildId": "build-17",
      "owner": "JohnDoe"
    },
    "status": "succeeded",
    "result": "pass",
    "action": {
      "name": "action-xyz",
      "action": "action-xyz",
      "description": "action-xyz",
      "value": "1"
    },
    "problem": {
    }
  }
}`

const sliTriggeredEvent = `
{
    "data": {
      "get-sli": {
        "customFilters": [],
        "end": "2021-01-15T15:09:45.000Z",
        "indicators": [
          "response_time_p95",
          "some_other_metric"
        ],
        "sliProvider": "keptn-service-template-go",
        "start": "2021-01-15T15:04:45.000Z"
      },
      "labels": null,
      "message": "",
      "project": "user-managed",
      "result": "",
      "service": "nginx",
      "stage": "dev",
      "status": ""
    },
    "id": "409539ae-c0b9-436e-abc6-c257292e28ff",
    "source": "test-events",
    "specversion": "1.0",
    "time": "2021-01-15T15:09:46.144Z",
    "type": "sh.keptn.event.get-sli.triggered",
    "shkeptncontext": "da7aec34-78c4-4182-a2c8-51eb88f5871d"
}`

func Test_ActionTriggered(t *testing.T) {
	if !isE2ETestingAllowed() {
		t.Skip("Skipping Test_ActionTriggered, not allowed by environment")
	}

	shipyardFilePath, err := CreateTmpShipyardFile(shipyard)

	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	actionEventFilePath, err := CreateTmpFile("event.yaml", actionTriggeredEvent)

	require.Nil(t, err)
	defer os.Remove(actionEventFilePath)

	// Setup the E2E test environment
	testEnv := setupE2ETTestEnvironment(t,
		actionEventFilePath,
		shipyardFilePath,
	)

	// Make sure project is delete after the tests are completed
	defer testEnv.CleanupFunc()

	// Send the event to keptn
	keptnContext, err := testEnv.API.SendEvent(testEnv.Event)
	require.NoError(t, err)

	// Checking if the service-template-go responded with a .started event
	requireWaitForEvent(t,
		testEnv.API,
		1*time.Minute,
		1*time.Second,
		keptnContext,
		"sh.keptn.event.action.started",
		func(_ *models.KeptnContextExtendedCE) bool {
			return true
		},
		"keptn-service-template-go",
	)

	// Checking if the service-template-go responded with a .finished event
	requireWaitForEvent(t,
		testEnv.API,
		1*time.Minute,
		1*time.Second,
		keptnContext,
		"sh.keptn.event.action.finished",
		func(_ *models.KeptnContextExtendedCE) bool {
			return true
		},
		"keptn-service-template-go",
	)
}

func Test_SLITriggered(t *testing.T) {
	if !isE2ETestingAllowed() {
		t.Skip("Skipping Test_SLITriggered, not allowed by environment")
	}

	shipyardFilePath, err := CreateTmpShipyardFile(shipyard)

	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	sliEventFilePath, err := CreateTmpFile("event.yaml", sliTriggeredEvent)

	require.Nil(t, err)
	defer os.Remove(sliEventFilePath)

	// Setup the E2E test environment
	testEnv := setupE2ETTestEnvironment(t,
		sliEventFilePath,
		shipyardFilePath,
	)

	// Make sure project is delete after the tests are completed
	defer testEnv.CleanupFunc()

	// Send the event to keptn
	keptnContext, err := testEnv.API.SendEvent(testEnv.Event)
	require.NoError(t, err)

	// Checking if the service-template-go responded with a .started event
	requireWaitForEvent(t,
		testEnv.API,
		1*time.Minute,
		1*time.Second,
		keptnContext,
		"sh.keptn.event.get-sli.started",
		func(_ *models.KeptnContextExtendedCE) bool {
			return true
		},
		"keptn-service-template-go",
	)

	// Checking if the service-template-go responded with a .finished event
	requireWaitForEvent(t,
		testEnv.API,
		1*time.Minute,
		1*time.Second,
		keptnContext,
		"sh.keptn.event.get-sli.finished",
		func(_ *models.KeptnContextExtendedCE) bool {
			return true
		},
		"keptn-service-template-go",
	)
}
