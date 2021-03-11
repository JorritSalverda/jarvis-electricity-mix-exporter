package main

import (
	"context"
	"runtime"

	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/bigquery"
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
	entsoeArea  = kingpin.Flag("entsoe-area", "Area to retrieve electricity mix for ").Envar("ENTSOE_AREA").Required().String()

	bigqueryEnable    = kingpin.Flag("bigquery-enable", "Toggle to enable or disable bigquery integration").Default("true").OverrideDefaultFromEnvar("BQ_ENABLE").Bool()
	bigqueryInit      = kingpin.Flag("bigquery-init", "Toggle to enable bigquery table initialization").Default("true").OverrideDefaultFromEnvar("BQ_INIT").Bool()
	bigqueryProjectID = kingpin.Flag("bigquery-project-id", "Google Cloud project id that contains the BigQuery dataset").Envar("BQ_PROJECT_ID").Required().String()
	bigqueryDataset   = kingpin.Flag("bigquery-dataset", "Name of the BigQuery dataset").Envar("BQ_DATASET").Required().String()
	bigqueryTable     = kingpin.Flag("bigquery-table", "Name of the BigQuery table").Envar("BQ_TABLE").Required().String()

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
	bigqueryClient, err := bigquery.NewClient(*bigqueryProjectID, *bigqueryEnable, *bigqueryDataset, *bigqueryTable)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating bigquery.Client")
	}

	// init bigquery table if it doesn't exist yet
	if *bigqueryInit {
		err = bigqueryClient.InitBigqueryTable()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed initializing bigquery table")
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

	stateClient, err := state.NewClient(ctx, kubeClientset, *measurementFilePath, *measurementFileConfigMapName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating state.Client")
	}

	entsoeClient, err := entsoe.NewClient(*entsoeToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating entsoe.client")
	}

	exporterService, err := exporter.NewService(bigqueryClient, stateClient, entsoeClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating exporter.Service")
	}

	gracefulShutdown, waitGroup := foundation.InitGracefulShutdownHandling()

	err = exporterService.Run(ctx, waitGroup, entsoe.Area(*entsoeArea))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed running export")
	}

	foundation.HandleGracefulShutdown(gracefulShutdown, waitGroup)
}
