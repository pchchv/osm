package osm

var (
	// CustomJSONMarshaler can be set, to use a different json marshaler than the default one in the stdlib.
	// One use case is to include `github.com/json-iterator/go`.
	// Note that any errors that occur during marshaling will be different.
	CustomJSONMarshaler interface {
		Marshal(v interface{}) ([]byte, error)
	}

	// CustomJSONUnmarshaler can be set, to use a different json unmarshaler than the default one in the stdlib.
	// One use case is to include `github.com/json-iterator/go`.
	// Note that any errors that occur during unmarshaling will be different.
	CustomJSONUnmarshaler interface {
		Unmarshal(data []byte, v interface{}) error
	}
)
