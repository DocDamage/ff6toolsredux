//go:build go1.18
// +build go1.18

package pr

import (
	"encoding/json"
	"testing"
)

func FuzzExtractStringArray(f *testing.F) {
	f.Add([]byte(`["a","b","c"]`))
	f.Add([]byte(`[1,2,3]`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			return
		}
		_, _ = ExtractStringArray(v)
	})
}

func FuzzExtractInt64Array(f *testing.F) {
	f.Add([]byte(`[1,2,3]`))
	f.Add([]byte(`["a","b"]`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			return
		}
		_, _ = ExtractInt64Array(v)
	})
}

func FuzzExtractFloat64Array(f *testing.F) {
	f.Add([]byte(`[1.1,2.2,3.3]`))
	f.Add([]byte(`["a","b"]`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			return
		}
		_, _ = ExtractFloat64Array(v)
	})
}

func FuzzExtractString(f *testing.F) {
	f.Add([]byte(`"hello"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			return
		}
		_, _ = ExtractString(v)
	})
}

func FuzzExtractInt64(f *testing.F) {
	f.Add([]byte(`123`))
	f.Add([]byte(`"notanint"`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			return
		}
		_, _ = ExtractInt64(v)
	})
}
