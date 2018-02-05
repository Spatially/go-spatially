package spatialdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	geojson "github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
)

//
type Features []*Feature

//
func NewFeatures() Features {
	return Features{}
}

//
type Feature struct {
	*geojson.Feature
}

//
func NewFeature() *Feature {
	return &Feature{
		Feature: &geojson.Feature{},
	}
}

//
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

type getFeaturesRequest struct {
	LayerID           string             `json:"layer"`
	SpatialConstraint *SpatialConstraint `json:"spatialConstraint"`
}

func (s spatialDB) GetFeatures(layerID string, spatialConstraint *SpatialConstraint) (features Features, err error) {
	requestBody := getFeaturesRequest{
		LayerID:           layerID,
		SpatialConstraint: spatialConstraint,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get features json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/features", body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get features prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get features http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get features read response body")
	}
	if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	features = NewFeatures()
	if err = json.Unmarshal(responseBody, &features); err != nil {
		return nil, errors.Wrap(err, "spatialdb get features parse response body json")
	}
	return
}

func (s spatialDB) GetFeature(id string) (feature *Feature, err error) {
	request, err := http.NewRequest("GET", SpatiallyAPI+"/spatialdb/feature/"+id, nil)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get feature prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get feature http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get feature read response body")
	}
	if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	feature = NewFeature()
	if err = json.Unmarshal(responseBody, feature); err != nil {
		return nil, errors.Wrap(err, "spatialdb get feature parse response body json")
	}
	return
}

type createFeatureRequest struct {
	LayerID string           `json:"layer"`
	Feature *geojson.Feature `json:"feature"`
}

func (s spatialDB) CreateFeature(layerID string, feature *Feature) (createdFeature *Feature, err error) {
	requestBody := createFeatureRequest{
		LayerID: layerID,
		Feature: feature.Feature,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create feature json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/feature", body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create feature prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create feature http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create feature read response body")
	}
	if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	createdFeature = NewFeature()
	if err = json.Unmarshal(responseBody, createdFeature); err != nil {
		return nil, errors.Wrap(err, "spatialdb create feature parse response body json")
	}
	return
}

type updateFeatureRequest struct {
	Properties map[string]interface{} `json:"properties"`
}

func (s spatialDB) UpdateFeature(id string, properties map[string]interface{}) (feature *Feature, err error) {
	requestBody := updateFeatureRequest{
		Properties: properties,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb update feature json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("PUT", SpatiallyAPI+"/spatialdb/feature/"+id, body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb update feature prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb update feature http put")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb update feature read response body")
	}
	if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	feature = NewFeature()
	if err = json.Unmarshal(responseBody, feature); err != nil {
		return nil, errors.Wrap(err, "spatialdb update feature parse response body json")
	}
	return
}

func (s spatialDB) DeleteFeature(id string) (err error) {
	request, err := http.NewRequest("DELETE", SpatiallyAPI+"/spatialdb/feature/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "spatialdb delete feature prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "spatialdb delete feature http delete")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("There was an unexpected error deleting feature ID: %v", id)
	}
	return
}
