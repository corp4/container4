package api

import "net/rpc"

type Api struct {
	client *rpc.Client
}

func NewAPI(address string) (*Api, error) {
	client, err := rpc.DialHTTP("tcp", address)
	return &Api{client}, err
}

func (a *Api) SSH() *SSH {
	return &SSH{a.client}
}

func (a *Api) Close() error {
	return a.client.Close()
}
