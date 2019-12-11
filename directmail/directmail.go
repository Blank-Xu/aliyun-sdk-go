// https://help.aliyun.com/document_detail/29434.html
package directmail

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Blank-Xu/aliyunsdk/utils"
)

var (
	ErrEmailNumbersLimited = errors.New("email numbers limited, SingleSendMail can only send 100 email numbers one time")
)

func NewApi(host, accessKeyId, accessSecret, regionId, version, accountName, fromAlias string) *Option {
	op := new(Option)
	op.host = host
	op.accessKeyId = accessKeyId
	op.accessSecret = []byte(accessSecret + "&") // 签名要求
	op.regionId = regionId
	op.version = version
	op.accountName = accountName
	op.fromAlias = fromAlias
	return op
}

type Option struct {
	host         string
	accessKeyId  string
	accessSecret []byte
	regionId     string
	version      string
	accountName  string
	fromAlias    string
}

func (p *Option) SingleSendMail(emails []string, subject, htmlBody, textBody string) (statusCode int, response []byte, err error) {
	if len(emails) > 100 {
		err = ErrEmailNumbersLimited
		return
	}

	// UTC time
	now := time.Now().UTC()

	params := map[string]string{
		// 公共参数
		"Format":           "JSON",
		"Version":          p.version,
		"AccessKeyId":      p.accessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        now.Format(time.RFC3339)[:20], // ISO8601 标准
		"SignatureVersion": "1.0",
		"SignatureNonce":   strconv.FormatInt(now.UnixNano(), 10), // 签名随机数
		"RegionId":         p.regionId,

		// 请求参数
		"Action":         "SingleSendMail",
		"AccountName":    p.accountName,
		"AddressType":    "1",
		"ReplyToAddress": "false",
		"FromAlias":      p.fromAlias,
		"ToAddress":      strings.Join(emails, ","),
		"Subject":        subject,
		"HtmlBody":       htmlBody,
		"TextBody":       textBody,
	}

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return
	}

	return utils.HttpRequest(p.host, http.MethodPost, param)
}

func (p *Option) BatchSendMail(receiversName, templateName, tagName string) (statusCode int, response []byte, err error) {
	// UTC time
	now := time.Now().UTC()

	params := map[string]string{
		// 公共参数
		"Format":           "JSON",
		"Version":          p.version,
		"AccessKeyId":      p.accessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        now.Format(time.RFC3339)[:20], // ISO8601 标准
		"SignatureVersion": "1.0",
		"SignatureNonce":   strconv.FormatInt(now.UnixNano(), 10), // 签名随机数
		"RegionId":         p.regionId,

		// 请求参数
		"Action":        "BatchSendMail",
		"AccountName":   p.accountName,
		"AddressType":   "1",
		"ReceiversName": receiversName,
		"TemplateName":  templateName,
		"TagName":       tagName,
	}

	param, err := utils.HmacSha1Base64(p.accessSecret, http.MethodPost, params)
	if err != nil {
		return
	}

	return utils.HttpRequest(p.host, http.MethodPost, param)
}
