//go:build go1.18
// +build go1.18

package pr

import (
	"encoding/json"
	"testing"
)

func FuzzDecodeInterfaceSlice(f *testing.F) {
	f.Add([]byte(`[1,2,3]`))
	f.Add([]byte(`{"foo":1}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`"string"`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		err := json.Unmarshal(data, &v)
		if err != nil {
			return
		}
		_, _ = decodeInterfaceSlice(v, "fuzz")
	})
}

func FuzzDecodeString(f *testing.F) {
	f.Add([]byte(`"hello"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		err := json.Unmarshal(data, &v)
		if err != nil {
			return
		}
		_, _ = decodeString(v, "fuzz")
	})
}

func FuzzDecodeRawMessage(f *testing.F) {
	f.Add([]byte(`"{\\\"foo\\\":1}"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var v interface{}
		err := json.Unmarshal(data, &v)
		if err != nil {
			return
		}
		_, _ = decodeRawMessage(v, "fuzz")
	})
}
