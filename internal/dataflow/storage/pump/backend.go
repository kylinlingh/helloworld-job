package pump

import (
	"context"
	"errors"
	"helloworld/internal/dataflow/datastructure"
	"time"
)

type PumpBackend interface {
	GetName() string
	New() PumpBackend
	Init(interface{}) error
	WriteData(ctx context.Context, keys []interface{}) error
	SetFilters(filters datastructure.AnalyticsFilters)
	GetFilters() datastructure.AnalyticsFilters
	SetTimeout(timeout time.Duration)
	GetTimeout() time.Duration
	SetOmitDetailedRecording(recording bool)
	GetOmitDetailedRecording() bool
}

func GetPumpBackendByName(name string) (PumpBackend, error) {
	if pump, ok := availablePumps[name]; ok && pump != nil {
		return pump, nil
	}
	return nil, errors.New("pump pumps: " + name + " not found")
}
