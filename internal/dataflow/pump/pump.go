package pump

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/vmihailenco/msgpack/v5"
	"helloworld/internal/dataflow/datastructure"
	"helloworld/internal/dataflow/storage/downloadfrom"
	"helloworld/internal/dataflow/storage/downloadfrom/memory"
	"helloworld/internal/dataflow/storage/pump"
	log "helloworld/pkg/logger"
	"sync"
	"time"
)

const anaylticsKeyName = "job-analytics"

// PumpConfig defines options for pump back-end.
type PumpConfig struct {
	Type                  string                         `json:"type"                    mapstructure:"type"`
	Filters               datastructure.AnalyticsFilters `json:"filters"                 mapstructure:"filters"`
	Timeout               time.Duration                  `json:"timeout"                 mapstructure:"timeout"`
	OmitDetailedRecording bool                           `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	Meta                  map[string]interface{}         `json:"meta"                    mapstructure:"meta"`
}

type pumpService struct {
	// 多久去拉取一次信息
	secInterval  time.Duration
	omitDetails  bool
	pumpBackends map[string]PumpConfig
	dlstorage    downloadfrom.DownloadStroage
}

func CreatePumpService(purgeDelay time.Duration, pc map[string]PumpConfig, cache *ristretto.Cache) *pumpService {
	service := &pumpService{
		secInterval:  purgeDelay,
		omitDetails:  false,
		pumpBackends: pc,
		dlstorage:    memory.NewDownloadMemStorage(cache),
	}
	return service
}

var pmps []pump.PumpBackend

func (s *pumpService) initialize() {
	pmps = make([]pump.PumpBackend, len(s.pumpBackends))
	i := 0
	for key, pmp := range s.pumpBackends {
		pumpTypeName := pmp.Type
		if pumpTypeName == "" {
			pumpTypeName = key
		}

		pmpType, err := pump.GetPumpBackendByName(pumpTypeName)
		if err != nil {
			log.Fatal().Err(err).Msg("pump pumps load error, you can register new pumps or delete it")
		} else {
			pmpIns := pmpType.New()
			initErr := pmpIns.Init(pmp.Meta)
			if initErr != nil {
				log.Fatal().Err(initErr).Msg("pump pumps initialized failed")
			} else {
				log.Info().Msgf("csv pump-backend initialized successfully")
				pmpIns.SetFilters(pmp.Filters)
				pmpIns.SetTimeout(pmp.Timeout)
				pmpIns.SetOmitDetailedRecording(pmp.OmitDetailedRecording)
				pmps[i] = pmpIns
			}
		}
		i++
	}
}

func (p *pumpService) PrepareRun() *preparedPumpService {
	p.initialize()
	return &preparedPumpService{p}
}

func (p *pumpService) pump() {
	analyticsValues := p.dlstorage.GetAndDeleteSet(anaylticsKeyName)
	if len(analyticsValues) == 0 {
		return
	}
	// Convert to something clean
	keys := make([]interface{}, len(analyticsValues))

	for i, v := range analyticsValues {
		decoded := datastructure.AnalyticsRecord{}
		err := msgpack.Unmarshal(v.([]byte), &decoded)
		log.Info().Msgf("data decoded successfully: %v", decoded)
		if err != nil {
			log.Error().Err(err).Msg("couldn't unmarshal analytics data")
		} else {
			if p.omitDetails {
				//decoded.Policies = ""
				//decoded.Deciders = ""
			}
			keys[i] = interface{}(decoded)
		}
	}

	// Send to pumps
	writeToPumps(keys, p.secInterval)
}

func writeToPumps(keys []interface{}, purgeDelay time.Duration) {
	// Send to pumps
	if pmps != nil {
		var wg sync.WaitGroup
		wg.Add(len(pmps))
		for _, pmp := range pmps {
			go execPumpWriting(&wg, pmp, &keys, purgeDelay)
		}
		wg.Wait()
	} else {
		log.Warn().Msg("no pumps defined!")
	}
}

func execPumpWriting(wg *sync.WaitGroup, pmp pump.PumpBackend, keys *[]interface{}, purgeDelay time.Duration) {
	timer := time.AfterFunc(purgeDelay, func() {
		if pmp.GetTimeout() == 0 {
			log.Warn().Msgf(
				"pump %s is taking more time than the value configured of purge_delay. You should try to set a timeout for this pump.",
				pmp.GetName(),
			)
		} else if pmp.GetTimeout() > purgeDelay {
			log.Warn().Msgf("pump %s is taking more time than the value configured of purge_delay. You should try lowering the timeout configured for this pump.", pmp.GetName())
		}
	})
	defer timer.Stop()
	defer wg.Done()

	log.Debug().Msgf("writing to: %s", pmp.GetName())

	ch := make(chan error, 1)
	var ctx context.Context
	var cancel context.CancelFunc
	// Initialize context depending if the pump has a configured timeout
	if tm := pmp.GetTimeout(); tm > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(tm)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	defer cancel()

	go func(ch chan error, ctx context.Context, pmp pump.PumpBackend, keys *[]interface{}) {
		filteredKeys := filterData(pmp, *keys)

		ch <- pmp.WriteData(ctx, filteredKeys)
	}(ch, ctx, pmp, keys)

	select {
	case err := <-ch:
		if err != nil {
			log.Warn().Msgf("error writing to: %s - Error: %s", pmp.GetName(), err.Error())
		}
	case <-ctx.Done():
		//nolint: errorlint
		switch ctx.Err() {
		case context.Canceled:
			log.Warn().Msgf("The writing to %s have got canceled.", pmp.GetName())
		case context.DeadlineExceeded:
			log.Warn().Msgf("Timeout Writing to: %s", pmp.GetName())
		}
	}
}

func filterData(pump pump.PumpBackend, keys []interface{}) []interface{} {
	filters := pump.GetFilters()
	if !filters.HasFilter() && !pump.GetOmitDetailedRecording() {
		return keys
	}
	filteredKeys := keys[:] //nolint: gocritic
	newLenght := 0

	for _, key := range filteredKeys {
		decoded, _ := key.(datastructure.AnalyticsRecord)
		if pump.GetOmitDetailedRecording() {
			//decoded.Policies = ""
			//decoded.Deciders = ""
		}
		if filters.ShouldFilter(decoded) {
			continue
		}
		filteredKeys[newLenght] = decoded
		newLenght++
	}
	filteredKeys = filteredKeys[:newLenght]

	return filteredKeys
}

type preparedPumpService struct {
	*pumpService
}

func (p *preparedPumpService) Run(stopCh <-chan struct{}) error {
	ticker := time.NewTicker(p.secInterval)
	defer ticker.Stop()

	log.Info().Msg("run loop to pump data")
	for {
		select {
		case <-ticker.C:
			p.pump()
		case <-stopCh:
			log.Info().Msg("stop purge loop")
			return nil
		}
	}
}
