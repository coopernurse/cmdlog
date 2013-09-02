package cmdlog

import (
	"fmt"
	"github.com/AeroNotix/libsmtp"
	"net/smtp"
)

type SMTPLogger struct {
	Host string
	Auth *smtp.Auth
	From string
	To   []string
}

func (s SMTPLogger) Log(result *Result, subject string, url string) error {
	return libsmtp.SendMailWithAttachments(s.Host, s.Auth, s.From, subject, s.To, toMsg(result, url), nil)
}

const msgTmpl string = `Job: %s
Date: %s
Elapsed seconds: %.2f
Host: %s
Exit Str: %s

%s`

func toMsg(result *Result, url string) []byte {
	return []byte(fmt.Sprintf(msgTmpl, result.Name, result.StartDate(dateTimeFmt), result.ElapsedSeconds(),
		result.Host, result.ExitStr, url))
}
