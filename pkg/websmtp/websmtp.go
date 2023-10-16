package websmtp

import (
	"github.com/emersion/go-smtp"
	"io"
	"fmt"
	"net"
	"net/mail"
	"strings"
	"time"
)

// SendRequest is a request to send an email to the specified 'to' addresses.
type SendRequest struct {
	To      []string `json:"to" binding:"required"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
}

// SendStatus tracks the status of a SendRequest. If status is "completed", then
// Failed are the destination addresses that could not be reached.
// TODO: keep error messages on failed addresses.
type SendStatus struct {
	ID     string   `json:"id" binding:"required"`
	Status string   `json:"status" binding:"required"`
	Failed []string `json:"failed" binding:"required"`
}

// FindHost looks for an email address' SMTP host by looking up DNS records and
// knocking on MX record entries. Returns empty string if not found.
func FindHost(addr string) string {
	if _, err := mail.ParseAddress(addr); err != nil {
		return ""
	}

	domain := strings.Split(addr, "@")
	if len(domain) != 2 {
		return ""
	}

	// Retrieve MX DNS records.
	mxrecords, err := net.LookupMX(domain[1])
	if err != nil || len(mxrecords) == 0 {
		return ""
	}

	// Knock on DNS mx records port 25.
	for _, mx := range mxrecords {
		host := mx.Host + ":25"

		// TODO: bottleneck. Maybe implement caching?
		conn, err := net.DialTimeout("tcp", host, 60*time.Millisecond)

		if err == nil {
			conn.Close()
			return host
		}
	}

	return ""
}

// SendMail fulfills the SendRequest. This might take a while to complete if
// you specify too many email addresses. If you need immediate responses, use
// Sender struct instead.
func SendMail(id, src string, req SendRequest) SendStatus {
	failed := []string{}
	for _, dest := range req.To {
		host := FindHost(dest)
		if host == "" {
			failed = append(failed, dest)
			continue
		}
		if err := smtp.SendMail(host, nil, src, []string{dest},
			getData(id, dest, src, req)); err != nil {
			failed = append(failed, dest)
		}
	}
	return SendStatus{id, "completed", failed}
}

const dataFmt = "Message-Id: <%s@refid.com>\r\nTo: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n %s"

func getData(id, dest, src string, req SendRequest) io.Reader {
	return strings.NewReader(fmt.Sprintf(dataFmt, id, dest, src,
		req.Subject, req.Body))
}
