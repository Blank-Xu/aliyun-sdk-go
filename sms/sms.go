// Package sms .
//
// api for https://help.aliyun.com/document_detail/101300.html
//
package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Blank-Xu/aliyun-sdk-go/utils"
)

var (
	// ErrPhoneNumbersLimited  .
	ErrPhoneNumbersLimited = errors.New("phone numbers limited, SendSms can only send 1000 phone numbers one time")
	// ErrPhoneNumbersIsNull  .
	ErrPhoneNumbersIsNull = errors.New("phone numbers required")
)

// NewAPI  .
func NewAPI(host, accessKeyID, accessSecret, regionID, version string) *Option {
	op := new(Option)

	op.host = host
	op.accessSecret = []byte(accessSecret + "&") // POP签名规则

	// prepare params
	op.pubParam = [][2]string{
		{"Format", "JSON"},
		{"Version", version},
		{"AccessKeyId", accessKeyID},
		{"SignatureMethod", "HMAC-SHA1"},
		{"SignatureVersion", "1.0"},
		{"RegionId", regionID},
	}

	return op
}

// Option  .
type Option struct {
	host         string
	accessSecret []byte

	pubParam [][2]string
}

// SendSms  .
func (p *Option) SendSms(phones []string, signName, templateCode, templateParam string) (*Response, error) {
	num := len(phones)
	if num > 1000 {
		return nil, ErrPhoneNumbersLimited
	}

	var phonesBuf bytes.Buffer

	phonesBuf.Grow(num * 15)
	for _, phone := range phones {
		l := len(phone)

		if l == 0 {
			continue
		} else if l == 14 && phone[:3] == "+86" {
			// for chinese phone number
			phonesBuf.WriteString(phone[3:])
		} else if phone[:1] == "+" {
			phonesBuf.WriteString(phone[1:])
		} else {
			phonesBuf.WriteString(phone)
		}

		phonesBuf.WriteByte(',')
	}

	l := phonesBuf.Len()
	if l == 0 {
		return nil, ErrPhoneNumbersIsNull
	}
	// delete last ','
	phonesBuf.Truncate(l - 1)

	// UTC time
	now := time.Now().UTC()

	params := append(p.pubParam,
		[][2]string{
			{"Timestamp", now.Format(time.RFC3339)[:20]},              // ISO8601 标准
			{"SignatureNonce", strconv.FormatInt(now.UnixNano(), 10)}, // 签名随机数
			{"Action", "SendSms"},
			{"PhoneNumbers", phonesBuf.String()},
			{"SignName", signName},
			{"TemplateCode", templateCode},
			{"TemplateParam", templateParam},
		}...)

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return nil, err
	}

	code, body, err := utils.HTTPPost(p.host, param)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("response failed, code: %d, body: %s, err: %v", code, body, err)
	}

	var resp Response

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SendBatchSms  .
func (p *Option) SendBatchSms(phones, signNames []string, templateCode string, templateParams []string) (*Response, error) {
	nums := len(phones)
	if nums == 0 {
		return nil, ErrPhoneNumbersIsNull
	}
	if nums != len(signNames) {
		return nil, errors.New("every phone number should have a sign name")
	}
	if nums != len(templateParams) {
		return nil, errors.New("every phone number should have a templateParamJson")
	}

	phoneNumbers := make([]string, 0, nums)

	for _, phone := range phones {
		l := len(phone)

		if l == 14 && phone[:3] == "+86" {
			// for chinese phone number
			phoneNumbers = append(phoneNumbers, phone[3:])
		} else if l > 0 && phone[:1] == "+" {
			phoneNumbers = append(phoneNumbers, phone[1:])
		} else {
			phoneNumbers = append(phoneNumbers, phone)
		}
	}

	phoneNumberJSON, err := json.Marshal(phoneNumbers)
	if err != nil {
		return nil, fmt.Errorf("json marshal phones failed, err: %v", err)
	}

	signNameJSON, err := json.Marshal(signNames)
	if err != nil {
		return nil, fmt.Errorf("json marshal signNames failed, err: %v", err)
	}

	templateParamJSON, err := json.Marshal(templateParams)
	if err != nil {
		return nil, fmt.Errorf("json marshal templateParams failed, err: %v", err)
	}

	// UTC time
	now := time.Now().UTC()

	params := append(p.pubParam,
		[][2]string{
			{"Timestamp", now.Format(time.RFC3339)[:20]},              // ISO8601 标准
			{"SignatureNonce", strconv.FormatInt(now.UnixNano(), 10)}, // 签名随机数
			{"Action", "SendBatchSms"},
			{"PhoneNumberJson", string(phoneNumberJSON)},
			{"SignNameJson", string(signNameJSON)},
			{"TemplateCode", templateCode},
			{"TemplateParamJson", string(templateParamJSON)},
		}...)

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return nil, err
	}

	code, body, err := utils.HTTPPost(p.host, param)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("response failed, code: %d, body: %s, err: %v", code, body, err)
	}

	var resp Response

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
