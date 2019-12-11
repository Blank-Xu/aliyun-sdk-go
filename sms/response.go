package sms

import (
	"fmt"
)

// Response of sms request.
type Response struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	BizId     string `json:"BizId"`
}

func (p *Response) CheckError() error {
	if p.Code != "OK" {
		return fmt.Errorf("send sms failed, code: %s, message: %s", p.Code, p.Message)
	}

	return nil
}
