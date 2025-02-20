package geojson

import "fmt"

// Properties defines the feature properties with some helper methods.
type Properties map[string]interface{}

// MustInt guarantees the return of an `int`
// (with optional default).
// This function useful when you explicitly want a
// `int` in a single value return context, ie:
//
//	myFunc(f.Properties.MustInt("param1"), f.Properties.MustInt("optional_param", 123))
//
// Will panic if the value is present but not a number.
func (p Properties) MustInt(key string, def ...int) int {
	v := p[key]
	if i, ok := v.(int); ok {
		return i
	}

	if f, ok := v.(float64); ok {
		return int(f)
	}

	if v != nil {
		panic(fmt.Sprintf("not a number, but a %T: %v", v, v))
	}

	if len(def) > 0 {
		return def[0]
	}

	panic("property not found")
}

// MustFloat64 guarantees the return of a `float64`
// (with optional default).
// This function useful when you explicitly want a
// `float64` in a single value return context, ie:
//
//	myFunc(f.Properties.MustFloat64("param1"), f.Properties.MustFloat64("optional_param", 10.1))
//
// Will panic if the value is present but not a number.
func (p Properties) MustFloat64(key string, def ...float64) float64 {
	v := p[key]
	if f, ok := v.(float64); ok {
		return f
	}

	if i, ok := v.(int); ok {
		return float64(i)
	}

	if v != nil {
		panic(fmt.Sprintf("not a number, but a %T: %v", v, v))
	}

	if len(def) > 0 {
		return def[0]
	}

	panic("property not found")
}
