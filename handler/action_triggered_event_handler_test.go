package handler

import (
	"encoding/json"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newEvent(filename string) keptnapi.KeptnContextExtendedCE {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := keptnapi.KeptnContextExtendedCE{}
	err = json.Unmarshal(content, &event)
	_ = err
	return event
}

func Test_Receiving_GetActionTriggeredEvent(t *testing.T) {
	var returnedStatusCode = 200
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if strings.Contains(r.URL.String(), "/admin/features/") {
				w.WriteHeader(returnedStatusCode)
				w.Write([]byte(`{}`))
				return
			}

			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			keptnCE := &keptnapi.KeptnContextExtendedCE{}

			_ = json.Unmarshal(body, keptnCE)

			w.WriteHeader(returnedStatusCode)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	fakeKeptn := sdk.NewFakeKeptn("test-service-template-svc")
	fakeKeptn.AddTaskHandler("sh.keptn.event.action.triggered", NewActionTriggeredEventHandler())

	fakeKeptn.NewEvent(newEvent("../test/events/action_triggered.json"))

	fakeKeptn.AssertNumberOfEventSent(t, 2)

	fakeKeptn.AssertSentEventType(t, 0, keptnv2.GetStartedEventType("action"))
	fakeKeptn.AssertSentEventType(t, 1, keptnv2.GetFinishedEventType("action"))

	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)
}
