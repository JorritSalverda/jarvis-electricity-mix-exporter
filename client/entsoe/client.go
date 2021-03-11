package entsoe

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sethgrid/pester"
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

	log.Info().Msgf("Getting aggregated generation per type for in domain %v and time interval %v...", area, timeInterval)

	getAggregatedGenerationPerTypeURL := fmt.Sprintf("%v?securityToken=%v&documentType=%v&processType=%v&in_Domain=%v&timeInterval=%v", c.apiBaseURL, c.securityToken, DocumentTypeActualGenerationPerType, ProcessTypeRealised, area, timeInterval.FormatAsParameter())

	resp, err := pester.Get(getAggregatedGenerationPerTypeURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Request returned unexpected status code %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	err = xml.Unmarshal(body, &response)
	if err != nil {
		return
	}

	return
}
