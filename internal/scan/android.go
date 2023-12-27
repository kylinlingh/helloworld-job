package scan

import (
	"context"
	log "helloworld/pkg/logger"
)

type AndroidScanJob struct {
}

func (a *AndroidScanJob) RunJob(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("android scan job started")
	return nil
}

func (a *AndroidScanJob) ScanStaticCode(ctx context.Context) error {
	return nil
}

func (a *AndroidScanJob) ScanBinaryArtifacts(ctx context.Context) error {
	return nil
}
