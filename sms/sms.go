// https://help.aliyun.com/document_detail/101300.html
package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Blank-Xu/aliyunsdk/utils"
)

var (
	ErrPhoneNumbersLimited = errors.New("phone numbers limited, SendSms can only send 1000 phone numbers one time")
	ErrPhoneNumbersIsNull  = errors.New("phone numbers required")
)

func NewApi(host, accessKeyId, accessSecret, regionId, signName, version string) *Option {
	op := new(Option)
	op.host = host
	op.accessKeyId = accessKeyId
	op.accessSecret = []byte(accessSecret + "&") // POP签名规则
	op.regionId = regionId
	op.version = version
	op.signName = signName
	return op
}

type Option struct {
	host         string
	accessKeyId  string
	accessSecret []byte
	regionId     string
	version      string
	signName     string
}

func (p *Option) SendSms(phones []string, templateCode, templateParam string) (*Response, error) {
	num := len(phones)
	if num > 1000 {
		return nil, ErrPhoneNumbersLimited
	}

	var phoneNumbersBuf bytes.Buffer
	phoneNumbersBuf.Grow(num * 15)
	for _, phone := range phones {
		l := len(phone)

		if l == 0 {
			continue
		} else if l > 3 && phone[:3] == "+86" {
			phoneNumbersBuf.WriteString(phone[3:])
		} else if phone[:1] == "+" {
			phoneNumbersBuf.WriteString(phone[1:])
		} else {
			phoneNumbersBuf.WriteString(phone)
		}

		phoneNumbersBuf.WriteByte(',')
	}

	l := phoneNumbersBuf.Len()
	if l == 0 {
		return nil, ErrPhoneNumbersIsNull
	}
	// delete last ','
	phoneNumbersBuf.Truncate(l - 1)

	// UTC time
	now := time.Now().UTC()

	params := map[string]string{
		// 公共参数
		"AccessKeyId":      p.accessKeyId,
		"Format":           "JSON",
		"RegionId":         p.regionId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   strconv.FormatInt(now.UnixNano(), 10), // 签名随机数
		"SignatureVersion": "1.0",
		"Timestamp":        now.Format(time.RFC3339)[:20], // ISO8601 标准
		"Version":          p.version,

		// 请求参数
		"Action":        "SendSms",
		"PhoneNumbers":  phoneNumbersBuf.String(),
		"SignName":      p.signName,
		"TemplateCode":  templateCode,
		"TemplateParam": templateParam,
	}

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return nil, err
	}

	_, body, err := utils.HttpRequest(p.host, http.MethodPost, param)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

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
		if l > 3 && phone[:3] == "+86" {
			phoneNumbers = append(phoneNumbers, phone[3:])
		} else if l > 0 && phone[:1] == "+" {
			phoneNumbers = append(phoneNumbers, phone[1:])
		} else {
			phoneNumbers = append(phoneNumbers, phone)
		}
	}

	phoneNumberJson, err := json.Marshal(phoneNumbers)
	if err != nil {
		return nil, fmt.Errorf("json marshal phones failed, err: %v", err)
	}

	signNameJson, err := json.Marshal(signNames)
	if err != nil {
		return nil, fmt.Errorf("json marshal signNames failed, err: %v", err)
	}

	templateParamJson, err := json.Marshal(templateParams)
	if err != nil {
		return nil, fmt.Errorf("json marshal templateParams failed, err: %v", err)
	}

	// UTC time
	now := time.Now().UTC()

	params := map[string]string{
		// 公共参数
		"AccessKeyId":      p.accessKeyId,
		"Format":           "JSON",
		"RegionId":         p.regionId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   strconv.FormatInt(now.UnixNano(), 10), // 签名随机数
		"SignatureVersion": "1.0",
		"Timestamp":        now.Format(time.RFC3339)[:20], // ISO8601 标准
		"Version":          p.version,

		// 请求参数
		"Action":            "SendBatchSms",
		"PhoneNumberJson":   string(phoneNumberJson),
		"SignNameJson":      string(signNameJson),
		"TemplateCode":      templateCode,
		"TemplateParamJson": string(templateParamJson),
	}

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return nil, err
	}

	_, body, err := utils.HttpRequest(p.host, http.MethodPost, param)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
