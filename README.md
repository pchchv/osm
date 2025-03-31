# OSM [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/osm)](https://goreportcard.com/report/github.com/pchchv/osm) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/osm)](https://pkg.go.dev/github.com/pchchv/osm)

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

In the `OSM` package, the core OSM data types are referred **Objects**. The Node, Way, Relation, Changeset, Note and User types implement the `osm.Object` interface and can be referenced using the `osm.ObjectID` type. As a result, it is possible to have a `[]osm.Object` slice containing nodes, changesets and users.

Individual versions of the core OSM map data types are referred **Elements**, and the set of versions for a given Node, Way or Relation is referred a **Feature**. For example, `osm.ElementID` might refer to "Node with ID 10 and version 3" and `osm.FeatureID` might refer to "all versions of a node with ID 10". In another way, features represent a road and how it has changed over time, and an element is a specific version of that feature.

A number of helper methods are provided for working with features and elements. The idea is to simplify working with, for example, Way and its member nodes.