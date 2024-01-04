package upload

type UploadStorage interface {
	Connect() bool
	AppendToSetPipelined(string, [][]byte)
}
