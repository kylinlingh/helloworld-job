package memory

type memory struct {
}

func (m *memory) Connect() bool {
	return true
}

func (m *memory) AppendToSetPipelined(string, [][]byte) {
	// 写入到 channel 中
}
