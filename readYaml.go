package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type StagesYaml = map[string]([]string)

func GetConfigPath(rootPath string) string {
	return filepath.Join(rootPath, "stages.yaml")
}

func readYamlFile(yamlPath string) []byte {
	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		Crashf(
			"Can not read stages.yaml due to: %v",
			err.Error(),
		)
	}
	return data
}

func parseYaml(yamlData []byte) StagesYaml {
	yamlMap := StagesYaml{}

	err := yaml.Unmarshal(yamlData, &yamlMap)
	if err != nil {
		Crashf(
			"Can not decode stages.yaml file due to: %v",
			err.Error(),
		)
	}

	return yamlMap
}

func ReadYaml(yamlPath string) StagesYaml {
	return parseYaml(
		readYamlFile(yamlPath),
	)
}
