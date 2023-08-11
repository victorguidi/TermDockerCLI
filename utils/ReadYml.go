package utils

import (
	"os"

	"github.com/victorguidi/TermDockerCLI/types"
	"gopkg.in/yaml.v3"
)

func ReadYml(filename string) (types.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return types.Config{}, err
	}
	defer file.Close()

	var config types.Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return types.Config{}, err
	}

	return config, nil
}
