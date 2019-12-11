package sms

import (
	"testing"
)

const (
	tHost         = ""
	tAccessKeyId  = ""
	tAccessSecret = ""
	tRegionId     = ""
	tSignName     = ""
	tVersion      = ""
)

var (
	tApi = NewApi(tHost, tAccessKeyId, tAccessSecret, tRegionId, tSignName, tVersion)

	tPhones       = []string{"+8615900000000"}
	tTemplateCode = ""

	tTemplateParam = ``

	tSignNames      = []string{"ali sms test"}
	tTemplateParams = []string{``}
)

func TestSendSms(t *testing.T) {
	resp, err := tApi.SendSms(tPhones, tTemplateCode, tTemplateParam)
	if err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}
	if err = resp.CheckError(); err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}

	t.Logf("send success, resp: %+v", resp)
}

func TestSendBatchSms(t *testing.T) {
	resp, err := tApi.SendBatchSms(tPhones, tSignNames, tTemplateCode, tTemplateParams)
	if err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}
	if err = resp.CheckError(); err != nil {
		t.Fatalf("send failed, resp: %s, err: %v", resp, err)
	}

	t.Logf("send success, resp: %+v", resp)
}
