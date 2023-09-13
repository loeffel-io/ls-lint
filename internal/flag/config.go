package flag

import "strings"

type Config []string

func (config *Config) String() string {
	return strings.Join(*config, ",")
}

func (config *Config) Set(value string) error {
	*config = append(*config, value)
	return nil
}
