# go-spatiallydb

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
db, err := spatiallydb.New(YOUR_APPLICATION_CODE, YOUR_APPLICATION_KEY)
if err != nil {
 log.Fatal(err)
}
```

### Create a layer & feature

```go
layer := spatiallydb.NewLayer()
if err := layer.Create(db, "businesses"); err != nil {
  log.Fatal(err)
}

feature := spatiallydb.NewFeature()
geometry := geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
properties := map[string]interface{}{
 "name": "Starbucks",
}
if err := feature.Create(db, layer.ID, geometry, properties); err != nil {
  log.Fatal(err)
}
```

### Get layer

```go
layer := spatiallydb.NewLayer()
if err := layer.Get(db, layerID); err != nil {
  log.Fatal(err)
}
```

### Get features in a polygon

```go
spatialConstraint := &spatiallydb.SpatialConstraint{
  WKT: "POLYGON((-71.06952667236328 42.35902554157146,-71.06420516967773 42.35902554157146,-71.06420516967773 42.3563616979687,-71.06952667236328 42.3563616979687,-71.06952667236328 42.35902554157146))"
}
features := spatiallydb.NewFeatures()
if err := features.GetBySpatialConstraint(db, layer.ID, spatialConstraint); err != nil {
  log.Fatal(err)
}
```

### Get features in a buffer (circle)

```go
spatialConstraint := &spatiallydb.SpatialConstraint{
  WKT: "POINT(-71.06042861938477 42.35686910545623)",
  Radius: 1000.0, // meters
}
features := spatiallydb.NewFeatures()
if err := features.GetBySpatialConstraint(db, layer.ID, spatialConstraint); err != nil {
  log.Fatal(err)
}
```

### Update a feature

```go
feature := spatiallydb.NewFeature()
if err := feature.Update(db, featureID, map[string]interface{}{
 "name": "Starbucks Boston",
}); err != nil {
 log.Fatal(err)
}
```

### Delete a feature

```go
feature := spatiallydb.NewFeature()
if err := feature.Delete(db, featureID); err != nil {
 log.Fatal(err)
}
```

### Delete a layer

```go
layer := spatiallydb.NewLayer()
if err := layer.Delete(db, layerID); err != nil {
 log.Fatal(err)
}
```
