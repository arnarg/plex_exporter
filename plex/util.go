package plex

import (
	"io/ioutil"
	"net/http"
)

// SendRequest sends a HTTP request according to provided method and url.
func sendRequest(m, u string, h map[string]string, c *http.Client) (*http.Response, []byte, error) {
	req, err := createRequest(m, u, h)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return resp, nil, err
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
func createRequest(m, u string, h map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range h {
		req.Header.Set(k, v)
	}

	return req, nil
}
