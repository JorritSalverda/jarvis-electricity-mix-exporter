package exporter

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/bigquery"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/config"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/entsoe"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/state"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Run(ctx context.Context, gracefulShutdown chan os.Signal, waitGroup *sync.WaitGroup) error
}

func NewService(generationBigqueryClient bigquery.Client, exchangeBigqueryClient bigquery.Client, configClient config.Client, stateClient state.Client, entsoeClient entsoe.Client) (Service, error) {
	return &service{
		generationBigqueryClient: generationBigqueryClient,
		exchangeBigqueryClient:   exchangeBigqueryClient,
		configClient:             configClient,
		stateClient:              stateClient,
		entsoeClient:             entsoeClient,
	}, nil
}

type service struct {
	generationBigqueryClient bigquery.Client
	exchangeBigqueryClient   bigquery.Client
	configClient             config.Client
	stateClient              state.Client
	entsoeClient             entsoe.Client
}

func (s *service) Run(ctx context.Context, gracefulShutdown chan os.Signal, waitGroup *sync.WaitGroup) error {

	config, err := s.configClient.ReadConfig()
	if err != nil {
		return err
	}

	// check if there's a previous measurement stored in state file
	lastState, err := s.stateClient.ReadState(ctx)
	if err != nil {
		return err
	}

	for _, areaConfig := range config.Areas {
		err = s.runForArea(ctx, gracefulShutdown, waitGroup, *areaConfig, lastState)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) runForArea(ctx context.Context, gracefulShutdown chan os.Signal, waitGroup *sync.WaitGroup, areaConfig apiv1.AreaConfig, lastState *apiv1.State) error {

	log.Info().Interface("areaConfig", areaConfig).Msgf("Retrieving measurements for area %v / country %v", areaConfig.Area, areaConfig.Country)

	for {
		now := time.Now().UTC().Round(time.Duration(areaConfig.ResolutionMinutes) * time.Minute)

		// if it's the first time begin a year ago, otherwise start at last stored value
		start := now.AddDate(-1*areaConfig.StartYearsAgo, -1*areaConfig.StartMonthsAgo, -1*areaConfig.StartDaysAgo)
		if lastState != nil && lastState.LastRetrievedGenerationTime != nil {
			if lastRetrievedGenerationTime, ok := lastState.LastRetrievedGenerationTime[areaConfig.Area]; ok {
				start = lastRetrievedGenerationTime.Add(time.Duration(areaConfig.ResolutionMinutes) * time.Minute)
			}
		}
		end := start.Add(time.Duration(4*24*areaConfig.ResolutionMinutes) * time.Minute)
		if end.After(now) {
			end = now
		}

		// don't continue, we're up to date
		if !start.Before(end) {
			log.Info().Msgf("Start - %v - and end - %v - are equal, exiting", start, end)
			return nil
		}

		// retrieve actual measurements
		response, err := s.entsoeClient.GetAggregatedGenerationPerType(areaConfig.Area, apiv1.TimeInterval{
			Start: start,
			End:   end,
		})
		if err != nil && !errors.Is(err, entsoe.ErrNoMatchingDataFound) {
			return err
		}
		if err != nil && errors.Is(err, entsoe.ErrNoMatchingDataFound) {
			log.Info().Msg("No data has been found, exiting")
			return nil
		}

		if len(response.TimeSeries) == 0 {
			log.Info().Msg("No timeseries have been returned, exiting")
			return nil
		}

		nrOfSlots := int(response.TimePeriod.End.Sub(response.TimePeriod.Start).Minutes() / float64(areaConfig.ResolutionMinutes))

		if nrOfSlots == 0 {
			log.Info().Msg("No new measurements were inserted, exiting")
			return nil
		}

		for i := 0; i < nrOfSlots; i++ {
			err := s.handleTimeSlot(ctx, response, i, areaConfig, waitGroup, lastState)
			if err != nil {
				return err
			}
		}

		log.Info().Msg("Sleeping for 15 seconds before retrieving more data, to avoid rate limiting")
		select {
		case signalReceived := <-gracefulShutdown:
			log.Warn().Msgf("Received signal %v. Waiting for running tasks to finish...", signalReceived)
			return nil
		case <-time.After(15 * time.Second):
		}
	}
}

func (s *service) mapToEnergyType(psrType apiv1.PsrType) apiv1.EnergyType {

	switch psrType {
	case apiv1.PsrTypeBiomass:
		return apiv1.EnergyTypeBiomass

	case apiv1.PsrTypeFossilHardCoal,
		apiv1.PsrTypeFossilBrownCoal:
		return apiv1.EnergyTypeCoal

	case apiv1.PsrTypeFossilCoalDerivedGas,
		apiv1.PsrTypeFossilGas:
		return apiv1.EnergyTypeGas

	case apiv1.PsrTypeFossilOil,
		apiv1.PsrTypeFossilOilShale,
		apiv1.PsrTypeFossilOilPeat:
		return apiv1.EnergyTypeOil

	case apiv1.PsrTypeGeothermal:
		return apiv1.EnergyTypeGeothermal

	case apiv1.PsrTypeHydroPumpedStorage,
		apiv1.PsrTypeHydroRunOfRiver,
		apiv1.PsrTypeHydroWaterReservoir,
		apiv1.PsrTypeMarin:
		return apiv1.EnergyTypeHydro

	case apiv1.PsrTypeNuclear:
		return apiv1.EnergyTypeNuclear

	case apiv1.PsrTypeOtherRenewable:
		return apiv1.EnergyTypeOtherRenewable

	case apiv1.PsrTypeSolar:
		return apiv1.EnergyTypeSolar
	case apiv1.PsrTypeWaste:
		return apiv1.EnergyTypeWaste

	case apiv1.PsrTypeWindOffshore:
		return apiv1.EnergyTypeWindOffshore

	case apiv1.PsrTypeWindOnshore:
		return apiv1.EnergyTypeWindOnshore
	}

	return apiv1.EnergyTypeUnknown
}

func (s *service) mapToSampleDirection(timeSerie apiv1.AggregatedGenerationTimeSerie) apiv1.SampleDirection {
	if timeSerie.InBiddingZone != apiv1.AreaUnknown {
		return apiv1.SampleDirectionIn
	}
	if timeSerie.OutBiddingZone != apiv1.AreaUnknown {
		return apiv1.SampleDirectionOut
	}

	return apiv1.SampleDirectionUnknown
}

func (s *service) mapToSampleUnit(measurementUnit apiv1.MeasurementUnit) apiv1.SampleUnit {
	if measurementUnit == apiv1.MeasurementUnitMegaWatt {
		return apiv1.SampleUnitMegaWatt
	}

	return apiv1.SampleUnitUnknown
}

func (s *service) createGenerationMeasurementForTimeSlot(response apiv1.GetAggregatedGenerationPerTypeResponse, timeSlotStartTime time.Time, areaConfig apiv1.AreaConfig) apiv1.GenerationMeasurement {
	measurement := apiv1.GenerationMeasurement{
		ID:             uuid.New().String(),
		Source:         "ENTSOE",
		Area:           areaConfig.Area.GetCountryCode(),
		MeasuredAtTime: timeSlotStartTime,
	}

	// insert all periods that started after last inserted one
	for _, ts := range response.TimeSeries {
		if ts.Period.TimeInterval.Start.After(timeSlotStartTime) {
			// log.Info().Msgf("Timeserie %v for psr type %v starts after time slot %v, continuing to next timeserie", ts.ID, ts.MktPsrType.PsrType, timeSlotStartTime)
			continue
		}
		if ts.Period.TimeInterval.End.Equal(timeSlotStartTime) {
			// log.Info().Msgf("Timeserie %v for psr type %v ends at time slot %v, continuing to next timeserie", ts.ID, ts.MktPsrType.PsrType, timeSlotStartTime)
			continue
		}
		if ts.Period.TimeInterval.End.Before(timeSlotStartTime) {
			// log.Info().Msgf("Timeserie %v for psr type %v ends before time slot %v, continuing to next timeserie", ts.ID, ts.MktPsrType.PsrType, timeSlotStartTime)
			continue
		}

		pointIndexForSlot := int(timeSlotStartTime.Sub(ts.Period.TimeInterval.Start).Minutes() / float64(areaConfig.ResolutionMinutes))

		energyType := s.mapToEnergyType(ts.MktPsrType.PsrType)
		if pointIndexForSlot < len(ts.Period.Points) {
			measurement.Samples = append(measurement.Samples, &apiv1.Sample{
				EnergyType:         energyType,
				OriginalEnergyType: string(ts.MktPsrType.PsrType),
				IsRenewable:        energyType.IsRenewable(),
				MetricType:         apiv1.MetricTypeGauge,
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

func (s *service) handleTimeSlot(ctx context.Context, response apiv1.GetAggregatedGenerationPerTypeResponse, timeSlotIndex int, areaConfig apiv1.AreaConfig, waitGroup *sync.WaitGroup, lastState *apiv1.State) (err error) {
	timeSlotStartTime := response.TimePeriod.Start.Add(time.Duration(timeSlotIndex*areaConfig.ResolutionMinutes) * time.Minute)

	measurement := s.createGenerationMeasurementForTimeSlot(response, timeSlotStartTime, areaConfig)

	// store measurement
	waitGroup.Add(1)
	defer waitGroup.Done()
	err = s.generationBigqueryClient.InsertMeasurement(measurement)
	if err != nil {
		return
	}

	// update state
	if lastState == nil {
		lastState = &apiv1.State{}
	}
	if lastState.LastRetrievedGenerationTime == nil {
		lastState.LastRetrievedGenerationTime = make(map[apiv1.Area]time.Time, 0)
	}
	lastState.LastRetrievedGenerationTime[areaConfig.Area] = measurement.MeasuredAtTime

	// store state
	err = s.stateClient.StoreState(ctx, *lastState)
	if err != nil {
		return
	}

	return
}
