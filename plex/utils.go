package plex

import (
	"io/ioutil"
	"net/http"
	"time"
)

// AddTokenHeader is a simple function that adds the token to a header.
func AddTokenHeader(h map[string]string, t string) map[string]string {
	h["X-Plex-Token"] = t
	return h
}

// SendRequest sends a HTTP request according to provided method and url.
func SendRequest(m, u string, h map[string]string) (*http.Response, []byte, error) {
	req, err := CreateRequest(m, u, h)
	if err != nil {
		return nil, nil, err
	}

	httpClient := &http.Client{Timeout: time.Second * 10}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, body, nil
}

// CreateRequest creates a HTTP request according to provided method and url.
// Headers are added to the request.
func CreateRequest(m, u string, h map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range h {
		req.Header.Set(k, v)
	}

	return req, nil
}
