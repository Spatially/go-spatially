package spatialdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Layers is a slice of Layer
type Layers []*Layer

// Layer represents a group of features in the database.
type Layer struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	FeatureCount int    `json:"featureCount"`
}

// NewLayer creates a new empty Layer
func NewLayer() *Layer {
	return &Layer{}
}

// NewLayers creates a new empty slice of layers
func NewLayers() Layers {
	return Layers{}
}

func (s spatialDB) GetLayers() (layers Layers, err error) {
	request, err := http.NewRequest("GET", SpatiallyAPI+"/spatialdb/layers", nil)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get layers prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get layers http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get layers read response body")
	}
	if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	layers = NewLayers()
	if err = json.Unmarshal(responseBody, &layers); err != nil {
		return nil, errors.Wrap(err, "spatialdb get layers parse response body json")
	}
	return
}

func (s spatialDB) GetLayer(id string) (layer *Layer, err error) {
	request, err := http.NewRequest("GET", SpatiallyAPI+"/spatialdb/layer/"+id, nil)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get layer prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get layer http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb get layer read response body")
	}
	if len(responseBody) == 0 {
		return nil, fmt.Errorf("Layer not found, ID: %v", id)
	} else if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	layer = NewLayer()
	if err = json.Unmarshal(responseBody, layer); err != nil {
		return nil, errors.Wrap(err, "spatialdb get layer parse response body json")
	}
	if layer.ID == "" {
		log.Println(string(responseBody))
		return nil, fmt.Errorf("There was an unexpected error getting layer ID: %v", id)
	}
	return
}

type createLayerRequest struct {
	Name string `json:"name"`
}

func (s spatialDB) CreateLayer(name string) (layer *Layer, err error) {
	requestBody := createLayerRequest{
		Name: name,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create layer json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/layer", body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create layer prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create layer http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb create layer read response body")
	}
	if resp.StatusCode != 200 {
		return nil, s.requestError(responseBody)
	}
	layer = NewLayer()
	if err = json.Unmarshal(responseBody, layer); err != nil {
		return nil, errors.Wrap(err, "spatialdb create layer parse response body json")
	}
	if layer.ID == "" {
		log.Println(string(responseBody))
		return nil, fmt.Errorf("There was an unexpected error creating layer: %v", name)
	}
	return
}

func (s spatialDB) DeleteLayer(id string) (err error) {
	request, err := http.NewRequest("DELETE", SpatiallyAPI+"/spatialdb/layer/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "spatialdb delete layer prepare http request")
	}
	s.prepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "spatialdb delete layer http delete")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("There was an unexpected error deleting layer ID: %v", id)
	}
	return
}
