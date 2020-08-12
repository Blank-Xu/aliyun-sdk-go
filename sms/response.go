package sms

import (
	"fmt"
)

// Response  for sms request.
type Response struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	BizID     string `json:"BizId"`
}

// Error  .
func (p Response) Error() error {
	if p.Code != "OK" {
		return fmt.Errorf("send sms failed, code: %s, message: %s", p.Code, p.Message)
	}

	return nil
}
