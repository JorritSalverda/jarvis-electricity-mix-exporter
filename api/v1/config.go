package api

import (
	"fmt"
)

type Config struct {
	Areas []*AreaConfig `yaml:"areas"`
}

type AreaConfig struct {
	Area              Area              `yaml:"area"`
	Country           CountryCode       `yaml:"country"`
	Source            Source            `yaml:"source"`
	ResolutionMinutes int               `yaml:"resolutionMinutes"`
	StartYearsAgo     int               `yaml:"startYearsAgo"`
	StartMonthsAgo    int               `yaml:"startMonthsAgo"`
	StartDaysAgo      int               `yaml:"startDaysAgo"`
	Exchanges         []*ExchangeConfig `yaml:"exchanges"`
}

type ExchangeConfig struct {
	Area              Area        `yaml:"area"`
	Country           CountryCode `yaml:"country"`
	Source            Source      `yaml:"source"`
	ResolutionMinutes int         `yaml:"resolutionMinutes"`
}

func (c *Config) SetDefaults() {
	for _, a := range c.Areas {
		a.SetDefaults()
	}
}

func (ac *AreaConfig) SetDefaults() {
	if ac.Source == SourceUnknown {
		ac.Source = SourceEntsoe
	}
	if ac.ResolutionMinutes == 0 {
		ac.ResolutionMinutes = 15
	}
	for _, e := range ac.Exchanges {
		e.SetDefaults()
	}
}

func (ec *ExchangeConfig) SetDefaults() {
	if ec.Source == SourceUnknown {
		ec.Source = SourceEntsoe
	}
	if ec.ResolutionMinutes == 0 {
		ec.ResolutionMinutes = 60
	}
}

func (c *Config) Validate() (valid bool, errors []error, warnings []string) {
	if len(c.Areas) == 0 {
		errors = append(errors, fmt.Errorf("No areas have been configured, set at least one area"))
	}
	for _, a := range c.Areas {
		e, w := a.validate()
		errors = append(errors, e...)
		warnings = append(warnings, w...)
	}

	return len(errors) == 0, errors, warnings
}

func (ac *AreaConfig) validate() (errors []error, warnings []string) {
	if ac.Source == SourceUnknown {
		errors = append(errors, fmt.Errorf("Source for area is unknown, set with `source: entsoe`"))
	}
	if ac.Country == CountryCodeUnknown {
		errors = append(errors, fmt.Errorf("Country for area is unknown, set with `country: NL`"))
	}
	if ac.ResolutionMinutes == 0 {
		errors = append(errors, fmt.Errorf("Resolution for area is unknown, set with `resolutionMinutes: 15`"))
	}
	for _, e := range ac.Exchanges {
		er, w := e.validate()
		errors = append(errors, er...)
		warnings = append(warnings, w...)
	}

	return errors, warnings
}

func (ec *ExchangeConfig) validate() (errors []error, warnings []string) {
	if ec.Source == SourceUnknown {
		errors = append(errors, fmt.Errorf("Source for area is unknown, set with `source: entsoe`"))
	}
	if ec.Country == CountryCodeUnknown {
		errors = append(errors, fmt.Errorf("Country for area is unknown, set with `country: NL`"))
	}
	if ec.ResolutionMinutes == 0 {
		errors = append(errors, fmt.Errorf("Resolution for area is unknown, set with `resolutionMinutes: 15`"))
	}

	return errors, warnings
}

type Source string

const (
	SourceUnknown Source = ""
	SourceEntsoe  Source = "entsoe"
)

type CountryCode string

const (
	CountryCodeUnknown      CountryCode = ""
	CountryCodeNetherlands  CountryCode = "NL"
	CountryCodeBelgium      CountryCode = "BE"
	CountryCodeGermany      CountryCode = "DE"
	CountryCodeDenmark      CountryCode = "DK"
	CountryCodeGreatBritain CountryCode = "GB"
	CountryCodeNorway       CountryCode = "NO"
)
