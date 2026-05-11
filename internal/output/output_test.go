package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	err := Render(map[string]any{"hello": "world"}, Options{Format: FormatJSON, Out: &buf})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"hello": "world"`) {
		t.Errorf("output missing JSON pair: %q", buf.String())
	}
}

func TestRenderJSONCompact(t *testing.T) {
	var buf bytes.Buffer
	err := Render(map[string]any{"a": 1}, Options{Format: FormatJSON, Compact: true, Out: &buf})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != `{"a":1}` {
		t.Errorf("compact output = %q", got)
	}
}

func TestRenderYAML(t *testing.T) {
	var buf bytes.Buffer
	err := Render(map[string]any{"name": "vps", "id": 42}, Options{Format: FormatYAML, Out: &buf})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "name: vps") || !strings.Contains(out, "id: 42") {
		t.Errorf("yaml output missing fields: %q", out)
	}
}

func TestRenderJQ(t *testing.T) {
	var buf bytes.Buffer
	payload := []map[string]any{{"name": "fra"}, {"name": "ams"}}
	err := Render(payload, Options{Format: FormatJSON, JQ: ".[0].name", Out: &buf})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != `"fra"` {
		t.Errorf("jq output = %q", got)
	}
}

func TestRenderTemplate(t *testing.T) {
	var buf bytes.Buffer
	err := Render(map[string]any{"x": "y"},
		Options{Format: FormatTemplate, Template: "x is {{.x}}", Out: &buf})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != "x is y" {
		t.Errorf("template output = %q", got)
	}
}
