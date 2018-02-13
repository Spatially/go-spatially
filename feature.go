package spatiallydb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	geojson "github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
)

// Features is a slice of Feature
type Features []*Feature

// NewFeatures creates a new empty slice of Features
func NewFeatures() Features {
	return Features{}
}

type getFeaturesRequest struct {
	LayerID           string             `json:"layer"`
	SpatialConstraint *SpatialConstraint `json:"spatialConstraint"`
}

// GetByLayer - Given a layer id, retrieves the features that belong to it and updates the slice receiver
func (f *Features) GetByLayer(db SpatiallyDB, layerID string) (err error) {
	requestBody := getFeaturesRequest{
		LayerID: layerID,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "get features by layer json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/features", body)
	if err != nil {
		return errors.Wrap(err, "get features by layer prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "get features by layer http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "get features by layer read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, f); err != nil {
		return errors.Wrap(err, "get features by layer parse response body json")
	}
	return
}

// GetBySpatialConstraint - Given a layer id and spatial constraint object, retrieves all features that satisfay the constraint
// and updates the slice receiver
func (f *Features) GetBySpatialConstraint(db SpatiallyDB, layerID string, spatialConstraint *SpatialConstraint) (err error) {
	requestBody := getFeaturesRequest{
		LayerID:           layerID,
		SpatialConstraint: spatialConstraint,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "get features by spatial constraint json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/features", body)
	if err != nil {
		return errors.Wrap(err, "get features by spatial constraint prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "get features by spatial constraint http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "get features by spatial constraint read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, f); err != nil {
		return errors.Wrap(err, "get features by spatial constraint parse response body json")
	}
	return
}

// Feature is a wrapped geojson.Feature
type Feature struct {
	*geojson.Feature
}

// NewFeature creates a new SpatiallyDB feature
func NewFeature() *Feature {
	return &Feature{
		Feature: &geojson.Feature{},
	}
}

// NewFeatureFromWKT will create a new feature given a Well Known Text shape
func NewFeatureFromWKT(wkt string) (feature *Feature, err error) {
	feature = &Feature{
		Feature: &geojson.Feature{},
	}
	g, err := WKTToGeometry(wkt)
	if err != nil {
		return feature, err
	}
	feature.Geometry = g
	return
}

// Get - Given a feature id, retrieves the feature and updates the receiver
func (f *Feature) Get(db SpatiallyDB, id string) (err error) {
	request, err := http.NewRequest("GET", SpatiallyAPI+"/spatialdb/feature/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "get feature prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "get feature http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "get feature read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, f); err != nil {
		return errors.Wrap(err, "get feature parse response body json")
	}
	return
}

type createFeatureRequest struct {
	LayerID string           `json:"layer"`
	Feature *geojson.Feature `json:"feature"`
}

// Create - given a layer id, geometry and properties - creates the feature and updates the receiver with the created feature.
// It also increases the layer feature count
func (f *Feature) Create(db SpatiallyDB, layerID string, geometry *geojson.Geometry, properties map[string]interface{}) (err error) {
	f.Geometry = geometry
	f.Properties = properties
	requestBody := createFeatureRequest{
		LayerID: layerID,
		Feature: f.Feature,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "create feature json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/feature", body)
	if err != nil {
		return errors.Wrap(err, "create feature prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "create feature http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "create feature read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, f); err != nil {
		return errors.Wrap(err, "create feature parse response body json")
	}
	return
}

type updateFeatureRequest struct {
	Properties map[string]interface{} `json:"properties"`
}

// Update - Given a feature id and properties, it updates the feature and receiver
func (f *Feature) Update(db SpatiallyDB, id string, properties map[string]interface{}) (err error) {
	requestBody := updateFeatureRequest{
		Properties: properties,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "update feature json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("PUT", SpatiallyAPI+"/spatialdb/feature/"+id, body)
	if err != nil {
		return errors.Wrap(err, "update feature prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "update feature http put")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "update feature read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, f); err != nil {
		return errors.Wrap(err, "update feature parse response body json")
	}
	return
}

// Delete - Given a feature id, it deletes the feature and decreases the layer's feature count
func (f *Feature) Delete(db SpatiallyDB, id string) (err error) {
	request, err := http.NewRequest("DELETE", SpatiallyAPI+"/spatialdb/feature/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "delete feature prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "delete feature http delete")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("There was an unexpected error deleting feature ID: %v", id)
	}
	return
}
