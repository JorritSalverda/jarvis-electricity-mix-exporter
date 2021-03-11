package exporter

import (
	"context"
	"time"

	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/bigquery"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/entsoe"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/state"
	contractsv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/contracts/v1"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Run(ctx context.Context, area entsoe.Area) (err error)
}

func NewService(bigqueryClient bigquery.Client, stateClient state.Client, entsoeClient entsoe.Client) (Service, error) {
	return &service{
		bigqueryClient: bigqueryClient,
		stateClient:    stateClient,
		entsoeClient:   entsoeClient,
	}, nil
}

type service struct {
	bigqueryClient bigquery.Client
	stateClient    state.Client
	entsoeClient   entsoe.Client
}

func (s *service) Run(ctx context.Context, area entsoe.Area) (err error) {

	// check if there's a previous measurement stored in state file
	lastMeasurement, err := s.stateClient.ReadState(ctx)
	if err != nil {
		return
	}

	// retrieve data from entsoe for the last 3 hours
	now := time.Now().UTC()
	nowRoundedToTimeSlotSize := now.Round(time.Duration(entsoe.TimeSlotsInMinutes * time.Minute))
	startTime := nowRoundedToTimeSlotSize.Add(time.Duration(-3 * time.Hour))
	endTime := nowRoundedToTimeSlotSize

	timeInterval := entsoe.TimeInterval{
		Start: startTime,
		End:   endTime,
	}

	response, err := s.entsoeClient.GetAggregatedGenerationPerType(area, timeInterval)
	if err != nil {
		return
	}

	// check which data is new and can be stored in bigquery
	startLastPeriod := response.TimePeriod.End.Add(-1 * entsoe.TimeSlotsInMinutes * time.Minute)
	if lastMeasurement != nil && !lastMeasurement.MeasuredAtTime.Before(startLastPeriod) {
		log.Info().Msgf("Last measurement was stored at %v and last retrieved period starts at %v, exiting", lastMeasurement.MeasuredAtTime, startLastPeriod)
		return nil
	}

	if len(response.TimeSeries) == 0 || len(response.TimeSeries[0].Period.Points) == 0 {
		log.Info().Msg("Not timeseries or points have been returned, exiting")
		return nil
	}

	nrOfPoints := len(response.TimeSeries[0].Period.Points)

	for i := 0; i < nrOfPoints; i++ {
		measurement := contractsv1.Measurement{
			Source:         "ENTSOE",
			Area:           area.GetCountryCode(),
			MeasuredAtTime: response.TimePeriod.Start.Add(time.Duration(i*entsoe.TimeSlotsInMinutes) * time.Minute),
		}

		// insert all periods that started after last inserted one
		for _, ts := range response.TimeSeries {
			if i < len(ts.Period.Points) {

				measurement.Samples = append(measurement.Samples, &contractsv1.Sample{
					EnergyType: s.mapToEnergyType(ts.MktPsrType.PsrType),
				})

				// ts.Period.Points[i].Quantity
			}

			// store measurement
			err = s.bigqueryClient.InsertMeasurement(measurement)
			if err != nil {
				return
			}
		}
	}

	// // store state
	// err = s.stateClient.StoreState(ctx, measurement)
	// if err != nil {
	// 	return
	// }

	return nil
}

func (s *service) mapToEnergyType(psrType entsoe.PsrType) contractsv1.EnergyType {

	switch psrType {
	case entsoe.PsrTypeBiomass:
		return contractsv1.EnergyTypeBiomass

	case entsoe.PsrTypeFossilHardCoal,
		entsoe.PsrTypeFossilBrownCoal:
		return contractsv1.EnergyTypeCoal

	case entsoe.PsrTypeFossilCoalDerivedGas,
		entsoe.PsrTypeFossilGas:
		return contractsv1.EnergyTypeGas

	case entsoe.PsrTypeFossilOil,
		entsoe.PsrTypeFossilOilShale,
		entsoe.PsrTypeFossilOilPeat:
		return contractsv1.EnergyTypeOil

	case entsoe.PsrTypeGeothermal:
		return contractsv1.EnergyTypeGeothermal

	case entsoe.PsrTypeHydroPumpedStorage,
		entsoe.PsrTypeHydroRunOfRiver,
		entsoe.PsrTypeHydroWaterReservoir,
		entsoe.PsrTypeMarin:
		return contractsv1.EnergyTypeHydro

	case entsoe.PsrTypeNuclear:
		return contractsv1.EnergyTypeNuclear

	case entsoe.PsrTypeOtherRenewable:
		return contractsv1.EnergyTypeOtherRenewable

	case entsoe.PsrTypeSolar:
		return contractsv1.EnergyTypeSolar
	case entsoe.PsrTypeWaste:
		return contractsv1.EnergyTypeWaste

	case entsoe.PsrTypeWindOffshore:
		return contractsv1.EnergyTypeWindOffshore

	case entsoe.PsrTypeWindOnshore:
		return contractsv1.EnergyTypeWindOnshore
	}

	return contractsv1.EnergyTypeUnknown
}
