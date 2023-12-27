package scan

import (
	"context"
	log "helloworld/pkg/logger"
)

type IOSScanJob struct {
}

func (a *IOSScanJob) RunJob(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ios scan job started")
	return nil
}

func (a *IOSScanJob) ScanStaticCode(ctx context.Context) error {
	return nil
}

func (a *IOSScanJob) ScanBinaryArtifacts(ctx context.Context) error {
	return nil
}
