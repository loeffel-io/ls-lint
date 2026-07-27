package config

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

type configFile struct {
	Extends extendsList `yaml:"extends"`
	Ls      Ls          `yaml:"ls"`
	Ignore  []string    `yaml:"ignore"`
}

type extendsList []string

func (list *extendsList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		if value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
			return fmt.Errorf("extends must be a non-empty string or list of non-empty strings")
		}

		*list = extendsList{value.Value}
		return nil
	case yaml.SequenceNode:
		paths := make(extendsList, 0, len(value.Content))
		for _, item := range value.Content {
			if item.Kind != yaml.ScalarNode || item.Tag != "!!str" || strings.TrimSpace(item.Value) == "" {
				return fmt.Errorf("extends must be a non-empty string or list of non-empty strings")
			}

			paths = append(paths, item.Value)
		}

		*list = paths
		return nil
	default:
		return fmt.Errorf("extends must be a non-empty string or list of non-empty strings")
	}
}

type loader struct {
	active map[string]int
	stack  []string
}

// Load reads and merges config files in order. Each config may extend other
// config files, which are merged before the config that declares them.
func Load(paths []string) (*Config, error) {
	result := NewConfig(make(Ls), make([]string, 0))
	configLoader := loader{
		active: make(map[string]int),
		stack:  make([]string, 0),
	}

	for _, path := range paths {
		loaded, err := configLoader.load(path)
		if err != nil {
			return nil, err
		}

		merge(result, loaded)
	}

	slices.Sort(result.Ignore)
	result.Ignore = slices.Compact(result.Ignore)

	return result, nil
}

func (configLoader *loader) load(path string) (*Config, error) {
	absolutePath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("resolve config %q: %w", path, err)
	}

	identity := absolutePath
	if evaluatedPath, evaluateErr := filepath.EvalSymlinks(absolutePath); evaluateErr == nil {
		identity = evaluatedPath
	}

	if index, ok := configLoader.active[identity]; ok {
		chain := append(slices.Clone(configLoader.stack[index:]), absolutePath)
		return nil, fmt.Errorf("config extends cycle: %s", strings.Join(chain, " -> "))
	}

	configLoader.active[identity] = len(configLoader.stack)
	configLoader.stack = append(configLoader.stack, absolutePath)
	defer func() {
		delete(configLoader.active, identity)
		configLoader.stack = configLoader.stack[:len(configLoader.stack)-1]
	}()

	configBytes, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", absolutePath, err)
	}

	var parsed configFile
	if err = yaml.Unmarshal(configBytes, &parsed); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", absolutePath, err)
	}

	result := NewConfig(make(Ls), make([]string, 0))
	for _, extendedPath := range parsed.Extends {
		resolvedPath := extendedPath
		if !filepath.IsAbs(resolvedPath) {
			resolvedPath = filepath.Join(filepath.Dir(absolutePath), resolvedPath)
		}

		extendedConfig, loadErr := configLoader.load(resolvedPath)
		if loadErr != nil {
			return nil, fmt.Errorf("config %q extends %q: %w", absolutePath, extendedPath, loadErr)
		}

		merge(result, extendedConfig)
	}

	merge(result, NewConfig(parsed.Ls, parsed.Ignore))

	return result, nil
}

func merge(target *Config, source *Config) {
	maps.Copy(target.Ls, source.Ls)
	target.Ignore = append(target.Ignore, source.Ignore...)
}
