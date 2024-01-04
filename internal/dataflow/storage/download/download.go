package download

type DownloadStroage interface {
	GetName() string
	Connect() bool
	GetAndDeleteSet(string) []interface{}
}
