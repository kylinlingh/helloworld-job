package pump

import (
	"helloworld/internal/dataflow/datastructure"

	"time"
)

// CommonPumpConfig defines common options used by all persistent store, like elasticsearch, kafka, mongo and etc.
type CommonPumpConfig struct {
	filters               datastructure.AnalyticsFilters
	timeout               time.Duration
	OmitDetailedRecording bool
}

// SetFilters set attributes `filters` for CommonPumpConfig.
func (p *CommonPumpConfig) SetFilters(filters datastructure.AnalyticsFilters) {
	p.filters = filters
}

// GetFilters get attributes `filters` for CommonPumpConfig.
func (p *CommonPumpConfig) GetFilters() datastructure.AnalyticsFilters {
	return p.filters
}

// SetTimeout set attributes `timeout` for CommonPumpConfig.
func (p *CommonPumpConfig) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

// GetTimeout get attributes `timeout` for CommonPumpConfig.
func (p *CommonPumpConfig) GetTimeout() time.Duration {
	return p.timeout
}

// SetOmitDetailedRecording set attributes `OmitDetailedRecording` for CommonPumpConfig.
func (p *CommonPumpConfig) SetOmitDetailedRecording(omitDetailedRecording bool) {
	p.OmitDetailedRecording = omitDetailedRecording
}

// GetOmitDetailedRecording get attributes `OmitDetailedRecording` for CommonPumpConfig.
func (p *CommonPumpConfig) GetOmitDetailedRecording() bool {
	return p.OmitDetailedRecording
}
