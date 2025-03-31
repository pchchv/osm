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