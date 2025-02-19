package osm

import "fmt"

type typeS struct {
	Type string `json:"type"`
}

func findType(index int, data []byte) (string, error) {
	ts := typeS{}
	if err := unmarshalJSON(data, &ts); err != nil {
		// should not happened due to previous decoding succeeded
		return "", err
	}

	if ts.Type == "" {
		return "", fmt.Errorf("could not find type in element index %d", index)
	}

	return ts.Type, nil
}
