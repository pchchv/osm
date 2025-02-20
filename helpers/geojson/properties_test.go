package geojson

import "testing"

func TestPropertiesClone(t *testing.T) {
	props := Properties{
		"one": 2,
	}

	clone := props.Clone()
	if clone["one"] != 2 {
		t.Errorf("should clone properties")
	}

	clone["one"] = 3
	if props["one"] != 2 {
		t.Errorf("should clone properties")
	}
}
