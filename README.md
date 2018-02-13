# go-spatialdb

TODO: SpatiallyDB and what is it

## Features

* Geospatial support for all GeoJSON feature types
* Layer support, feature count and aggregation
* Group features by layer
* Read and write GeoJSON features from layer
* Intersect and buffer query support

## Coming Soon

* Within query support
* More aggregations than count, avg, sum, min, max, etc
* Bulk ingest support
* Property search/filtering

## Getting Started

1. Go to [SpatiallyDB](https://www.spatially.com/db) and create an account. Free accounts are available.
2. After creating an account, you'll be able to retrieve your API code & key. Save this in a secure place.
3. Read this documentation and godocs.
4. Get involved, provide us with feedback. We are in beta and want to hear from you.

### Create a SpatiallyDB instance

```go
spatiallyDB, err := spatialdb.New(YOUR_APPLICATION_CODE, YOUR_APPLICATION_KEY)
if err != nil {
 log.Fatal(err)
}
```

### Create a layer & feature

```go
layer, err := spatiallyDB.CreateLayer("businesses")
if err != nil {
 log.Fatal(err)
}

feature := spatialdb.NewFeature()
feature.Geometry = geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
feature.Properties = map[string]interface{}{
 "name": "Starbucks",
}

createdFeature, err := spatiallyDB.CreateFeature(layer.ID, feature)
if err != nil {
 log.Fatal(err)
}
```

### Get features in a polygon

```go
spatialConstraint := &spatialdb.SpatialConstraint{
  WKT: "POLYGON((-71.06952667236328 42.35902554157146,-71.06420516967773 42.35902554157146,-71.06420516967773 42.3563616979687,-71.06952667236328 42.3563616979687,-71.06952667236328 42.35902554157146))"
}
features, err := spatiallyDB.GetFeatures(layer.ID, spatialConstraint)
if err != nil {
  log.Fatal(err)
}
```

### Get features in a buffer (circle)

```go
spatialConstraint := &spatialdb.SpatialConstraint{
  WKT: "POINT(-71.06042861938477 42.35686910545623)",
  Radius: 1000.0 // meters
}
features, err := spatiallyDB.GetFeatures(layer.ID, spatialConstraint)
if err != nil {
  log.Fatal(err)
}
```

### Update a feature

```go
updatedFeature, err := spatiallyDB.UpdateFeature(createdFeature.ID, map[string]interface{}{
 "name": "Starbucks Boston",
})
if err != nil {
 log.Fatal(err)
}
```

### Delete a feature

```go
err := spatiallyDB.DeleteFeature(updatedFeature.ID)
if err != nil {
 log.Fatal(err)
}
```

### Delete a layer

```go
err := spatiallyDB.DeleteLayer(layer.ID)
if err != nil {
 log.Fatal(err)
}
```
