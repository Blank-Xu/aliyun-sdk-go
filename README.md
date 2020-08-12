# Aliyun-sdk

[![Go Report Card](https://goreportcard.com/badge/github.com/Blank-Xu/aliyun-sdk-go)](https://goreportcard.com/report/github.com/Blank-Xu/aliyun-sdk-go)
[![PkgGoDev](https://pkg.go.dev/badge/Blank-Xu/aliyun-sdk-go)](https://pkg.go.dev/github.com/Blank-Xu/aliyun-sdk-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

---

An easy way to use Aliyun [DirectMail](https://help.aliyun.com/document_detail/29434.html) to send emails,

use Aliyun [Short Message Service](https://help.aliyun.com/document_detail/101300.html) to send messages.

## Installation

    go get github.com/Blank-Xu/aliyun-sdk-go
    
## Simple Example
```go
package main

import (
	directmail "github.com/Blank-Xu/aliyun-sdk-go/directmail"
	sms "github.com/Blank-Xu/aliyun-sdk-go/sms"
	
	"fmt"
)

func main() {
	// email
	directmailAPI := directmail.NewAPI(host, accessKeyID, accessSecret, regionID, version, accountName, fromAlias)

	code, body, err := directmailAPI.SingleSendMail(emails, subject, "", textBody)
	fmt.Printf("code: %d, body: %s, err: %v\n", code, body, err)

	// sms
	smsAPI := sms.NewAPI(host, accessKeyID, accessSecret, regionID, version)
    
	resp, err := smsAPI.SendSms(phones, signName, templateCode, templateParam)
	fmt.Printf("resp: %+v, err: %v\n", resp, err)
}
```

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.