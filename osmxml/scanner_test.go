package osmxml

import (
	"bytes"
	"io"
)

func changesetReader() io.Reader {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<osm version="0.6" generator="replicate_changesets.rb" copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" license="http://opendatacommons.org/licenses/odbl/1-0/">
  <changeset id="41226352" created_at="2016-08-03T22:40:15Z" closed_at="2016-08-04T01:41:27Z" open="false" num_changes="112" user="dragon_ear" uid="321257" min_lat="36.496286" max_lat="36.6110983" min_lon="136.6138636" max_lon="136.644462" comments_count="0">
    <tag k="comment" v="updated fire hydrant details with OsmHydrant"/>
    <tag k="created_by" v="OsmHydrant / http://yapafo.net v0.3"/>
  </changeset>
  <changeset id="41227987" created_at="2016-08-04T01:41:04Z" closed_at="2016-08-04T01:41:07Z" open="false" num_changes="7" user="MapAnalyser465" uid="3077038" min_lat="-33.7963817" max_lat="-33.7881945" min_lon="151.2527542" max_lon="151.2667464" comments_count="0">
    <tag k="comment" v="Updated Burnt Creek Deviation to Motorway Standard"/>
    <tag k="locale" v="en"/>
    <tag k="host" v="https://www.openstreetmap.org/id"/>
    <tag k="imagery_used" v="Bing"/>
    <tag k="created_by" v="iD 1.9.7"/>
  </changeset>
</osm>`)

	return bytes.NewReader(data)
}

func changesetReaderErr() io.Reader {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<osm version="0.6" generator="replicate_changesets.rb" copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" license="http://opendatacommons.org/licenses/odbl/1-0/">
  <changeset id="41226352" created_at="2016-08-03T22:40:15Z" closed_at="2016-08-04T01:41:27Z" open="false" num_changes="112" user="dragon_ear" uid="321257" min_lat="36.496286" max_lat="36.6110983" min_lon="136.6138636" max_lon="136.644462" comments_count="0">
    <tag k="comment" v="updated fire hydrant details with OsmHydrant"/>
    <tag k="created_by" v="OsmHydrant / http://yapafo.net v0.3"/>
  </changeset>
  <changeset id="41227987" created_at="2016-08-04T01:41:04Z" closed_at="2016-08-04T01:41:07Z" open="false" num_changes="7" user="MapAnalyser465" uid="3077038" min_lat="-33.7963817" max_lat="-33.7881945" min_lon="151.2527542" max_lon="151.2667464" comments_count="0">
    <tag k="comment" v="Updated Burnt Creek Deviation to Motorway Standard"/>`)

	return bytes.NewReader(data)
}

func boundsReader() io.Reader {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<osm>
	<bounds minlat="1" minlon="2" maxlat="3" maxlon="4"/>
</osm>`)

	return bytes.NewReader(data)
}

func userNoteReader() io.Reader {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<osm>
  <user id="1"></user>
  <note><id>2</id></note>
</osm>`)

	return bytes.NewReader(data)
}
