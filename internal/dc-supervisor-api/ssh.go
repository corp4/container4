package api

import "net/rpc"

type SSH struct {
	client *rpc.Client
}

func (s *SSH) AddAuthorizedKey(pubKey string) (success bool, err error) {
	err = s.client.Call("SSH.AddAuthorizedKey", pubKey, &success)
	return
}

func (s *SSH) HasAuthorizedKey(pubkey string) (hasKey bool, err error) {
	err = s.client.Call("SSH.HasAuthorizedKey", pubkey, &hasKey)
	return
}
