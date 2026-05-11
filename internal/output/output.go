// Package output renders payloads as tables, JSON, YAML, or templates,
// auto-detecting the right default based on whether stdout is a TTY.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/itchyny/gojq"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v3"
)

// Format identifies the chosen output renderer.
type Format string

const (
	FormatAuto     Format = "auto"
	FormatTable    Format = "table"
	FormatJSON     Format = "json"
	FormatYAML     Format = "yaml"
	FormatTemplate Format = "template"
)

// Options configures Render.
type Options struct {
	Format   Format
	Template string // used when Format == FormatTemplate
	JQ       string // optional jq filter applied to the JSON projection
	NoColor  bool
	Compact  bool
	Out      io.Writer
}

// Resolve picks an effective format based on Options and the TTY state of out.
func Resolve(opts Options) Format {
	if opts.Format != "" && opts.Format != FormatAuto {
		return opts.Format
	}
	if env := os.Getenv("HOSTINGER_OUTPUT"); env != "" {
		return Format(env)
	}
	if isTTY(opts.Out) && !isCI() {
		return FormatTable
	}
	return FormatJSON
}

// Render writes payload according to opts. Tables are only used when a
// TableFn has been registered for the payload's Go type; otherwise it falls
// back to JSON. Errors from the rendering pipeline (e.g. invalid jq) are
// returned to the caller and should be reported on stderr.
func Render(payload any, opts Options) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if payload == nil {
		return nil
	}

	if opts.JQ != "" {
		filtered, err := applyJQ(payload, opts.JQ)
		if err != nil {
			return err
		}
		payload = filtered
		// jq output is data-shaped; force JSON unless the caller asked for
		// something specific.
		if opts.Format == "" || opts.Format == FormatAuto {
			opts.Format = FormatJSON
		}
	}

	switch Resolve(opts) {
	case FormatJSON:
		return writeJSON(opts.Out, payload, opts.Compact)
	case FormatYAML:
		return writeYAML(opts.Out, payload)
	case FormatTemplate:
		return writeTemplate(opts.Out, payload, opts.Template)
	case FormatTable:
		if fn, ok := lookupTable(payload); ok {
			return fn(opts.Out, payload, opts.NoColor)
		}
		// No registered table renderer — fall back to indented JSON so we
		// still produce something useful instead of erroring out.
		return writeJSON(opts.Out, payload, false)
	default:
		return fmt.Errorf("unknown output format: %q", opts.Format)
	}
}

func writeJSON(w io.Writer, v any, compact bool) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if !compact {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}

func writeYAML(w io.Writer, v any) error {
	// Round-trip through JSON so generated types (which only have json tags)
	// produce sensible YAML keys instead of Go field names.
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return enc.Close()
}

func writeTemplate(w io.Writer, v any, tmpl string) error {
	if tmpl == "" {
		return fmt.Errorf("--output template requires --template")
	}
	t, err := template.New("out").Parse(tmpl)
	if err != nil {
		return err
	}
	if err := t.Execute(w, v); err != nil {
		return err
	}
	_, _ = io.WriteString(w, "\n")
	return nil
}

func applyJQ(payload any, expr string) (any, error) {
	q, err := gojq.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("parse --jq: %w", err)
	}
	// Normalize payload through JSON so gojq sees plain map/slice/scalar values.
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	iter := q.Run(data)
	var out []any
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			return nil, err
		}
		out = append(out, v)
	}
	if len(out) == 1 {
		return out[0], nil
	}
	return out, nil
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

func isCI() bool {
	for _, k := range []string{"CI", "GITHUB_ACTIONS", "BUILDKITE", "GITLAB_CI"} {
		if os.Getenv(k) != "" {
			return true
		}
	}
	return false
}
