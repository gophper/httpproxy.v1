package utils

import (
	"net"
	"time"
)

var cpmp = make(map[int64]ClientPair)

type ClientPair struct {
	InputClient  *ClientConn
	OutputClient *ClientConn
	ConnectId    int64
	Expire       time.Time
}

func NewClientPair(cs, ss net.Conn, connectId int64) ClientPair {
	return ClientPair{InputClient: NewClientConn(cs), OutputClient: NewClientConn(ss), ConnectId: connectId}
}

func init() {
	run()
}

func run() {
	go func() {
		var tk = time.NewTicker(time.Second)

		select {
		case <-tk.C:
			for key, value := range cpmp {
				if time.Now().After(value.Expire.Add(time.Minute)) {
					UnsetPair(key)
				}
			}
		}
	}()
}

func GetPair(connectId int64) *ClientPair {
	if r, ok := cpmp[connectId]; ok {
		return &r
	}
	return nil
}

func RebuildPair(connectId int64, cs net.Conn) ClientPair {
	a := cpmp[connectId]
	a.InputClient = NewClientConn(cs)
	return a
}

func UnsetPair(connectId int64) {
	delete(cpmp, connectId)
	return
}

func DelayUnsetPair(cp ClientPair) {
	cp.Expire = time.Now()
	cpmp[cp.ConnectId] = cp
	return
}

func RegisterPair(cp ClientPair) *ClientPair {
	cpmp[cp.ConnectId] = cp
	return &cp
}
