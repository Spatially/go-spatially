package spatialdb

import geojson "github.com/paulmach/go.geojson"

//
type SpatialDB interface {
	// Layers
	GetLayers() Layers
	GetLayer(name string) *Layer
	CreateLayer(name string) *Layer
	DeleteLayer(name string)

	// Features
	GetFeatures(layer string, spatialConstraint *SpatialConstraint) Features
	GetFeature(id string) *Feature
	CreateFeature(layer string, feature *Feature) *Feature
	UpdateFeature(id string, properties map[string]interface{}) *Feature
	DeleteFeature(id string)

	// Grid - TODO(josebalius): Fledge out
}

//
type SpatialConstraint struct {
	WKT    string                `json:"wkt"`
	Radius float64               `json:"radius"`
	Type   SpatialConstraintType `json:"type"`
}

//
type SpatialConstraintType int

//
const (
	SpatialConstraintWithin SpatialConstraintType = iota
	SpatialConstraintIntersect
)

//
type Layers []*Layer

//
type Layer struct {
	ID   string
	Name string
}

//
type Features []*Feature

//
type Feature struct {
	*geojson.Feature
}

//
func New(apiKey string) {

}
