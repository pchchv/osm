package osm

import "encoding/json"

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

func marshalJSON(v interface{}) ([]byte, error) {
	if CustomJSONMarshaler == nil {
		return json.Marshal(v)
	}

	return CustomJSONMarshaler.Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	if CustomJSONUnmarshaler == nil {
		return json.Unmarshal(data, v)
	}

	return CustomJSONUnmarshaler.Unmarshal(data, v)
}
