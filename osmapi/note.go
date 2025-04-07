package osmapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pchchv/osm"
)

var (
	_ NotesOption = Limit(1)
	_ NotesOption = MaxDaysClosed(1)
)

// Note returns the note from the osm rest api.
func (ds *Datasource) Note(ctx context.Context, id osm.NoteID) (*osm.Note, error) {
	o := &osm.OSM{}
	url := fmt.Sprintf("%s/notes/%d", ds.baseURL(), id)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	if l := len(o.Notes); l != 1 {
		return nil, fmt.Errorf("wrong number of notes, expected 1, got %v", l)
	}

	return o.Notes[0], nil
}

// Notes returns the notes in a bounding box.
// Can provide options to limit the results or change what it means to be "closed".
// See the options or osm api v0.6 docs for details.
func (ds *Datasource) Notes(ctx context.Context, bounds *osm.Bounds, opts ...NotesOption) (osm.Notes, error) {
	var err error
	params := make([]string, 0, 1+len(opts))
	params = append(params, fmt.Sprintf("bbox=%f,%f,%f,%f", bounds.MinLon, bounds.MinLat, bounds.MaxLon, bounds.MaxLat))
	for _, o := range opts {
		params, err = o.applyNotes(params)
		if err != nil {
			return nil, err
		}
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/notes?%s", ds.baseURL(), strings.Join(params, "&"))
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Notes, nil
}

// NotesSearch returns the notes whose text matches the query.
// Can provide options to limit the results or change what it means to be "closed".
// See the options or osm api v0.6 docs for details.
func (ds *Datasource) NotesSearch(ctx context.Context, query string, opts ...NotesOption) (osm.Notes, error) {
	var err error
	params := make([]string, 0, 1+len(opts))
	params = append(params, fmt.Sprintf("q=%s", url.QueryEscape(query)))
	for _, o := range opts {
		params, err = o.applyNotes(params)
		if err != nil {
			return nil, err
		}
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/notes/search?%s", ds.baseURL(), strings.Join(params, "&"))
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Notes, nil
}

// Note returns the note from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func Note(ctx context.Context, id osm.NoteID) (*osm.Note, error) {
	return DefaultDatasource.Note(ctx, id)
}

// Notes returns the notes in a bounding box. Can provide options to limit the results
// or change what it means to be "closed". See the options or osm api v0.6 docs for details.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func Notes(ctx context.Context, bounds *osm.Bounds, opts ...NotesOption) (osm.Notes, error) {
	return DefaultDatasource.Notes(ctx, bounds, opts...)
}

// NotesSearch returns the notes in a bounding box whose text matches the query.
// Can provide options to limit the results or change what it means to be "closed".
// See the options or osm api v0.6 docs for details.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func NotesSearch(ctx context.Context, query string, opts ...NotesOption) (osm.Notes, error) {
	return DefaultDatasource.NotesSearch(ctx, query, opts...)
}
