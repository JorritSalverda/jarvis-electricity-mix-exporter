package config

import (
	"errors"
	"io/ioutil"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var (
	ErrConfigNotValid = errors.New("Config is invalid")
)

type Client interface {
	ReadConfig() (config apiv1.Config, err error)
	ReadConfigFromFile(path string) (config apiv1.Config, err error)
}

func NewClient(configPath string) (Client, error) {
	return &client{
		configPath: configPath,
	}, nil
}

type client struct {
	configPath string
}

func (c *client) ReadConfig() (config apiv1.Config, err error) {
	return c.ReadConfigFromFile(c.configPath)
}

func (c *client) ReadConfigFromFile(path string) (config apiv1.Config, err error) {
	log.Debug().Msgf("Reading %v file...", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return config, err
	}

	// ensure all defaults are set
	config.SetDefaults()

	// validate config
	valid, errors, warnings := config.Validate()

	if len(warnings) > 0 {
		log.Warn().Interface("warnings", warnings).Msgf("Config file at %v has warnings", path)
	}

	if !valid {
		log.Warn().Interface("errors", errors).Msgf("Config file at %v is not valid, it has errors", path)
		return config, ErrConfigNotValid
	}

	return
}
