package main

import (
	"context"
	"runtime"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/bigquery"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/config"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/entsoe"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/state"
	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/service/exporter"
	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	// set when building the application
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()

	// application specific config
	entsoeToken = kingpin.Flag("entsoe-token", "Api token for https://transparency.entsoe.eu/api").Envar("ENTSOE_TOKEN").Required().String()

	bigqueryEnable          = kingpin.Flag("bigquery-enable", "Toggle to enable or disable bigquery integration").Default("true").OverrideDefaultFromEnvar("BQ_ENABLE").Bool()
	bigqueryInit            = kingpin.Flag("bigquery-init", "Toggle to enable bigquery table initialization").Default("true").OverrideDefaultFromEnvar("BQ_INIT").Bool()
	bigqueryProjectID       = kingpin.Flag("bigquery-project-id", "Google Cloud project id that contains the BigQuery dataset").Envar("BQ_PROJECT_ID").Required().String()
	bigqueryDataset         = kingpin.Flag("bigquery-dataset", "Name of the BigQuery dataset").Envar("BQ_DATASET").Required().String()
	bigqueryGenerationTable = kingpin.Flag("bigquery-generation-table", "Name of the BigQuery table with generation measurements").Envar("BQ_GENERATION_TABLE").Required().String()
	bigqueryExchangeTable   = kingpin.Flag("bigquery-exchange-table", "Name of the BigQuery table with generation measurements").Envar("BQ_EXCHANGE_TABLE").Required().String()

	configPath                   = kingpin.Flag("config-path", "Path to the config.yaml file").Default("/configs/config.yaml").OverrideDefaultFromEnvar("CONFIG_PATH").String()
	measurementFilePath          = kingpin.Flag("state-file-path", "Path to file with state.").Default("/configs/last-measurement.json").OverrideDefaultFromEnvar("MEASUREMENT_FILE_PATH").String()
	measurementFileConfigMapName = kingpin.Flag("state-file-configmap-name", "Name of the configmap with state file.").Default("jarvis-electricity-mix-exporter").OverrideDefaultFromEnvar("MEASUREMENT_FILE_CONFIG_MAP_NAME").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(foundation.NewApplicationInfo(appgroup, app, version, branch, revision, buildDate))

	// create context to cancel commands on sigterm
	ctx := foundation.InitCancellationContext(context.Background())

	// init bigquery client
	generationBigqueryClient, err := bigquery.NewClient(*bigqueryProjectID, *bigqueryEnable, *bigqueryDataset, *bigqueryGenerationTable, apiv1.GenerationMeasurement{}, "MeasuredAtTime")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating bigquery.Client for GenerationMeasurement")
	}
	exchangeBigqueryClient, err := bigquery.NewClient(*bigqueryProjectID, *bigqueryEnable, *bigqueryDataset, *bigqueryGenerationTable, apiv1.ExchangeMeasurement{}, "MeasuredAtTime")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating bigquery.Client for ExchangeMeasurement")
	}

	// init bigquery table if it doesn't exist yet
	if *bigqueryInit {
		err = generationBigqueryClient.InitBigqueryTable()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed initializing bigquery table for GenerationMeasurement")
		}
		err = exchangeBigqueryClient.InitBigqueryTable()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed initializing bigquery table for ExchangeMeasurement")
		}
	}

	// create kubernetes api client
	kubeClientConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed retrieving kubeClientConfig")
	}
	// creates the clientset
	kubeClientset, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating kubeClientset")
	}

	configClient, err := config.NewClient(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating config.Client")
	}

	stateClient, err := state.NewClient(kubeClientset, *measurementFilePath, *measurementFileConfigMapName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating state.Client")
	}

	entsoeClient, err := entsoe.NewClient(*entsoeToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating entsoe.client")
	}

	exporterService, err := exporter.NewService(generationBigqueryClient, exchangeBigqueryClient, configClient, stateClient, entsoeClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating exporter.Service")
	}

	gracefulShutdown, waitGroup := foundation.InitGracefulShutdownHandling()

	err = exporterService.Run(ctx, gracefulShutdown, waitGroup)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed running export")
	}

	waitGroup.Wait()
}
