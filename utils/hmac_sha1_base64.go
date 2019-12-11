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

func HmacSha1Base64(accessSecret []byte, httpMethod string, params map[string]string) (string, error) {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.Grow(128)
	for _, key := range keys {
		buf.WriteString(QueryEscape(key))
		buf.WriteByte('=')
		buf.WriteString(QueryEscape(params[key]))
		buf.WriteByte('&')
	}
	buf.Truncate(buf.Len() - 1)
	param := buf.String()

	var signatureBuf bytes.Buffer
	buf.Grow(256)
	signatureBuf.WriteString(httpMethod)
	// QueryEscape("&/&") = "&%2F&"
	signatureBuf.WriteString("&%2F&")
	signatureBuf.WriteString(QueryEscape(param))

	// HmacSHA1
	hmacSha1 := hmac.New(crypto.SHA1.New, accessSecret)
	if _, err := hmacSha1.Write(signatureBuf.Bytes()); err != nil {
		return "", err
	}

	// Base64
	signature := base64.StdEncoding.EncodeToString(hmacSha1.Sum(nil))
	// 签名需要特殊URL编码
	signature = QueryEscape(signature)

	param += "&Signature=" + signature

	return param, nil
}

func QueryEscape(param string) string {
	param = url.QueryEscape(param)

	param = strings.Replace(param, "+", "%20", -1)
	param = strings.Replace(param, "*", "%2A", -1)
	param = strings.Replace(param, "%7E", "~", -1)

	return param
}
