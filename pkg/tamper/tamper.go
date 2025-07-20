package tamper

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func Load(tamperFilepath string) (Tamper, error) {
	p, err := plugin.Open(tamperFilepath)
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("Tamper")
	if err != nil {
		return nil, err
	}

	tamper, ok := sym.(Tamper)
	if !ok {
		return nil, fmt.Errorf("invalid plugin type: %v", tamperFilepath)
	}

	return tamper, nil
}

// ValidatePlugin tries to open the plugin file and verify
// it exports a Tamper symbol implementing the interface.
/* func ValidatePlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	// Assume symbol name is "Tamper"
	sym, err := p.Lookup("Tamper")
	if err != nil {
		return errors.New("symbol 'Tamper' not found in plugin")
	}

	_, ok := sym.(Tamper)
	if !ok {
		return errors.New("symbol 'Tamper' does not implement Tamper interface")
	}
	return nil
} */

// GetAvailableTampers loads all .so files from dir,
// validates them, and returns a list of their Names.
func GetAvailableTampers(dir string) (map[string]Tamper, error) {
	var tampers = make(map[string]Tamper)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".so") {
			continue
		}

		path := filepath.Join(dir, e.Name())

		p, err := plugin.Open(path)
		if err != nil {
			continue // skip invalid plugin
		}

		sym, err := p.Lookup("Tamper")
		if err != nil {
			continue
		}

		tamper, ok := sym.(Tamper)
		if !ok {
			continue
		}

		tampers[tamper.Name()] = tamper

	}

	return tampers, nil
}
