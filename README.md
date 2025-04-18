# OSM [![CI](https://github.com/pchchv/osm/workflows/CI/badge.svg)](https://github.com/pchchv/osm/actions?query=workflow%3ACI+event%3Apush) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/osm)](https://goreportcard.com/report/github.com/pchchv/osm) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/osm)](https://pkg.go.dev/github.com/pchchv/osm)

**OSM** package is a general purpose library for reading, writing and working with [OpenStreetMap](https://osm.org) data in Go.   
It has the ability to:
- read/write [OSM XML](https://wiki.openstreetmap.org/wiki/OSM_XML)
- read/write [OSM JSON](https://wiki.openstreetmap.org/wiki/OSM_JSON), a format returned by the Overpass API.
- efficiently parse [OSM PBF](https://wiki.openstreetmap.org/wiki/PBF_Format) data files available at [planet.osm.org](https://planet.osm.org/)

Made available by the package are the following types:
- Node
- Way
- Relation
- Changeset
- Note
- User

Following “container” types:
- OSM - container returned via API
- Change - used by the replication API
- Diff - corresponds to [Overpass Augmented Diffs](https://wiki.openstreetmap.org/wiki/Overpass_API/Augmented_Diffs)

## Concepts

In the `osm` package, the core OSM data types are referred **Objects**. The Node, Way, Relation, Changeset, Note and User types implement the `osm.Object` interface and can be referenced using the `osm.ObjectID` type. As a result, it is possible to have a `[]osm.Object` slice containing nodes, changesets and users.

Individual versions of the core OSM map data types are referred **Elements**, and the set of versions for a given Node, Way or Relation is referred a **Feature**. For example, `osm.ElementID` might refer to "Node with ID 10 and version 3" and `osm.FeatureID` might refer to "all versions of a node with ID 10". In another way, features represent a road and how it has changed over time, and an element is a specific version of that feature.

A number of helper methods are provided for working with features and elements. The idea is to simplify working with, for example, Way and its member nodes.

## Working with JSON

`osm` package supports reading and writing [OSM JSON](https://wiki.openstreetmap.org/wiki/OSM_JSON). This format is returned by the Overpass API and can be optionally returned by the [OSM API](https://wiki.openstreetmap.org/wiki/API_v0.6#JSON_Format).

If performance is important, third party "encoding/json" replacements such as [github.com/json-iterator/go](https://github.com/json-iterator/go) are supported.  
They can be enabled with something like this:

```go
import (
  jsoniter "github.com/json-iterator/go"
  "github.com/pchchv/osm"
)

var c = jsoniter.Config{
  EscapeHTML:              true,
  SortMapKeys:             false,
  MarshalFloatWith6Digits: true,
}.Froze()

osm.CustomJSONMarshaler = c
osm.CustomJSONUnmarshaler = c
```

The above change can have significant performance implications, see the benchmarks below on a large OSM Change object.

```
benchmark                            old ns/op     new ns/op     delta
BenchmarkChange_MarshalJSON-12       604496        461349        -23.68%
BenchmarkChange_UnmarshalJSON-12     1633834       1051630       -35.63%

benchmark                            old allocs    new allocs    delta
BenchmarkChange_MarshalJSON-12       1277          1081          -15.35%
BenchmarkChange_UnmarshalJSON-12     5133          8580          +67.15%

benchmark                            old bytes     new bytes     delta
BenchmarkChange_MarshalJSON-12       180583        162727        -9.89%
BenchmarkChange_UnmarshalJSON-12     287707        317723        +10.43%
```

## Scanning large data files

For small data it is possible to use the `encoding/xml` package in the Go stdlib to marshal/unmarshal the data. This is typically done using the `osm.OSM` or `osm.Change` "container" structs.

For large data the package defines the `Scanner` interface implemented in both the [osmxml](osmxml) and [osmpbf](osmpbf) sub-packages.

```go
type osm.Scanner interface {
	Scan() bool
	Object() osm.Object
	Err() error
	Close() error
}
```

This interface is designed to mimic the [bufio.Scanner](https://golang.org/pkg/bufio/#Scanner)
interface found in the Go standard library.

Example usage:

```go
f, err := os.Open("./delaware-latest.osm.pbf")
if err != nil {
	panic(err)
}
defer f.Close()

scanner := osmpbf.New(context.Background(), f, 3)
defer scanner.Close()

for scanner.Scan() {
	o := scanner.Object()
	// do something
}

scanErr := scanner.Err()
if scanErr != nil {
	panic(scanErr)
}
```

**Note:** Scanners are **not** safe for parallel use. Objects must be fed into the channel and workers must read from it.

## CGO and zlib

OSM PBF data comes in blocks, each block is zlib compressed. Decompressing this data takes about 33% of the total read time. [DataDog/czlib](https://github.com/DataDog/czlib) is used to speed this process. See [osmpbf/README.md](osmpbf#using-cgoczlib-for-decompression) for more details.

As a result, a C compiler is necessary to install this module. On macOS this may require installing pkg-config using something like `brew install pkg-config`

CGO can be disabled at build time using the `CGO_ENABLED` ENV variable. For example, `CGO_ENABLED=0 go build`. The code will fallback to the stdlib implementation of zlib.

## List of sub-package utilities

- [`osmapi`](osmapi) - supports all the v0.6 read/data endpoints
- [`osmpbf`](osmpbf) - stream processing of `*.osm.pbf` files
- [`osmxml`](osmxml) - stream processing of `*.osm` xml files
- [`annotate`](annotate) - adds lon/lat, version, changeset and orientation data to way and relation members
- [`osmgeojson`](osmgeojson) - converts OSM data to GeoJSON
- [`replication`](replication) - fetch replication state and change files