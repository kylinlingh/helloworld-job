package scan

import (
	"context"
)

// ScanJob
type ScanJob interface {
	// RunJob returns a function that can be run via an errgroup to perform the scan service.
	RunJob(ctx context.Context) error
}
