package spatialdb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// SpatialDB is the API interface to intersect with Spatially DB
type SpatialDB interface {
	// GetLayers returns a slice of all layers currently in the database
	GetLayers() (layers Layers, err error)

	// GetLayer given a layer id, returns the requested layer
	GetLayer(id string) (layer *Layer, err error)

	// CreateLayer given a name, creates a new layer in the database and returns it
	CreateLayer(name string) (layer *Layer, err error)

	// DeleteLayer given an id, delete the layer from the database
	DeleteLayer(id string) (err error)

	// GetFeatures is used to query the database by layer id. SpatialConstraint is an optional
	// parameter that when provided restricts the results to features that satisfy the spatial constraint
	// boundary
	GetFeatures(layerID string, spatialConstraint *SpatialConstraint) (features Features, err error)

	// GetFeature given an feature id, returns the feature
	GetFeature(id string) (feature *Feature, err error)

	// CreateFeature given a layer id, and feature creates a new feature in the database under the given
	// layer
	CreateFeature(layerID string, feature *Feature) (createdFeature *Feature, err error)

	// UpdateFeature given a feature id and a map of properties, updates the current feature properties
	UpdateFeature(id string, properties map[string]interface{}) (feature *Feature, err error)

	// DeleteFeature given a feature id, deletes the feature from its layer
	DeleteFeature(id string) (err error)
}

// SpatialConstraint is an object used to describe and boundary and intersection type from which
// to query features with
type SpatialConstraint struct {
	WKT    string                `json:"wkt"`
	Radius float64               `json:"radius"`
	Type   SpatialConstraintType `json:"type"`
}

// SpatialConstraintType is the type of spatial intersection to do on features
type SpatialConstraintType int

const (
	// SpatialConstraintIntersect is a SpatialContraintType that only selects features that intersect
	// with the given boundary
	SpatialConstraintIntersect SpatialConstraintType = iota
)

// SpatiallyAPI - Spatially's API URL
const SpatiallyAPI = "http://localhost:8000"

type spatialDB struct {
	Token string
}

type gatewayRequest struct {
	Code string `json:"code"`
	Key  string `json:"key"`
}

type gatewayResponse struct {
	Token string `json:"token"`
}

// New created a new instance of SpatialDB. The parameters are the api code & key
// provided by Spatially. It generates a token with the API.
func New(apiCode, apiKey string) (SpatialDB, error) {
	request := gatewayRequest{
		Code: apiCode,
		Key:  apiKey,
	}
	j, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb json marshal gateway request")
	}
	body := bytes.NewReader(j)
	resp, err := http.Post(SpatiallyAPI+"/gateway/client", "application/json", body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb gateway request")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb read gateway request response body")
	}
	var response gatewayResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, errors.Wrap(err, "spatialdb json unmarshal gateway request")
	}
	if len(response.Token) == 0 {
		return nil, errors.New("SpatialDB was not able to generate a valid token")
	}
	return spatialDB{Token: response.Token}, nil
}

func (s spatialDB) prepareRequest(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+s.Token)
}

type requestError struct {
	Message string `json:"message"`
}

func (s spatialDB) requestError(responseBody []byte) error {
	var err requestError
	if err := json.Unmarshal(responseBody, &err); err != nil {
		return err
	}
	return errors.New(err.Message)
}
