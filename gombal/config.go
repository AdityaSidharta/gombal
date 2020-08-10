package gombal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Strategy string `yaml:"Strategy"`
}

func LoadConfig(path string) (Config, error) {
	c := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	logrus.Info(fmt.Sprintf("Config successfully loaded from %v", path))
	logrus.Info(fmt.Sprintf("Config : %v", c))

	_, ok := Find(supportedStrategy, c.Strategy)
	if !ok {
		return c, invalidStrategyError
	}

	return c, nil
}
