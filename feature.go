package spatialdb

import geojson "github.com/paulmach/go.geojson"

//
type Features []*Feature

//
type Feature struct {
	*geojson.Feature
}

func (s spatialDB) GetFeatures(layer string, spatialConstraint *SpatialConstraint) (features Features, err error) {
	return
}

func (s spatialDB) GetFeature(id string) (feature *Feature, err error) {
	return
}

func (s spatialDB) CreateFeature(layer string, feature *Feature) (createdFeature *Feature, err error) {
	return
}

func (s spatialDB) UpdateFeature(id string, properties map[string]interface{}) (feature *Feature, err error) {
	return
}

func (s spatialDB) DeleteFeature(id string) (err error) {
	return
}
