package utils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	defaultClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 30,
	}
)

// HTTPGet  .
func HTTPGet(host string, param string) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, host+"?"+param, nil)
	if err != nil {
		return 0, nil, err
	}

	return HTTPRequest(req)
}

// HTTPPost  .
func HTTPPost(host string, param string) (int, []byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		host,
		strings.NewReader(param),
	)
	if err != nil {
		return 0, nil, err
	}

	return HTTPRequest(req)
}

// HTTPRequest  .
func HTTPRequest(req *http.Request) (int, []byte, error) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, body, err
}
