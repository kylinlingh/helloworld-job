package downloadfrom

// DownloadHandler defines the analytics downloadfrom interface.
type DownloadHandler interface {
	Init(config interface{}) error
	GetName() string
	Connect() bool
	GetAndDeleteSet(string) []interface{}
}

const (
	// AnalyticsKeyName defines the key name in redis which used to analytics.
	AnalyticsKeyName string = "job-analytics"
)
