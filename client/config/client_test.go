package config

import (
	"testing"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFromFile(t *testing.T) {

	t.Run("ReturnsConfig", func(t *testing.T) {

		client, _ := NewClient("./test-config.yaml")

		// act
		config, err := client.ReadConfigFromFile("./test-config.yaml")

		assert.Nil(t, err)
		assert.Equal(t, 6, len(config.Areas))
		assert.Equal(t, apiv1.AreaNetherlands, config.Areas[0].Area)
		assert.Equal(t, apiv1.CountryCodeNetherlands, config.Areas[0].Country)
		assert.Equal(t, 1, config.Areas[0].StartYearsAgo)
		assert.Equal(t, 2, config.Areas[0].StartMonthsAgo)
		assert.Equal(t, 3, config.Areas[0].StartDaysAgo)
		assert.Equal(t, 15, config.Areas[0].ResolutionMinutes)

		assert.Equal(t, 5, len(config.Areas[0].Exchanges))
		assert.Equal(t, apiv1.AreaBelgium, config.Areas[0].Exchanges[0].Area)
		assert.Equal(t, apiv1.CountryCodeBelgium, config.Areas[0].Exchanges[0].Country)
		assert.Equal(t, 60, config.Areas[0].Exchanges[0].ResolutionMinutes)
	})
}

func TestReadConfig(t *testing.T) {

	t.Run("ReturnsConfig", func(t *testing.T) {

		client, _ := NewClient("./test-config.yaml")

		// act
		config, err := client.ReadConfig()

		assert.Nil(t, err)
		assert.Equal(t, 6, len(config.Areas))
		assert.Equal(t, apiv1.AreaNetherlands, config.Areas[0].Area)
		assert.Equal(t, apiv1.CountryCodeNetherlands, config.Areas[0].Country)
		assert.Equal(t, 15, config.Areas[0].ResolutionMinutes)

		assert.Equal(t, 5, len(config.Areas[0].Exchanges))
		assert.Equal(t, apiv1.AreaBelgium, config.Areas[0].Exchanges[0].Area)
		assert.Equal(t, apiv1.CountryCodeBelgium, config.Areas[0].Exchanges[0].Country)
		assert.Equal(t, 60, config.Areas[0].Exchanges[0].ResolutionMinutes)
	})
}
