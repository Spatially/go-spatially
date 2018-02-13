package spatiallydb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	geojson "github.com/paulmach/go.geojson"
	"github.com/pborman/uuid"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestCreateFeaturePoint(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	feature := NewFeature()
	geometry := geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
	properties := map[string]interface{}{
		"name": "Starbucks",
	}
	layerID := uuid.NewUUID().String()
	featureID := uuid.NewUUID().String()
	mockCreateFeatureEndpoint(t, layerID, featureID)
	if err := feature.Create(sdb, layerID, geometry, properties); err != nil {
		t.Error(err)
	}
	if feature.PropertyMustString("name") != "Starbucks" {
		t.Error("Invalid feature name")
	}
	createdFeatureID, isString := feature.ID.(string)
	if !isString {
		t.Error("Invalid feature id")
	}
	if feature.Geometry.IsPoint() == false {
		t.Error("Invalid feature geometry")
	}
	if createdFeatureID != featureID {
		t.Error("Feature ID", createdFeatureID, "does not match", featureID)
	}
}

func mockCreateFeatureEndpoint(t *testing.T, layerID, featureID string) {
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/spatialdb/feature", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var request createFeatureRequest
		if err := json.Unmarshal(body, &request); err != nil {
			return nil, err
		}
		if request.LayerID != layerID {
			t.Error("Invalid layer id in request")
		}
		if request.Feature == nil {
			t.Error("Request feature is nil")
		}
		if request.Feature.Geometry == nil {
			t.Error("Request feature has a nil geometry")
		}
		if request.Feature.PropertyMustString("name") != "Starbucks" {
			t.Error("Request feature name does not match")
		}
		request.Feature.ID = featureID
		return httpmock.NewJsonResponse(200, request.Feature)
	})
}

func TestCreateFeaturePolygon(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	feature := NewFeature()
	geometry, err := WKTToGeometry("POLYGON((-71.06296062469482 42.362336359418954,-71.05918407440186 42.358277337975814,-71.06665134429932 42.35979950174449,-71.06296062469482 42.362336359418954))")
	if err != nil {
		t.Error(err)
	}
	properties := map[string]interface{}{
		"name": "Starbucks",
	}
	layerID := uuid.NewUUID().String()
	featureID := uuid.NewUUID().String()
	mockCreateFeatureEndpoint(t, layerID, featureID)
	if err := feature.Create(sdb, layerID, geometry, properties); err != nil {
		t.Error(err)
	}
	if feature.PropertyMustString("name") != "Starbucks" {
		t.Error("Invalid feature name")
	}
	createdFeatureID, isString := feature.ID.(string)
	if !isString {
		t.Error("Invalid feature id")
	}
	if feature.Geometry.IsPolygon() == false {
		t.Error("Invalid feature geometry")
	}
	if createdFeatureID != featureID {
		t.Error("Feature ID", createdFeatureID, "does not match", featureID)
	}
}

func TestGetFeaturesSimple(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	layerID := uuid.NewUUID().String()
	mockGetFeaturesEndpoint(t, layerID, nil)
	features := NewFeatures()
	if err := features.GetByLayer(sdb, layerID); err != nil {
		t.Error(err)
	}
	if len(features) == 0 {
		t.Error("Expected features length to be 1, it is", len(features))
	}
	feature := features[0]
	if feature.PropertyMustString("name") != "Starbucks" {
		t.Error("Expected feature name to be Starbucks")
	}
}

func mockGetFeaturesEndpoint(t *testing.T, layerID string, sp *SpatialConstraint) {
	httpmock.RegisterResponder("POST", SpatiallyAPI+"/spatialdb/features", func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var request getFeaturesRequest
		if err := json.Unmarshal(body, &request); err != nil {
			return nil, err
		}
		if request.LayerID != layerID {
			t.Error("Invalid request layer id")
		}
		if request.SpatialConstraint != nil {
			if sp.WKT != "" && request.SpatialConstraint.WKT != sp.WKT {
				t.Error("Invalid request spatial constraint")
			}
			if sp.Radius != 0.0 && request.SpatialConstraint.Radius != sp.Radius {
				t.Error("Invalid request spatial constraint radius")
			}
		}
		features := NewFeatures()
		feature := NewFeature()
		feature.Geometry = geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
		feature.Properties = map[string]interface{}{"name": "Starbucks"}
		features = append(features, feature)
		return httpmock.NewJsonResponse(200, features)
	})
}

func TestGetFeaturesSpatialConstraint(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	sp := &SpatialConstraint{
		WKT: "POLYGON((-71.06253147125244 42.34286335316219,-71.07102870941162 42.3436246263197,-71.07017040252686 42.346669626771906,-71.0627031326294 42.34559120597906,-71.06253147125244 42.34286335316219))",
	}
	layerID := uuid.NewUUID().String()
	mockGetFeaturesEndpoint(t, layerID, sp)
	features := NewFeatures()
	if err := features.GetBySpatialConstraint(sdb, layerID, sp); err != nil {
		t.Error(err)
	}
	if len(features) == 0 {
		t.Error("Expected features length to be 1, it is", len(features))
	}
	feature := features[0]
	if feature.PropertyMustString("name") != "Starbucks" {
		t.Error("Expected feature name to be Starbucks")
	}
}

func TestGetFeaturesSpatialConstraintBuffer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	sp := &SpatialConstraint{
		WKT:    "POINT (-71.06772422790527 42.35848049347556)",
		Radius: 1000.0,
	}
	layerID := uuid.NewUUID().String()
	mockGetFeaturesEndpoint(t, layerID, sp)
	features := NewFeatures()
	if err := features.GetBySpatialConstraint(sdb, layerID, sp); err != nil {
		t.Error(err)
	}
	if len(features) == 0 {
		t.Error("Expected features length to be 1, it is", len(features))
	}
	feature := features[0]
	if feature.PropertyMustString("name") != "Starbucks" {
		t.Error("Expected feature name to be Starbucks")
	}
}

func TestGetFeatureByID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	featureID := uuid.NewUUID().String()
	httpmock.RegisterResponder("GET", SpatiallyAPI+"/spatialdb/feature/"+featureID, func(req *http.Request) (*http.Response, error) {
		feature := NewFeature()
		feature.Geometry = geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
		feature.Properties = map[string]interface{}{"name": "Starbucks"}
		return httpmock.NewJsonResponse(200, feature)
	})
	feature := NewFeature()
	if err := feature.Get(sdb, featureID); err != nil {
		t.Error(err)
	}
	if feature.Geometry.IsPoint() == false {
		t.Error("Invalid feature geometry, expected Point")
	}
	if feature.PropertyMustString("name") != "Starbucks" {
		t.Error("Invalid feature name, expected Starbucks")
	}
}

func TestUpdateFeature(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	featureID := uuid.NewUUID().String()
	httpmock.RegisterResponder("PUT", SpatiallyAPI+"/spatialdb/feature/"+featureID, func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var request updateFeatureRequest
		if err := json.Unmarshal(body, &request); err != nil {
			return nil, err
		}
		if name, exists := request.Properties["name"].(string); exists {
			if name != "Starbucks Boston" {
				t.Error("Invalid request feature name")
			}
		} else {
			t.Error("Invalid request feature properties")
		}
		feature := NewFeature()
		feature.Geometry = geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
		feature.Properties = map[string]interface{}{"name": "Starbucks Boston"}
		return httpmock.NewJsonResponse(200, feature)
	})
	feature := NewFeature()
	if err := feature.Update(sdb, featureID, map[string]interface{}{
		"name": "Starbucks Boston",
	}); err != nil {
		t.Error(err)
	}
	if feature.Geometry.IsPoint() == false {
		t.Error("Invalid feature geometry, expected Point")
	}
	if feature.PropertyMustString("name") != "Starbucks Boston" {
		t.Error("Invalid feature name, expected Starbucks Boston")
	}
}

func TestDeleteFeature(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockGatewayEndpoint(t)
	sdb, err := New(applicationCode, applicationKey)
	if err != nil {
		t.Error(err)
	}
	featureID := uuid.NewUUID().String()
	httpmock.RegisterResponder("DELETE", SpatiallyAPI+"/spatialdb/feature/"+featureID, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, ""), nil
	})
	feature := NewFeature()
	if err := feature.Delete(sdb, featureID); err != nil {
		t.Error(err)
	}
}
