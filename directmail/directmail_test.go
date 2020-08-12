package directmail

import (
	"testing"
)

const (
	host         = ""
	accessKeyID  = ""
	accessSecret = ""
	regionID     = ""
	version      = ""
	accountName  = ""

	fromAlias = ""

	receiversName = ""
	templateName  = ""
	tagName       = ""
)

var (
	api = NewAPI(host, accessKeyID, accessSecret, regionID, version, accountName, fromAlias)

	emails   = []string{"test@example.com"}
	subject  = "Aliyun directmail test"
	textBody = "text body"
	htmlBody = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Aliyun DirectMail</title>
</head>
<body>
html body
</body>
</html>`
)

func TestSingleSendMail(t *testing.T) {
	code, resp, err := api.SingleSendMail(emails, subject, htmlBody, textBody)
	if err != nil || code != 200 {
		t.Fatalf("send failed, code: %d, resp: %s, err: %v", code, resp, err)
	}

	t.Logf("send success, resp: %s", resp)
}

func TestBatchSendMail(t *testing.T) {
	code, resp, err := api.BatchSendMail(receiversName, templateName, tagName)
	if err != nil || code != 200 {
		t.Fatalf("send failed, code: %d, resp: %s, err: %v", code, resp, err)
	}

	t.Logf("send success, resp: %s", resp)
}
