package spatially

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

// Get - Given a layer id, retrieves the layer and updates receiver
func (l *Layer) Get(db Database, id string) (err error) {
	request, err := http.NewRequest("GET", SpatiallyAPI+"/spatialdb/layer/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "get layer prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "get layer http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "get layer read response body")
	}
	if len(responseBody) == 0 {
		return fmt.Errorf("Layer not found, ID: %v", id)
	} else if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, l); err != nil {
		return errors.Wrap(err, "get layer parse response body json")
	}
	if l.ID == "" {
		log.Println(string(responseBody))
		return fmt.Errorf("There was an unexpected error getting layer ID: %v", id)
	}
	return
}

type createLayerRequest struct {
	Name string `json:"name"`
}

// Create - Given a layer name, creates the layer and updates receiver
func (l *Layer) Create(db Database, name string) (err error) {
	requestBody := createLayerRequest{
		Name: name,
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "spatialdb create layer json marshal request body")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/spatialdb/layer", body)
	if err != nil {
		return errors.Wrap(err, "spatialdb create layer prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "spatialdb create layer http post")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "spatialdb create layer read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, l); err != nil {
		return errors.Wrap(err, "spatialdb create layer parse response body json")
	}
	if l.ID == "" {
		log.Println(string(responseBody))
		return fmt.Errorf("There was an unexpected error creating layer: %v", name)
	}
	return
}

// Delete - Given a layer id, deletes the layer
func (l *Layer) Delete(db Database, id string) (err error) {
	request, err := http.NewRequest("DELETE", SpatiallyAPI+"/spatialdb/layer/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "spatialdb delete layer prepare http request")
	}
	db.PrepareRequest(request)
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

// NewLayers creates a new empty slice of layers
func NewLayers() Layers {
	return Layers{}
}

// Get - Retrieves all layers that belong to this user and updates the slice receiver.
// It does not retrieve layer features
func (l *Layers) Get(db Database) (err error) {
	request, err := http.NewRequest("GET", SpatiallyAPI+"/spatialdb/layers", nil)
	if err != nil {
		return errors.Wrap(err, "get layers prepare http request")
	}
	db.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "get layers http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "get layers read response body")
	}
	if resp.StatusCode != 200 {
		return db.Error(responseBody)
	}
	if err = json.Unmarshal(responseBody, l); err != nil {
		return errors.Wrap(err, "get layers parse response body json")
	}
	return
}
