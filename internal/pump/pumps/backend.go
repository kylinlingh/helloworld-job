package pumps

import (
	"errors"
	"helloworld/internal/pump/analytics"
)

type PumpBackend interface {
	GetName() string
	New() PumpBackend
	Init(interface{}) error
	SetFilters(filters analytics.AnalyticsFilters)
	GetFilters() analytics.AnalyticsFilters
	SetTimeout(timeout int)
	GetTimeout() int
	SetOmitDetailedRecording(recording bool)
	GetOmitDetailedRecording() bool
}

func GetPumpBackendByName(name string) (PumpBackend, error) {
	if pump, ok := availablePumps[name]; ok && pump != nil {
		return pump, nil
	}
	return nil, errors.New("pump pumps: " + name + " not found")
}
