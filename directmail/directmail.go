// Package directmail  .
//
// api for https://help.aliyun.com/document_detail/29434.html
//
package directmail

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Blank-Xu/aliyun-sdk-go/utils"
)

var (
	// ErrEmailNumbersLimited  .
	ErrEmailNumbersLimited = errors.New("email numbers limited, SingleSendMail can only send 100 email numbers one time")
)

// NewAPI  .
func NewAPI(host, accessKeyID, accessSecret, regionID, version, accountName, fromAlias string) *Option {
	op := new(Option)

	op.host = host
	op.accessSecret = []byte(accessSecret + "&") // 签名要求

	// prepare params
	pubParam := [][2]string{
		// 公共参数
		{"Format", "JSON"},
		{"Version", version},
		{"AccessKeyId", accessKeyID},
		{"SignatureMethod", "HMAC-SHA1"},
		{"SignatureVersion", "1.0"},
		{"RegionId", regionID},
		{"AccountName", accountName},
		{"AddressType", "1"},
	}

	op.paramSingleMail = append(pubParam, [][2]string{
		{"ReplyToAddress", "false"},
		{"FromAlias", fromAlias},
		{"Action", "SingleSendMail"},
	}...)

	op.paramBatchMail = append(pubParam, [][2]string{
		{"Action", "BatchSendMail"},
	}...)

	return op
}

// Option  .
type Option struct {
	host         string
	accessSecret []byte

	paramSingleMail [][2]string
	paramBatchMail  [][2]string
}

// SingleSendMail  .
func (p *Option) SingleSendMail(addresses []string, subject, htmlBody, textBody string) (int, []byte, error) {
	if len(addresses) > 100 {
		return 0, nil, ErrEmailNumbersLimited
	}

	// UTC time
	now := time.Now().UTC()

	params := make([][2]string, 0, 18)
	params = append(params, p.paramSingleMail...)
	params = append(params, [][2]string{
		{"Timestamp", now.Format(time.RFC3339)[:20]},              // ISO8601 标准
		{"SignatureNonce", strconv.FormatInt(now.UnixNano(), 10)}, // 签名随机数
		{"ToAddress", strings.Join(addresses, ",")},
		{"Subject", subject},
		{"HtmlBody", htmlBody},
		{"TextBody", textBody},
	}...)

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return 0, nil, err
	}

	return utils.HTTPPost(p.host, param)
}

// BatchSendMail  .
func (p *Option) BatchSendMail(receiversName, templateName, tagName string) (int, []byte, error) {
	// UTC time
	now := time.Now().UTC()

	params := make([][2]string, 0, 18)
	params = append(params, p.paramSingleMail...)
	params = append(params, [][2]string{
		{"Timestamp", now.Format(time.RFC3339)[:20]},              // ISO8601 标准
		{"SignatureNonce", strconv.FormatInt(now.UnixNano(), 10)}, // 签名随机数
		{"ReceiversName", receiversName},
		{"TemplateName", templateName},
		{"TagName", tagName},
	}...)

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return 0, nil, err
	}

	return utils.HTTPPost(p.host, param)
}
