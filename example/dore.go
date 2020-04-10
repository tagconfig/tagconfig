package main

import (
	"net/http"
	"time"
)

type TokenDoer struct {
	Token  string
	Client *http.Client
}

func (t *TokenDoer) Do(req *http.Request) (*http.Response, error) {
	if t.Client == nil {
		t.Client = &http.Client{
			Timeout: time.Second * 60,
		}
	}

	query := req.URL.Query()
	query.Add("token", t.Token)
	req.URL.RawQuery = query.Encode()
	return t.Client.Do(req)
}
