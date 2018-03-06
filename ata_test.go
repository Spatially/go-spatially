package spatially

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/Spatially/go-geometry"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestNewATAHome(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	api, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	pointWKT := "POINT(-71.064156780428 42.35862883483673)"
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/ads/science/ata", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		var request ataRequest
		j, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(j, &request); err != nil {
			return nil, err
		}
		if request.AreaType != "ATA" {
			t.Error("Invalid Area Type")
		}
		if request.Buffer != 100 {
			t.Error("Invalid buffer value")
		}
		if len(request.LocationType) == 0 || request.LocationType[0] != "Home" {
			t.Error("Invalid location type")
		}
		if request.PointWKT != pointWKT {
			t.Error("Invalid point wkt")
		}
		if request.TimeOfDay != AllDay.String() {
			t.Error("Invalid Time of day")
		}
		response := ataResponse{
			Buffer:            100,
			FeatureCollection: &geometry.FeatureCollection{},
		}
		return httpmock.NewJsonResponse(200, response)
	})
	ata, err := NewATA(api, pointWKT, nil)
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", *ata)
}

func TestNewATAWork(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	api, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	pointWKT := "POINT(-71.064156780428 42.35862883483673)"
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/ads/science/ata", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		var request ataRequest
		j, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(j, &request); err != nil {
			return nil, err
		}
		if request.AreaType != "ATA" {
			t.Error("Invalid Area Type")
		}
		if request.Buffer != 100 {
			t.Error("Invalid buffer value")
		}
		if len(request.LocationType) == 0 || request.LocationType[0] != "Work" {
			t.Error("Invalid location type")
		}
		if request.PointWKT != pointWKT {
			t.Error("Invalid point wkt")
		}
		if request.TimeOfDay != AllDay.String() {
			t.Error("Invalid Time of day")
		}
		response := ataResponse{
			Buffer:            100,
			FeatureCollection: &geometry.FeatureCollection{},
		}
		return httpmock.NewJsonResponse(200, response)
	})
	ata, err := NewATA(api, pointWKT, &ATAOptions{
		LocationType: Work,
	})
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", *ata)
}

func TestNewATAHomeAndWork(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	api, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	pointWKT := "POINT(-71.064156780428 42.35862883483673)"
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/ads/science/ata", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		var request ataRequest
		j, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(j, &request); err != nil {
			return nil, err
		}
		if request.AreaType != "ATA" {
			t.Error("Invalid Area Type")
		}
		if request.Buffer != 100 {
			t.Error("Invalid buffer value")
		}
		if len(request.LocationType) < 2 || request.LocationType[0] != "Home" || request.LocationType[1] != "Work" {
			t.Error("Invalid location type")
		}
		if request.PointWKT != pointWKT {
			t.Error("Invalid point wkt")
		}
		if request.TimeOfDay != AllDay.String() {
			t.Error("Invalid Time of day")
		}
		response := ataResponse{
			Buffer:            100,
			FeatureCollection: &geometry.FeatureCollection{},
		}
		return httpmock.NewJsonResponse(200, response)
	})
	ata, err := NewATA(api, pointWKT, &ATAOptions{
		LocationType: HomeAndWork,
	})
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", *ata)
}

func TestNewATATimeOfDay(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	api, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	pointWKT := "POINT(-71.064156780428 42.35862883483673)"
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/ads/science/ata", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		var request ataRequest
		j, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(j, &request); err != nil {
			return nil, err
		}
		if request.AreaType != "ATA" {
			t.Error("Invalid Area Type")
		}
		if request.Buffer != 100 {
			t.Error("Invalid buffer value")
		}
		if len(request.LocationType) == 0 || request.LocationType[0] != "Home" {
			t.Error("Invalid location type")
		}
		if request.PointWKT != pointWKT {
			t.Error("Invalid point wkt")
		}
		if request.TimeOfDay != Morning.String() {
			t.Error("Invalid Time of day")
		}
		response := ataResponse{
			Buffer:            100,
			FeatureCollection: &geometry.FeatureCollection{},
		}
		return httpmock.NewJsonResponse(200, response)
	})
	ata, err := NewATA(api, pointWKT, &ATAOptions{
		TimeOfDay: Morning,
	})
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", *ata)
}

func TestNewATAGeofence(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	api, err := NewAPI(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	pointWKT := "POINT(-71.064156780428 42.35862883483673)"
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/ads/science/ata", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		var request ataRequest
		j, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(j, &request); err != nil {
			return nil, err
		}
		if request.AreaType != "ATA" {
			t.Error("Invalid Area Type")
		}
		if request.Buffer != 100 {
			t.Error("Invalid buffer value")
		}
		if len(request.LocationType) == 0 || request.LocationType[0] != "Home" {
			t.Error("Invalid location type")
		}
		if request.PointWKT != pointWKT {
			t.Error("Invalid point wkt")
		}
		if request.TimeOfDay != AllDay.String() {
			t.Error("Invalid Time of day")
		}
		if request.GeoFence != true {
			t.Error("Invalid Geo fence")
		}
		response := ataResponse{
			Buffer:            100,
			FeatureCollection: &geometry.FeatureCollection{},
		}
		return httpmock.NewJsonResponse(200, response)
	})
	ata, err := NewATA(api, pointWKT, &ATAOptions{
		GeoFence: true,
	})
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", *ata)
}
