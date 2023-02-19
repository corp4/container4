package api

import (
	"net/rpc"
	"time"
)

type Supervisor struct {
	client *rpc.Client
}

func (s *Supervisor) GetStatus() (status string, err error) {
	err = s.client.Call("Supervisor.GetStatus", struct{}{}, &status)
	return
}

func (s *Supervisor) GetTime() (t time.Time, err error) {
	err = s.client.Call("Supervisor.GetTime", struct{}{}, &t)
	return
}
