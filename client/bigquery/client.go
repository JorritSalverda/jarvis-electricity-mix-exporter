package bigquery

import (
	"context"
	"fmt"
	"time"

	googlebigquery "cloud.google.com/go/bigquery"
	contractsv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/contracts/v1"
	"github.com/rs/zerolog/log"
)

// Client is the interface for connecting to bigquery
type Client interface {
	CheckIfDatasetExists() (exists bool)
	CheckIfTableExists() (exists bool)
	CreateTable(typeForSchema interface{}, partitionField string, waitReady bool) (err error)
	UpdateTableSchema(typeForSchema interface{}) (err error)
	DeleteTable() (err error)
	InsertMeasurement(measurement contractsv1.Measurement) (err error)
	InitBigqueryTable() (err error)
}

// NewClient returns new bigquery.Client
func NewClient(projectID string, enable bool, dataset, table string) (Client, error) {

	ctx := context.Background()

	bigqueryClient, err := googlebigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &client{
		projectID: projectID,
		client:    bigqueryClient,
		enable:    enable,
		dataset:   dataset,
		table:     table,
	}, nil
}

type client struct {
	projectID string
	client    *googlebigquery.Client
	enable    bool
	dataset   string
	table     string
}

func (c *client) CheckIfDatasetExists() (exists bool) {

	if !c.enable {
		return false
	}

	ds := c.client.Dataset(c.dataset)

	md, err := ds.Metadata(context.Background())

	log.Error().Err(err).Msgf("Error retrieving metadata for dataset %v", c.dataset)

	return md != nil
}

func (c *client) CheckIfTableExists() (exists bool) {

	if !c.enable {
		return false
	}

	tbl := c.client.Dataset(c.dataset).Table(c.table)

	md, _ := tbl.Metadata(context.Background())

	// log.Error().Err(err).Msgf("Error retrieving metadata for table %v", table)

	return md != nil
}

func (c *client) CreateTable(typeForSchema interface{}, partitionField string, waitReady bool) (err error) {

	if !c.enable {
		return nil
	}

	tbl := c.client.Dataset(c.dataset).Table(c.table)

	// infer the schema of the type
	schema, err := googlebigquery.InferSchema(typeForSchema)
	if err != nil {
		return err
	}

	tableMetadata := &googlebigquery.TableMetadata{
		Schema: schema,
	}

	// if partitionField is set use it for time partitioning
	if partitionField != "" {
		tableMetadata.TimePartitioning = &googlebigquery.TimePartitioning{
			Field: partitionField,
		}
	}

	// create the table
	err = tbl.Create(context.Background(), tableMetadata)
	if err != nil {
		return err
	}

	if waitReady {
		for {
			if c.CheckIfTableExists() {
				break
			}
			time.Sleep(time.Second)
		}
	}

	return nil
}

func (c *client) UpdateTableSchema(typeForSchema interface{}) (err error) {

	if !c.enable {
		return nil
	}

	tbl := c.client.Dataset(c.dataset).Table(c.table)

	// infer the schema of the type
	schema, err := googlebigquery.InferSchema(typeForSchema)
	if err != nil {
		return err
	}

	meta, err := tbl.Metadata(context.Background())
	if err != nil {
		return err
	}

	update := googlebigquery.TableMetadataToUpdate{
		Schema: schema,
	}
	if _, err := tbl.Update(context.Background(), update, meta.ETag); err != nil {
		return err
	}

	return nil
}

func (c *client) DeleteTable() (err error) {

	if !c.enable {
		return nil
	}

	tbl := c.client.Dataset(c.dataset).Table(c.table)

	// delete the table
	err = tbl.Delete(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (c *client) InsertMeasurement(measurement contractsv1.Measurement) (err error) {

	if !c.enable {
		return nil
	}

	tbl := c.client.Dataset(c.dataset).Table(c.table)

	u := tbl.Uploader()

	if err := u.Put(context.Background(), measurement); err != nil {
		return err
	}

	return nil
}

func (c *client) InitBigqueryTable() (err error) {

	log.Debug().Msgf("Checking if table %v.%v.%v exists...", c.projectID, c.dataset, c.table)
	tableExist := c.CheckIfTableExists()

	if !tableExist {
		log.Debug().Msgf("Creating table %v.%v.%v...", c.projectID, c.dataset, c.table)
		err := c.CreateTable(contractsv1.Measurement{}, "MeasuredAtTime", true)
		if err != nil {
			return fmt.Errorf("Failed creating bigquery table: %w", err)
		}
	} else {
		log.Debug().Msgf("Trying to update table %v.%v.%v schema...", c.projectID, c.dataset, c.table)
		err := c.UpdateTableSchema(contractsv1.Measurement{})
		if err != nil {
			return fmt.Errorf("Failed updating bigquery table schema: %w", err)
		}
	}

	return nil
}
