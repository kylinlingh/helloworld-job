package scan

import (
	"context"
	"errors"
	"fmt"
	"helloworld/internal/datastore"
	log "helloworld/pkg/logger"
)

var jobMap map[string]func(factory datastore.DBFactory) ApplicationScanJob

const (
	ANDROID_SOURCECODE = "android_sc"
	ANDROID_ARTIFACT   = "android_af"
	IOS_SCOURCECODE    = "ios_sc"
	IOS_ARTIFACT       = "ios_af"
)

type ApplicationScanJob interface {
	// RunJob returns a function that can be run via an errgroup to perform the scan job.
	RunJob(ctx context.Context) error
}

func initJobMap() {
	jobMap = make(map[string]func(factory datastore.DBFactory) ApplicationScanJob)
	jobMap[ANDROID_SOURCECODE] = NewAndroidSourceCodeScanJob
	jobMap[ANDROID_ARTIFACT] = NewAndroidArtifactScanJob
	jobMap[IOS_SCOURCECODE] = NewIOSSourceCodeScanJob
	jobMap[IOS_ARTIFACT] = NewIOSArtifactScanJob
}

func CreateScanJobs(ctx context.Context, dbInstance datastore.DBFactory, scanTargets []string) []ApplicationScanJob {
	initJobMap()
	// 创建多个扫描任务
	scanJobs := make([]ApplicationScanJob, 0, len(scanTargets))
	for _, st := range scanTargets {
		if initFunc, ok := jobMap[st]; !ok {
			log.Ctx(ctx).Fatal().Err(errors.New(fmt.Sprintf("invalid scan target: %s [%s, %s, %s, %s]", st, ANDROID_SOURCECODE, ANDROID_ARTIFACT, IOS_ARTIFACT, IOS_SCOURCECODE)))
		} else {
			scanJobs = append(scanJobs, initFunc(dbInstance))
		}
	}
	return scanJobs
}
