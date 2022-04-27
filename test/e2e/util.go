package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"testing"
	"time"
)

// KeptnConnectionDetails contains the endpoint and the API token for Keptn
type KeptnConnectionDetails struct {
	Endpoint string
	APIToken string
}

// readKeptnConnectionDetailsFromEnv parses the environment variables and creates a KeptnConnectionDetails
func readKeptnConnectionDetailsFromEnv() KeptnConnectionDetails {
	return KeptnConnectionDetails{
		Endpoint: os.Getenv("KEPTN_ENDPOINT"),
		APIToken: os.Getenv("KEPTN_API_TOKEN"),
	}
}

// isE2ETestingAllowed checks if the E2E tests are allowed to run by parsing environment variables
func isE2ETestingAllowed() bool {
	boolean, err := strconv.ParseBool(os.Getenv("ENABLE_E2E_TEST"))
	if err != nil {
		return false
	}

	return boolean
}

// convertKeptnModelToErrorString transforms the models.Error structure to an error string
func convertKeptnModelToErrorString(keptnError *models.Error) string {
	if keptnError == nil {
		return ""
	}

	if keptnError.Message != nil {
		return fmt.Sprintf("%d, %s", keptnError.Code, *keptnError.Message)
	}

	return fmt.Sprintf("%d <no error message>", keptnError.Code)
}

// readKeptnContextExtendedCE reads a file from a given path and returnes the parsed models.KeptnContextExtendedCE struct
func readKeptnContextExtendedCE(path string) (*models.KeptnContextExtendedCE, error) {
	fileContents, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	var keptnContextExtendedCE models.KeptnContextExtendedCE
	err = json.Unmarshal(fileContents, &keptnContextExtendedCE)

	if err != nil {
		return nil, fmt.Errorf("unable to parse event: %w", err)
	}

	return &keptnContextExtendedCE, nil
}

// eventData structure contains common fields in the data part of a  models.KeptnContextExtendedCE struct that are needed by E2E tests
type eventData struct {
	Message string `mapstruct:"message,omitempty"`
	Project string `mapstruct:"project,omitempty"`
	Result  string `mapstruct:"result,omitempty"`
	Service string `mapstruct:"service,omitempty"`
	Stage   string `mapstruct:"stage,omitempty"`
	Status  string `mapstruct:"status,omitempty"`
}

// parseKeptnEventData parse the Data field of the models.KeptnContextExtendedCE structure into a form, which is more
// convenient to work with
func parseKeptnEventData(ce *models.KeptnContextExtendedCE) (*eventData, error) {
	var eventData eventData
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &eventData,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create mapstructure decoder: %w", err)
	}

	err = decoder.Decode(ce.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to decode event data: %w", err)
	}

	return &eventData, nil
}

// createK8sSecret creates a k8s secret from a json file and uploads it into the give namespace
func createK8sSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace string, jsonFilePath string) (func(ctx2 context.Context), error) {

	// read the file from the given path
	file, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets file: %w", err)
	}

	// unmarshal the contents from the file, since we are using k8s classes it must be a json
	var secret v1.Secret
	err = json.Unmarshal(file, &secret)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal secrets json: %s", err)
	}

	// create the secret in k8s
	_, err = clientset.CoreV1().Secrets(namespace).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s secret: %w", err)
	}

	// return a function which can be used to delete the secret after the tests have finished
	return func(ctx2 context.Context) {
		err := clientset.CoreV1().Secrets(namespace).Delete(ctx2, secret.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Errorf("Unable to delete secret!")
		}
	}, nil
}

// testEnvironment structure holds different structures and information that are commonly used by the E2E test environment
type testEnvironment struct {
	K8s         *kubernetes.Clientset
	API         KeptnAPI
	EventData   *eventData
	Event       *models.KeptnContextExtendedCE
	Namespace   string
	CleanupFunc func()
}

// setupE2ETTestEnvironment creates the basic e2e test environment, which includes creating a service and a project in Keptn,
// additionally also the given job configuration is uploaded to Keptn such that simple E2E tests can continue by sending
// the desired events or continue to customize the project
func setupE2ETTestEnvironment(t *testing.T, eventJSONFilePath string, shipyardPath string) testEnvironment {
	// Just test if we can connect to the cluster
	clientset, err := keptnutils.GetClientset(false)
	require.NoError(t, err)
	assert.NotNil(t, clientset)

	// Create a new Keptn api for the use of the E2E test
	keptnAPI := NewKeptAPI(readKeptnConnectionDetailsFromEnv())

	// Read the event we want to trigger and extract the project, service and stage
	keptnEvent, err := readKeptnContextExtendedCE(eventJSONFilePath)
	require.NoError(t, err)

	eventData, err := parseKeptnEventData(keptnEvent)
	require.NoError(t, err)

	// Load shipyard file and create the project in Keptn
	shipyardFile, err := ioutil.ReadFile(shipyardPath)
	require.NoError(t, err)

	err = keptnAPI.CreateProject(eventData.Project, shipyardFile)
	require.NoError(t, err)

	// deferred function must be called by the caller
	deleteProjectFunc := func() {
		if err := keptnAPI.DeleteProject(eventData.Project); err != nil {
			t.Log(err.Error())
		}
	}

	// Create a service in Keptn
	err = keptnAPI.CreateService(eventData.Project, eventData.Service)
	require.NoError(t, err)

	return testEnvironment{
		K8s:         clientset,
		API:         keptnAPI,
		EventData:   eventData,
		Event:       keptnEvent,
		Namespace:   "keptn",
		CleanupFunc: deleteProjectFunc,
	}
}

// requireWaitForEvent checks if an event occurred in a specific time frame while polling the event bus of keptn, the eventValidator
// should return true if the desired event was found
func requireWaitForEvent(t *testing.T, api KeptnAPI, waitFor time.Duration, tick time.Duration, keptnContext *models.EventContext, eventType string, eventValidator func(c *models.KeptnContextExtendedCE) bool, source string) {
	checkForEventsToMatch := func() bool {
		events, err := api.GetEvents(keptnContext.KeptnContext)
		require.NoError(t, err)

		// for each event we have to check if the type is the correct one and if
		// the source of the event matches the job executor, if that is the case
		// the event can be checked by the eventValidator
		for _, event := range events {
			if *event.Type == eventType && *event.Source == source {
				if eventValidator(event) {
					return true
				}
			}
		}

		return false
	}

	// We require waiting for a keptn event, this is useful to exit out tests if no .started event occurred.
	// It doesn't make sense in these cases to wait for a .finished or other .triggered events ...
	require.Eventuallyf(t, checkForEventsToMatch, waitFor, tick, "did not receive keptn event: %s", eventType)
}

// CreateTmpShipyardFile creates a temporary shipyard file from the provided YAML content and returns the name of the file
func CreateTmpShipyardFile(shipyardContent string) (string, error) {
	return CreateTmpFile("shipyard-*.yaml", shipyardContent)
}

// CreateTmpFile creates a temporary file using the provided file content
func CreateTmpFile(fileNamePattern, fileContent string) (string, error) {
	file, err := ioutil.TempFile("", fileNamePattern)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file.Name(), []byte(fileContent), os.ModeAppend); err != nil {
		err = os.Remove(file.Name())
		if err != nil {
			return "", err
		}
		return "", err
	}
	return file.Name(), nil
}

// CreateTmpDir creates a temporary directory on the file system
func CreateTmpDir() (string, error) {
	return ioutil.TempDir("", "")
}
