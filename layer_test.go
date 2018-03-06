package spatially

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pborman/uuid"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestCreateLayer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/spatialdb/layer", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var request createLayerRequest
		if err := json.Unmarshal(body, &request); err != nil {
			return nil, err
		}
		if request.Name != "layer1" {
			t.Error("Invalid layer name in request")
		}
		layer := NewLayer()
		layer.ID = uuid.NewUUID().String()
		layer.Name = request.Name
		return httpmock.NewJsonResponse(200, layer)
	})
	layer := NewLayer()
	if err := layer.Create(sdb, "layer1"); err != nil {
		t.Error(err)
	}
	if len(layer.ID) == 0 {
		t.Error("Invalid layer id")
	}
	if layer.Name != "layer1" {
		t.Error("Invalid layer name")
	}
}

func TestGetLayers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterResponder("GET", SpatiallyAPI+"/spatialdb/layers", func(req *http.Request) (*http.Response, error) {
		layers := NewLayers()
		layers = append(layers, &Layer{
			ID:   uuid.NewUUID().String(),
			Name: "layer1",
		})
		return httpmock.NewJsonResponse(200, layers)
	})
	layers := NewLayers()
	if err := layers.Get(sdb); err != nil {
		t.Error(err)
	}
	if len(layers) == 0 {
		t.Error("Layer count should not be zero")
	}
	layer := layers[0]
	if layer.Name != "layer1" {
		t.Error("First layer's name should be 'layer1'")
	}
}

func TestGetLayer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	layerID := uuid.NewUUID().String()
	httpmock.RegisterResponder("GET", SpatiallyAPI+"/spatialdb/layer/"+layerID, func(req *http.Request) (*http.Response, error) {
		layer := NewLayer()
		layer.ID = layerID
		layer.Name = "layer1"
		return httpmock.NewJsonResponse(200, layer)
	})
	layer := NewLayer()
	if err := layer.Get(sdb, layerID); err != nil {
		t.Error(err)
	}
	if layer.ID != layerID {
		t.Error("Invalid layer id")
	}
	if layer.Name != "layer1" {
		t.Error("Invalid layer name")
	}
}

func TestDeleteLayer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	layerID := uuid.NewUUID().String()
	httpmock.RegisterResponder("DELETE", SpatiallyAPI+"/spatialdb/layer/"+layerID, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, ""), nil
	})
	layer := NewLayer()
	if err := layer.Delete(sdb, layerID); err != nil {
		t.Error(err)
	}
}
