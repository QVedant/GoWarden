package registry

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Language struct {
	DisplayName    string   `yaml:"display_name"`
	Extension      string   `yaml:"extension"`
	Compile        []string `yaml:"compile"`
	Run            []string `yaml:"run"`
	TimeoutSeconds int      `yaml:"timeout_seconds"`
	MemoryMB       int      `yaml:"memory_mb"`
	ImageToolchain string   `yaml:"image_toolchain"`
}

type Registry struct {
	Languages map[string]Language `yaml:"languages"`
}

func Load(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading registry file: %w", err)
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parsing registry YAML: %w", err)
	}

	if err := reg.validate(); err != nil {
		return nil, fmt.Errorf("invalid registry: %w", err)
	}

	return &reg, nil
}

func (r *Registry) validate() error {
	if len(r.Languages) == 0 {
		return fmt.Errorf("registry contains no languages")
	}

	for name, lang := range r.Languages {
		if lang.DisplayName == "" {
			return fmt.Errorf("language %q: display_name is required", name)
		}
		if !strings.HasPrefix(lang.Extension, ".") {
			return fmt.Errorf("language %q: extension must start with '.', got %q", name, lang.Extension)
		}
		if len(lang.Run) == 0 {
			return fmt.Errorf("language %q: run command must not be empty", name)
		}
		if lang.TimeoutSeconds <= 0 {
			return fmt.Errorf("language %q: timeout_seconds must be positive", name)
		}
		if lang.MemoryMB <= 0 {
			return fmt.Errorf("language %q: memory_mb must be positive", name)
		}
	}

	return nil
}

func (r *Registry) Get(name string) (Language, bool) {
	lang, ok := r.Languages[name]
	return lang, ok
}

func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.Languages))
	for name := range r.Languages {
		names = append(names, name)
	}
	return names
}
