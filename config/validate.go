package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func ValidateConfigFile(file string) ([]toml.Key, error) {
	config := &Config{}
	md, err := toml.DecodeFile(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to test config: %s", err)
	}
	return md.Undecoded(), nil
}
