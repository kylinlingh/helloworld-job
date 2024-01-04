package datastructure

import "sync"

type MessageList struct {
	ValList [][]byte
	Mutext  sync.Mutex
}
