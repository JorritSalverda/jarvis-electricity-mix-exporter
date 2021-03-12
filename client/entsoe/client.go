package entsoe

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sethgrid/pester"
)

var (
	ErrNoMatchingDataFound = errors.New("No matching data found")
)

type Client interface {
	GetAggregatedGenerationPerType(area Area, timeInterval TimeInterval) (response GetAggregatedGenerationPerTypeResponse, err error)
}

func NewClient(securityToken string) (Client, error) {
	if securityToken == "" {
		return nil, fmt.Errorf("Token is empty, please provide a valid api token for transparency.entsoe.eu")
	}

	return &client{
		apiBaseURL:    "https://transparency.entsoe.eu/api",
		securityToken: securityToken,
	}, nil
}

type client struct {
	apiBaseURL    string
	securityToken string
}

func (c *client) GetAggregatedGenerationPerType(area Area, timeInterval TimeInterval) (response GetAggregatedGenerationPerTypeResponse, err error) {

	// https://transparency.entsoe.eu/content/static_content/Static%20content/web%20api/Guide.html#_aggregated_generation_per_type_16_1_b_c

	// 4.4.8. Aggregated Generation per Type
	// - One year range limit applies
	// - Minimum time interval in query response is one MTU period
	// - Mandatory parameters
	// 	 - DocumentType
	// 	 - ProcessType
	// 	 - In_Domain
	// 	 - TimeInterval or combination of PeriodStart and PeriodEnd
	// - Optional parameters
	//   - PsrType (When used, only queried production type is returned)

	// Please note the followings:
	// - Response from API is same irrespective of querying for Document Types A74 - Wind & Solar & A75 - Actual  Generation Per Type
	// - Time series with inBiddingZone_Domain attribute reflects Generation values while outBiddingZone_Domain reflects Consumption values.

	log.Info().Msgf("Getting aggregated generation per type for in domain %v and time interval %v to %v...", area, timeInterval.Start, timeInterval.End)

	getAggregatedGenerationPerTypeURL := fmt.Sprintf("%v?securityToken=%v&documentType=%v&processType=%v&in_Domain=%v&timeInterval=%v", c.apiBaseURL, c.securityToken, DocumentTypeActualGenerationPerType, ProcessTypeRealised, area, timeInterval.FormatAsParameter())

	log.Debug().Msgf("GET %v", strings.Replace(getAggregatedGenerationPerTypeURL, c.securityToken, "***", -1))

	resp, err := pester.Get(getAggregatedGenerationPerTypeURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {

		log.Debug().Str("body", string(body)).Msgf("%v GET %v", resp.StatusCode, strings.Replace(getAggregatedGenerationPerTypeURL, c.securityToken, "***", -1))

		if resp.StatusCode == http.StatusBadRequest && strings.Contains(string(body), "No matching data found") {
			return response, ErrNoMatchingDataFound
		}

		return response, fmt.Errorf("Request returned unexpected status code %v", resp.StatusCode)
	}

	err = xml.Unmarshal(body, &response)
	if err != nil {
		return
	}

	return
}
