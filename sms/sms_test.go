package sms

import (
	"testing"
)

const (
	host         = ""
	accessKeyID  = ""
	accessSecret = ""
	regionID     = ""
	version      = ""

	signName      = ""
	templateCode  = ""
	templateParam = ``
)

var (
	api = NewAPI(host, accessKeyID, accessSecret, regionID, version)

	phones = []string{"+8615900000000"}

	signNames      = []string{""}
	templateParams = []string{``}
)

func TestSendSms(t *testing.T) {
	resp, err := api.SendSms(phones, signName, templateCode, templateParam)
	if err != nil {
		t.Fatalf("send failed, resp: %+v, err: %v", resp, err)
	}

	if err = resp.Error(); err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}

	t.Logf("send success, resp: %+v", resp)
}

func TestSendBatchSms(t *testing.T) {
	resp, err := api.SendBatchSms(phones, signNames, templateCode, templateParams)
	if err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}

	if err = resp.Error(); err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}

	t.Logf("send success, resp: %+v", resp)
}
