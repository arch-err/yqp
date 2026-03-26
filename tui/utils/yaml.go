package utils

import (
	"bytes"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"gopkg.in/yaml.v3"
)

// IsValidYAML checks if the input bytes are valid YAML.
func IsValidYAML(data []byte) error {
	var out any
	return yaml.Unmarshal(data, &out)
}

// IsValidInput checks the validity of input data as YAML.
// Returns true if valid, along with any error.
func IsValidInput(data []byte) (bool, error) {
	if len(data) == 0 {
		return false, ErrInvalidYAML
	}

	if err := IsValidYAML(data); err != nil {
		return false, ErrInvalidYAML
	}

	return true, nil
}

// ErrInvalidYAML is returned when the input data is not valid YAML.
var ErrInvalidYAML = &invalidYAMLError{}

type invalidYAMLError struct{}

func (e *invalidYAMLError) Error() string {
	return "data is not valid YAML"
}

func highlightYAML(w io.Writer, source string, style *chroma.Style) error {
	l := lexers.Get("yaml")
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	f := formatters.Get(getTerminalColorSupport())
	if f == nil {
		f = formatters.Fallback
	}

	if style == nil {
		style = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return f.Format(w, style, it)
}

// Prettify takes YAML input bytes and returns syntax-highlighted output.
func Prettify(inputYAML []byte, chromaStyle *chroma.Style) (*bytes.Buffer, error) {
	// YAML is already human-readable, no indentation step needed.
	// Just syntax-highlight the raw input.
	var buf bytes.Buffer
	err := highlightYAML(&buf, string(inputYAML), chromaStyle)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
