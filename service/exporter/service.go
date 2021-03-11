package exporter

import (
	"context"
	"time"

	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/bigquery"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/entsoe"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/state"
	contractsv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/contracts/v1"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Run(ctx context.Context, area entsoe.Area) error
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

func (s *service) Run(ctx context.Context, area entsoe.Area) error {

	// check if there's a previous measurement stored in state file
	lastMeasurement, err := s.stateClient.ReadState(ctx)
	if err != nil {
		return err
	}

	for {
		now := time.Now().UTC().Round(time.Duration(entsoe.TimeSlotsInMinutes * time.Minute))

		// if it's the first time begin a year ago, otherwise start at last stored value
		var start time.Time
		if lastMeasurement != nil {
			start = lastMeasurement.MeasuredAtTime.Add(time.Duration(entsoe.TimeSlotsInMinutes) * time.Minute)
		} else {
			start = now.AddDate(-1, 0, 0)
		}
		end := start.Add(time.Duration(4*24*entsoe.TimeSlotsInMinutes) * time.Minute)
		if end.After(now) {
			end = now
		}

		// don't continue, we're up to date
		if !start.Before(end) {
			log.Info().Msgf("Start - %v - and end - %v - are equal, exiting", start, end)
			return nil
		}

		// retrieve actual measurements
		response, err := s.entsoeClient.GetAggregatedGenerationPerType(area, entsoe.TimeInterval{
			Start: start,
			End:   end,
		})
		if err != nil {
			return err
		}

		if len(response.TimeSeries) == 0 {
			log.Info().Msg("No timeseries have been returned, exiting")
			return nil
		}

		nrOfSlots := int(response.TimePeriod.End.Sub(response.TimePeriod.Start).Minutes() / entsoe.TimeSlotsInMinutes)

		var lastStoredMeasurement *contractsv1.Measurement
		for i := 0; i < nrOfSlots; i++ {
			timeSlotStartTime := response.TimePeriod.Start.Add(time.Duration(i*entsoe.TimeSlotsInMinutes) * time.Minute)
			if lastMeasurement != nil && timeSlotStartTime.Equal(lastMeasurement.MeasuredAtTime) {
				log.Info().Msgf("Timeslot start at %v has already been stored, continuing to next timeslot", timeSlotStartTime)
				continue
			}

			measurement := s.createMeasurementForTimeSlot(response, timeSlotStartTime, area)

			// store measurement
			err = s.bigqueryClient.InsertMeasurement(measurement)
			if err != nil {
				return err
			}

			lastStoredMeasurement = &measurement
		}

		if lastStoredMeasurement != nil {
			// store state
			err = s.stateClient.StoreState(ctx, *lastStoredMeasurement)
			if err != nil {
				return err
			}

			lastMeasurement = lastStoredMeasurement
		} else {
			log.Info().Msg("No new measurements were inserted, exiting")
			return nil
		}

		// otherwise wait a bit before starting next loop iteration to avoid hitting rate limits
		log.Info().Msg("Sleeping for 15 seconds before retrieving more data, to avoid rate limiting")
		time.Sleep(time.Duration(15) * time.Second)
	}
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

func (s *service) mapToSampleDirection(timeSerie entsoe.TimeSerie) contractsv1.SampleDirection {
	if timeSerie.InBiddingZone != entsoe.AreaUnknown {
		return contractsv1.SampleDirectionIn
	}
	if timeSerie.OutBiddingZone != entsoe.AreaUnknown {
		return contractsv1.SampleDirectionOut
	}

	return contractsv1.SampleDirectionUnknown
}

func (s *service) mapToSampleUnit(measurementUnit entsoe.MeasurementUnit) contractsv1.SampleUnit {
	if measurementUnit == entsoe.MeasurementUnitMegaWatt {
		return contractsv1.SampleUnitMegaWatt
	}

	return contractsv1.SampleUnitUnknown
}

func (s *service) createMeasurementForTimeSlot(response entsoe.GetAggregatedGenerationPerTypeResponse, timeSlotStartTime time.Time, area entsoe.Area) contractsv1.Measurement {

	measurement := contractsv1.Measurement{
		ID:             uuid.New().String(),
		Source:         "ENTSOE",
		Area:           area.GetCountryCode(),
		MeasuredAtTime: timeSlotStartTime,
	}

	// insert all periods that started after last inserted one
	for _, ts := range response.TimeSeries {
		if ts.Period.TimeInterval.Start.After(timeSlotStartTime) {
			log.Info().Msgf("Timeserie %v for psr type %v starts after time slot %v, continuing to next timeserie", ts.ID, ts.MktPsrType.PsrType, timeSlotStartTime)
			continue
		}
		if ts.Period.TimeInterval.End.Equal(timeSlotStartTime) {
			log.Info().Msgf("Timeserie %v for psr type %v ends at time slot %v, continuing to next timeserie", ts.ID, ts.MktPsrType.PsrType, timeSlotStartTime)
			continue
		}
		if ts.Period.TimeInterval.End.Before(timeSlotStartTime) {
			log.Info().Msgf("Timeserie %v for psr type %v ends before time slot %v, continuing to next timeserie", ts.ID, ts.MktPsrType.PsrType, timeSlotStartTime)
			continue
		}

		pointIndexForSlot := int(timeSlotStartTime.Sub(ts.Period.TimeInterval.Start).Minutes() / entsoe.TimeSlotsInMinutes)

		energyType := s.mapToEnergyType(ts.MktPsrType.PsrType)
		if pointIndexForSlot < len(ts.Period.Points) {
			measurement.Samples = append(measurement.Samples, &contractsv1.Sample{
				EnergyType:         energyType,
				OriginalEnergyType: string(ts.MktPsrType.PsrType),
				IsRenewable:        energyType.IsRenewable(),
				MetricType:         contractsv1.MetricTypeGauge,
				SampleDirection:    s.mapToSampleDirection(ts),
				SampleUnit:         s.mapToSampleUnit(ts.QuanityMeasurementUnit),
				Value:              ts.Period.Points[pointIndexForSlot].Quantity,
			})
		} else {
			// this timeserie seems to have less points, what to do now?
			log.Warn().Msgf("Timeserie %v for psr type %v only has %v points, while index %v should be retrieved", ts.ID, ts.MktPsrType.PsrType, len(ts.Period.Points), pointIndexForSlot)
		}
	}

	return measurement
}
