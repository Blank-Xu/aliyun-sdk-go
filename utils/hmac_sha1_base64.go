package utils

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

// HmacSha1Base64  for aliyun sdk
func HmacSha1Base64(accessSecret []byte, httpMethod string, params [][2]string) (string, error) {
	// sort first
	sort.Slice(params, func(i, j int) bool {
		return params[i][0] < params[j][0]
	})

	var paramBuf bytes.Buffer

	paramBuf.Grow(256)
	for _, param := range params {
		paramBuf.WriteString(param[0])
		paramBuf.WriteByte('=')
		paramBuf.WriteString(QueryEscape(param[1]))
		paramBuf.WriteByte('&')
	}

	l := paramBuf.Len()
	paramBuf.Truncate(l - 1)

	param := QueryEscape(paramBuf.String())

	var signatureBuf bytes.Buffer

	signatureBuf.Grow(l + 8)
	signatureBuf.WriteString(httpMethod)
	signatureBuf.WriteString("&%2F&")
	signatureBuf.WriteString(param)

	// HmacSHA1
	hmacSha1 := hmac.New(crypto.SHA1.New, accessSecret)
	if _, err := hmacSha1.Write(signatureBuf.Bytes()); err != nil {
		return "", err
	}

	// Base64
	signature := base64.StdEncoding.EncodeToString(hmacSha1.Sum(nil))
	// 签名需要特殊URL编码
	signature = QueryEscape(signature)

	paramBuf.WriteString("&Signature=")
	paramBuf.WriteString(signature)

	return paramBuf.String(), nil
}

// QueryEscape  .
func QueryEscape(param string) string {
	if param == "" {
		return ""
	}

	param = url.QueryEscape(param)

	param = strings.Replace(param, "+", "%20", -1)
	param = strings.Replace(param, "*", "%2A", -1)
	param = strings.Replace(param, "%7E", "~", -1)

	return param
}
