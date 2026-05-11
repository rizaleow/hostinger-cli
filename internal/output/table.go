package output

import (
	"encoding/json"
	"io"
	"reflect"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// TableFn renders a payload as a human-friendly table.
type TableFn func(w io.Writer, payload any, noColor bool) error

// Column describes how to project one field out of a row into a table cell.
type Column struct {
	Header  string
	Path    string // dotted JSON path within a row, e.g. "id" or "specs.cpu"
	Hidden  bool
	MaxLen  int
	Default string
}

var (
	tablesMu sync.RWMutex
	tables   = map[reflect.Type]TableFn{}
)

// RegisterTable associates a renderer with a concrete Go type (typically a
// generated response wrapper). Pointers and non-pointers are looked up
// equivalently.
func RegisterTable(sample any, fn TableFn) {
	tablesMu.Lock()
	defer tablesMu.Unlock()
	tables[derefType(reflect.TypeOf(sample))] = fn
}

// RegisterCollectionTable is a shorthand for registering a renderer that
// projects a slice payload into a fixed set of columns. The payload may be a
// slice, a wrapper struct containing a `.Data` slice, or a pointer to either.
func RegisterCollectionTable(sample any, cols []Column) {
	RegisterTable(sample, func(w io.Writer, payload any, noColor bool) error {
		rows, err := toRows(payload)
		if err != nil {
			return err
		}
		return renderTable(w, cols, rows, noColor)
	})
}

func lookupTable(payload any) (TableFn, bool) {
	tablesMu.RLock()
	defer tablesMu.RUnlock()
	fn, ok := tables[derefType(reflect.TypeOf(payload))]
	return fn, ok
}

func derefType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func toRows(payload any) ([]map[string]any, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// First try a direct slice.
	var slice []map[string]any
	if err := json.Unmarshal(b, &slice); err == nil {
		return slice, nil
	}
	// Then try the common `{ "data": [...] }` shape.
	var wrapper struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(b, &wrapper); err == nil && wrapper.Data != nil {
		return wrapper.Data, nil
	}
	// Fall back to a single-object row.
	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err == nil {
		return []map[string]any{obj}, nil
	}
	return nil, nil
}

func renderTable(w io.Writer, cols []Column, rows []map[string]any, noColor bool) error {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	if noColor {
		t.SetStyle(table.StyleDefault)
	} else {
		s := table.StyleColoredDark
		s.Color.Header = text.Colors{text.FgHiCyan, text.Bold}
		t.SetStyle(s)
	}
	header := table.Row{}
	for _, c := range cols {
		if c.Hidden {
			continue
		}
		header = append(header, c.Header)
	}
	t.AppendHeader(header)
	for _, row := range rows {
		out := table.Row{}
		for _, c := range cols {
			if c.Hidden {
				continue
			}
			v := lookupPath(row, c.Path)
			if v == nil && c.Default != "" {
				out = append(out, c.Default)
				continue
			}
			s := formatCell(v)
			if c.MaxLen > 0 && len(s) > c.MaxLen {
				s = s[:c.MaxLen-1] + "…"
			}
			out = append(out, s)
		}
		t.AppendRow(out)
	}
	t.Render()
	return nil
}

func lookupPath(row map[string]any, path string) any {
	if path == "" {
		return nil
	}
	cur := any(row)
	for _, seg := range splitPath(path) {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[seg]
	}
	return cur
}

func splitPath(p string) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(p); i++ {
		if p[i] == '.' {
			out = append(out, p[start:i])
			start = i + 1
		}
	}
	out = append(out, p[start:])
	return out
}

func formatCell(v any) string {
	switch x := v.(type) {
	case nil:
		return "-"
	case string:
		if x == "" {
			return "-"
		}
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "?"
		}
		return string(b)
	}
}
