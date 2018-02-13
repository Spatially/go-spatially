package spatiallydb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

const applicationCode = "9e2425a9-084f-11e8-9040-acde48001122"
const applicationKey = "9e24259e-084f-11e8-9040-acde48001122"

func TestNew(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	_, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
}

func mockGatewayEndpoint(t *testing.T) {
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/gateway/client", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var request gatewayRequest
		if err := json.Unmarshal(body, &request); err != nil {
			return nil, err
		}
		if request.Key != applicationKey {
			t.Error("Request Key:", request.Key, "does not equal the application key", applicationKey)
		}
		if request.Code != applicationCode {
			t.Error("Request Code:", request.Code, "does not equal the application code", applicationCode)
		}
		return httpmock.NewJsonResponse(200, gatewayResponse{"authToken"})
	})
}
