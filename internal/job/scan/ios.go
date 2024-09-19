package scan

import (
	"context"
	"helloworld/internal/datastore"
	log "helloworld/pkg/logger"
)

type IOSScanJob struct {
	store datastore.DBFactory
}

type IOSSourceCodeScanJob struct {
	IOSScanJob
}

func NewIOSSourceCodeScanJob(ds datastore.DBFactory) ApplicationScanJob {
	return &IOSSourceCodeScanJob{IOSScanJob{store: ds}}
}

func (i *IOSSourceCodeScanJob) RunJob(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ios source code scan job started")

	return nil
}

type IOSArtifactScanJob struct {
	IOSScanJob
}

func NewIOSArtifactScanJob(ds datastore.DBFactory) ApplicationScanJob {
	return &IOSArtifactScanJob{IOSScanJob{store: ds}}
}

func (i *IOSArtifactScanJob) RunJob(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ios artifact scan job started")

	return nil
}

func (i *IOSScanJob) RunJob(ctx context.Context) error {
	return nil
}
