# go-spatially

## Features

* Active Trade Area generations based on mobile data observations and machine learning
* Active Trade Area Geofence support - get the output you need for Google & Facebook
* Geospatial Database support for all GeoJSON feature types
* Layer support, feature count and aggregation
* Group features by layer
* Read and write GeoJSON features from layer
* Intersect and buffer query support

## Coming Soon

* Within query support
* More aggregations than count, avg, sum, min, max, etc
* Bulk ingest support
* Property search/filtering
* Grid search
* Geofencing Support & Notifications

## Getting Started

1.  Go to [Spatially Labs](https://labs.spatially.com) and create an account. Free accounts are available.
2.  After creating an account, you'll be able to retrieve your API code & key. Save this in a secure place.
3.  Read this documentation and godocs.
4.  Get involved, provide us with feedback. We are in beta and want to hear from you.

### Create a spatially instance

```go
api, err := spatially.NewAPI(YOUR_APPLICATION_CODE, YOUR_APPLICATION_KEY)
if err != nil {
 log.Fatal(err)
}
```

### Create an ATA (Active Trade Area) of Home locations

```go
ata, err := spatially.NewATA(api, "POINT(-71.064156780428 42.35862883483673)", nil)
if err != nil {
  log.Fatal(err)
}
log.Printf("%+v", *ata.FeatureCollection)
```

### Create an ATA (Active Trade Area) of Work Locations

```go
ata, err := spatially.NewATA(api, "POINT(-71.064156780428 42.35862883483673)", &spatially.ATAOptions{
  LocationType: spatially.Work,
})
if err != nil {
  log.Fatal(err)
}
log.Printf("%+v", *ata.FeatureCollection)
```

### Create an ATA (Active Trade Area) of Home And Work Locations based on Morning observations

```go
ata, err := spatially.NewATA(api, "POINT(-71.064156780428 42.35862883483673)", &spatially.ATAOptions{
  LocationType: spatially.HomeAndWork,
  TimeOfDay: spatially.Morning,
})
if err != nil {
  log.Fatal(err)
}
log.Printf("%+v", *ata.FeatureCollection)
```

### Create an ATA - GeoFenced

This method returns the same ATA you would receive usually but instead it has calculated the necessary geofences (point + radius buffers) to cover the ATA. This output is usually used for advertising networks such as Google or Facebook.

```go
ata, err := spatially.NewATA(api, "POINT(-71.064156780428 42.35862883483673)", &spatially.ATAOptions{
  GeoFence: true,
})
if err != nil {
  log.Fatal(err)
}
log.Printf("%+v", *ata.FeatureCollection)
```

### Create a layer & feature

```go
layer := spatially.NewLayer()
if err := layer.Create(api, "businesses"); err != nil {
  log.Fatal(err)
}

feature := spatially.NewFeature()
geometry := geojson.NewPointGeometry([]float64{-71.06772422790527, 42.35848049347556})
properties := map[string]interface{}{
 "name": "Starbucks",
}
if err := feature.Create(api, layer.ID, geometry, properties); err != nil {
  log.Fatal(err)
}
```

### Get layer

```go
layer := spatially.NewLayer()
if err := layer.Get(api, layerID); err != nil {
  log.Fatal(err)
}
```

### Get features in a polygon

```go
spatialConstraint := &spatially.SpatialConstraint{
  WKT: "POLYGON((-71.06952667236328 42.35902554157146,-71.06420516967773 42.35902554157146,-71.06420516967773 42.3563616979687,-71.06952667236328 42.3563616979687,-71.06952667236328 42.35902554157146))"
}
features := spatially.NewFeatures()
if err := features.GetBySpatialConstraint(api, layer.ID, spatialConstraint); err != nil {
  log.Fatal(err)
}
```

### Get features in a buffer (circle)

```go
spatialConstraint := &spatially.SpatialConstraint{
  WKT: "POINT(-71.06042861938477 42.35686910545623)",
  Radius: 1000.0, // meters
}
features := spatially.NewFeatures()
if err := features.GetBySpatialConstraint(api, layer.ID, spatialConstraint); err != nil {
  log.Fatal(err)
}
```

### Update a feature

```go
feature := spatially.NewFeature()
if err := feature.Update(api, featureID, map[string]interface{}{
 "name": "Starbucks Boston",
}); err != nil {
 log.Fatal(err)
}
```

### Delete a feature

```go
feature := spatially.NewFeature()
if err := feature.Delete(api, featureID); err != nil {
 log.Fatal(err)
}
```

### Delete a layer

```go
layer := spatially.NewLayer()
if err := layer.Delete(api, layerID); err != nil {
 log.Fatal(err)
}
```

### Grid - Population

```go
pop, err := spatially.Population(api, "POINT(-71.064156780428 42.35862883483673)", 150)
if err != nil {
  log.Fatal(err)
}
```

### Grid - Trade Area Market Size

```go
marketSize, err := spatially.TradeAreaMarketSize(api, ata)
if err != nil {
  log.Fatal(err)
}
```

### Grid - Popular Times

```go
times, err := spatially.PopularTimes(api, "POINT(-71.064156780428 42.35862883483673)", 150)
if err != nil {
  log.Fatal(err)
}
```

### Grid - Distance Sensitivity

```go
ds, err := spatially.DistanceSensitivity(api, "POINT(-71.064156780428 42.35862883483673)", 150)
if err != nil {
  log.Fatal(err)
}
```

### Grid - Demographics

```go
demographics, err := spatially.Demographics(api, "POINT(-71.064156780428 42.35862883483673)", 150)
if err != nil {
  log.Fatal(err)
}
```
