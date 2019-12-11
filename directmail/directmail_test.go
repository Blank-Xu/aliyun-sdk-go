package directmail

import (
	"testing"
)

const (
	tHost         = ""
	tAccessKeyId  = ""
	tAccessSecret = ""
	tRegionId     = ""
	tVersion      = ""
	tAccountName  = ""
	tFromAlias    = ""
)

var (
	tApi = NewApi(tHost, tAccessKeyId, tAccessSecret, tRegionId, tVersion, tAccountName, tFromAlias)

	tEmails   = []string{"test@example.com"}
	tSubject  = "ali email test"
	tTextBody = "text body"
	tHtmlBody = `<!DOCTYPE html>
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
	code, resp, err := tApi.SingleSendMail(tEmails, tSubject, tHtmlBody, tTextBody)
	if err != nil || code >= 300 {
		t.Fatalf("send failed, code: %d, resp: %s, err: %v", code, resp, err)
	}

	t.Logf("send success, resp: %s", resp)
}

func TestBatchSendMail(t *testing.T) {
	// TODO:
}
