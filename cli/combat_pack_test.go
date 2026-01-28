package cli

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput redirects stdout for the duration of fn and returns what was printed.
func captureOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	fnErr := fn()
	_ = w.Close()
	os.Stdout = old

	b, _ := io.ReadAll(r)
	_ = r.Close()

	return string(b), fnErr
}

func TestCombatPackSmokeCLI(t *testing.T) {
	cli := &CLI{args: []string{"combat-pack", "--mode", "smoke"}}

	out, err := captureOutput(func() error {
		return cli.Run()
	})
	if err != nil {
		t.Fatalf("combat-pack smoke run failed: %v", err)
	}

	// Expect JSON with passed/failed/total keys.
	dec := json.NewDecoder(strings.NewReader(out))
	var payload map[string]interface{}
	if err := dec.Decode(&payload); err != nil {
		t.Fatalf("failed to decode smoke JSON output: %v; raw=%s", err, out)
	}
	if _, ok := payload["passed"]; !ok {
		t.Fatalf("smoke output missing 'passed' field: %v", payload)
	}
	if failed, ok := payload["failed"]; ok {
		if n, ok := failed.(float64); ok && n != 0 {
			t.Fatalf("smoke tests reported failures: %v", payload)
		}
	}
}
