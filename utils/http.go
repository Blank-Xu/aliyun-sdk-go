package utils

import (
	"crypto/tls"
	"fmt"
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

func HttpRequest(host, method, param string) (int, []byte, error) {
	var (
		req *http.Request
		err error
	)

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(
			http.MethodGet,
			host+"?"+param,
			nil,
		)

	case http.MethodPost:
		req, err = http.NewRequest(
			http.MethodPost,
			host,
			strings.NewReader(param),
		)

	default:
		return 0, nil, fmt.Errorf("not support method: %s", method)
	}

	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, body, err
}
