package datastructure

import "sync"

type MessageList struct {
	Count   int
	ValList [][]byte
	Mutext  sync.Mutex
}
