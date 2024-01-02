package pump

import (
	"helloworld/internal/pump/analytics"
	"helloworld/internal/pump/pumps"
	log "helloworld/pkg/logger"
	"time"
)

const anaylticsKeyName = "job-analytics"

// PumpConfig defines options for pump back-end.
type PumpConfig struct {
	Type                  string                     `json:"type"                    mapstructure:"type"`
	Filters               analytics.AnalyticsFilters `json:"filters"                 mapstructure:"filters"`
	Timeout               int                        `json:"timeout"                 mapstructure:"timeout"`
	OmitDetailedRecording bool                       `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	Meta                  map[string]interface{}     `json:"meta"                    mapstructure:"meta"`
}

type pumpService struct {
	secInterval  int
	omitDetails  bool
	pumpBackends map[string]PumpConfig
}

func (s *pumpService) initialize() {
	pmpBackes := make([]pumps.PumpBackend, len(s.pumpBackends))
	i := 0
	for key, pmp := range s.pumpBackends {
		pumpTypeName := pmp.Type
		if pumpTypeName == "" {
			pumpTypeName = key
		}

		pmpType, err := pumps.GetPumpBackendByName(pumpTypeName)
		if err != nil {
			log.Fatal().Err(err).Msg("pump pumps load error, you can register new pumps or delete it")
		} else {
			pmpIns := pmpType.New()
			initErr := pmpIns.Init(pmp.Meta)
			if initErr != nil {
				log.Fatal().Err(initErr).Msg("pump pumps initialized failed")
			} else {
				log.Info().Str("pump pumps", pmpIns.GetName()).Msg("initialized successfully")
				pmpIns.SetFilters(pmp.Filters)
				pmpIns.SetTimeout(pmp.Timeout)
				pmpIns.SetOmitDetailedRecording(pmp.OmitDetailedRecording)
				pmpBackes[i] = pmpIns
			}
		}
		i++
	}
}

func (p *pumpService) PrepareRun() preparedPumpService {
	p.initialize()
	return preparedPumpService{p}
}

func (p *pumpService) pump() {
	analyticsValues := p.pumpBackends.GetAndDeleteSet(anaylticsKeyName)
	if len(analyticsValues) == 0 {
		return
	}
}

type preparedPumpService struct {
	*pumpService
}

func (p *preparedPumpService) Run(stopCh <-chan struct{}) error {
	ticker := time.NewTicker(time.Duration(p.secInterval) * time.Second)
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

func createPumpService() (*pumpService, error) {

}
