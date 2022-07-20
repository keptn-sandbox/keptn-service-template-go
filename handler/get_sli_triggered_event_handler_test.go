package handler

import (
	"encoding/json"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Receiving_GetSliTriggeredEvent(t *testing.T) {
	ch := make(chan *keptnapi.KeptnContextExtendedCE)

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
			go func() { ch <- keptnCE }()
		}),
	)
	defer ts.Close()

	fakeKeptn := sdk.NewFakeKeptn("test-service-template-svc")
	fakeKeptn.AddTaskHandler("sh.keptn.event.get-sli.triggered", NewGetSliEventHandler())

	fakeKeptn.NewEvent(newEvent("../test/events/get_sli_triggered.json"))

	fakeKeptn.AssertNumberOfEventSent(t, 2)

	fakeKeptn.AssertSentEventType(t, 0, keptnv2.GetStartedEventType("get-sli"))
	fakeKeptn.AssertSentEventType(t, 1, keptnv2.GetFinishedEventType("get-sli"))

	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)
}
