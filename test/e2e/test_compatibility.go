package e2e

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func Test_ActionTriggeredTest(t *testing.T) {
	projectName := "user-managed"
	serviceName := "nginx"
	stageName := "dev"
	sequenceName := "delivery"
	serviceChartPath := "https://charts.bitnami.com/bitnami/nginx-8.9.0.tgz"

	shipyardFilePath, err := CreateTmpShipyardFile(shipyard)

	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("wget %s -O chart.tgz", serviceChartPath))
	require.Nil(t, err)

	defer os.Remove("chart.tgz")

	// make sure the namespace from a previous test run has been deleted properly
	exists, err := keptnkubeutils.ExistsNamespace(false, projectName+"-dev")
	if exists {
		t.Logf("Deleting namespace %s-dev from previous test execution", projectName)
		clientset, err := keptnkubeutils.GetClientset(false)
		require.Nil(t, err)
		err = clientset.CoreV1().Namespaces().Delete(context.TODO(), projectName+"-dev", v1.DeleteOptions{})
		require.Nil(t, err)
	}

	require.Eventually(t, func() bool {
		t.Logf("Checking if namespace %s-dev is still there", projectName)
		exists, err := keptnkubeutils.ExistsNamespace(false, projectName+"-dev")
		if err != nil || exists {
			t.Logf("Namespace %s-dev is still there", projectName)
			return false
		}
		t.Logf("Namespace %s-dev has been removed - proceeding with test execution", projectName)
		return true
	}, 60*time.Second, 5*time.Second)

	// check if the project is already available - if not, delete it before creating it again
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// create the service
	t.Logf("Creating service %s in project %s", serviceName, projectName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))
	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// upload the service's helm chart
	t.Logf("Uploading the helm chart of service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --service=%s --project=%s --all-stages --resource=./chart.tgz --resourceUri=helm/%s.tgz", serviceName, projectName, serviceName))
	require.Nil(t, err)

	// trigger the sequence without defining custom endpoints first
	t.Logf("Triggering the first delivery sequence without providing custom endpoints")
	keptnContextID, err := TriggerSequence(projectName, serviceName, stageName, sequenceName, nil)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContextID)

	// wait until we get a action.triggered event
	var actionTriggeredEvent *models.KeptnContextExtendedCE
	t.Log("Waiting for deployment to complete")
	require.Eventually(t, func() bool {
		actionTriggeredEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetStartedEventType(keptnv2.ActionTaskName))
		if err != nil || actionTriggeredEvent == nil {
			t.Log("Action has not been started yet... Waiting a couple of seconds before checking again")
			return false
		}
		return true
	}, 60*time.Second, 5*time.Second)
	t.Log("Deployment has been completed")
}
