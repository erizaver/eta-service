package app

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const defaultConfig = "values/values_local.yaml"

func init() {
	var localFile string
	for i, v := range os.Args {
		if v == "--local-config" && i+1 < len(os.Args) {
			localFile = os.Args[i+1]
			break
		}
		if strings.HasPrefix(v, "--local-config=") {
			parts := strings.SplitN(v, "=", 2)
			localFile = parts[1]
			break
		}
	}

	if localFile == "" {
		localFile = defaultConfig
	}
	if err := loadEnvFromValuesFile(localFile); err != nil {
		panic(err)
	}
}

func loadEnvFromValuesFile(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "unable to read yaml file")
	}

	var config valuesYamlConfig
	if err := yaml.Unmarshal(b, &config); err != nil {
		return errors.Wrap(err, "error unmarshaling file")
	}

	for _, v := range config.Env {
		if err := os.Setenv(v.Name, v.Value); err != nil {
			return errors.Wrapf(err, "set env %s='%s'", v.Name, v.Value)
		}
	}

	return nil
}

type valuesYamlConfig struct {
	Env []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"env"`
}
