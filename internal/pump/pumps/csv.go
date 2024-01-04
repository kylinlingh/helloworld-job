package pumps

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"helloworld/internal/pump/analytics"
	log "helloworld/pkg/logger"
	"os"
	"path"
	"time"
)

// CSVPump defines a csv pumps with csv specific options and common options.
type CSVPump struct {
	csvConf *CSVConf
	CommonPumpConfig
}

// CSVConf defines csv specific options.
type CSVConf struct {
	// Specify the directory used to store automatically generated csv file which contains analyzed data.
	CSVDir string `mapstructure:"csv_dir"`
}

// New create a csv pumps instance.
func (c *CSVPump) New() PumpBackend {
	newPump := CSVPump{}

	return &newPump
}

// GetName returns the csv pumps name.
func (c *CSVPump) GetName() string {
	return "CSV Pump"
}

// Init initialize the csv pumps instance.
func (c *CSVPump) Init(conf interface{}) error {
	c.csvConf = &CSVConf{}
	err := mapstructure.Decode(conf, &c.csvConf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decode configuration")

	}

	ferr := os.MkdirAll(c.csvConf.CSVDir, 0o777)
	if ferr != nil {
		log.Error().Msg(ferr.Error())
	}

	log.Debug().Msg("csv initialized successfully")

	return nil
}

// WriteData write analyzed data to csv persistent back-end downloadfrom.
func (c *CSVPump) WriteData(ctx context.Context, data []interface{}) error {
	curtime := time.Now()
	fname := fmt.Sprintf("%d-%s-%d-%d.csv", curtime.Year(), curtime.Month().String(), curtime.Day(), curtime.Hour())
	fname = path.Join(c.csvConf.CSVDir, fname)

	var outfile *os.File
	var appendHeader bool

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		var createErr error
		outfile, createErr = os.Create(fname)
		if createErr != nil {
			log.Error().Err(createErr).Msg("failed to create new CSV file")
		}
		appendHeader = true
	} else {
		var appendErr error
		outfile, appendErr = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0o600)
		if appendErr != nil {
			log.Error().Err(appendErr).Msg("failed to open CSV file")
		}
	}

	defer outfile.Close()
	writer := csv.NewWriter(outfile)

	if appendHeader {
		startRecord := analytics.AnalyticsRecord{}
		headers := startRecord.GetFieldNames()

		err := writer.Write(headers)
		if err != nil {
			log.Error().Err(err).Msg("failed to write file headers")

			return errors.Wrap(err, "failed to write file headers")
		}
	}

	for _, v := range data {
		decoded, _ := v.(analytics.AnalyticsRecord)

		toWrite := decoded.GetLineValues()
		err := writer.Write(toWrite)
		if err != nil {
			log.Error().Msg("File write failed!")
			log.Error().Msg(err.Error())
		}
	}

	writer.Flush()

	return nil
}
