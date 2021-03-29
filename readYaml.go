package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type StagesYaml = map[string]interface{}
type BasicYaml = map[string]interface{}

type Config struct {
	watch    []string
	engine   string
	debounce int
}

type Actions = map[string]([]string)

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

var defaultConfig = Config{
	watch: []string{
		"src/**/*",
		"lib/**/*",
		"app/**/*",
	},
	engine:   "/bin/sh",
	debounce: 2000,
}

func parseStringArray(val interface{}) ([]string, error) {
	array, ok := val.([]interface{})
	if ok {
		var strings []string
		for i, item := range array {
			str, ok := item.(string)
			if ok {
				strings = append(strings, str)
			} else {
				return nil, fmt.Errorf("Bad formatting at action number \"%v\"", i)
			}
		}
		return strings, nil
	} else {
		return nil, fmt.Errorf("Bad formatting.")
	}
}

func parseEngine(engStr interface{}) string {
	engine, ok := engStr.(string)
	if ok {
		return engine
	} else {
		return defaultConfig.engine
	}
}

func parseWatch(watchStr interface{}) []string {
	watchPattern, err := parseStringArray(watchStr)
	if err == nil {
		return watchPattern
	} else {
		fmt.Println("\033[31m", fmt.Sprintf("Watch pattern parsing error due to '%v'", err))
		return defaultConfig.watch
	}
}

func parseDebounce(debounceObj interface{}) int {
	debounce, ok := debounceObj.(int)
	if ok {
		return debounce
	} else {
		return defaultConfig.debounce
	}
}

func parseConfig(yamlData StagesYaml) Config {
	configYaml, ok := yamlData["_config"].(map[interface{}]interface{})
	if ok {
		return Config{
			watch:    parseWatch(configYaml["watch"]),
			engine:   parseEngine(configYaml["engine"]),
			debounce: parseDebounce(configYaml["debounce"]),
		}
	} else {
		return defaultConfig
	}
}

func parseActions(yamlData StagesYaml) Actions {
	delete(yamlData, "_config")
	actions := make(map[string]([]string))

	for key, val := range yamlData {
		steps, err := parseStringArray(val)
		if err == nil {
			actions[key] = steps
		} else {
			fmt.Println("\033[31m", fmt.Sprintf("Action \"key\" parsing error due: ", err.Error()))
		}
	}

	return actions
}

func ReadYaml(yamlPath string) (Config, Actions) {
	yamlFile := parseYaml(
		readYamlFile(yamlPath),
	)
	return parseConfig(yamlFile), parseActions(yamlFile)
}
