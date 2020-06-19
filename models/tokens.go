package models

import "encoding/json"

type TokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (t TokenPair) Marhsal() ([]byte, error) {
	return json.Marshal(t)
}
