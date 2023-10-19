package websmtp

import (
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"io"
	"net"
	"net/mail"
	"strings"
	"sync"
	"time"
)

// hostCache caches MX records in memory to speed up looking for email hosts.
var hostCache = map[string]string{}
var cacheMutex sync.Mutex

// SendRequest is a request to send an email to the specified 'to' addresses.
type SendRequest struct {
	From    string   `json:"from" binding:"required"`    // Email source address.
	To      []string `json:"to" binding:"required"`      // List of destination email addresses.
	Subject string   `json:"subject" binding:"required"` // Email subject.
	Body    string   `json:"body" binding:"required"`    // Email body.
}

// SendStatus tracks the status of a SendRequest. If status is "completed", then
// Failed are the destination addresses that could not be reached.
type SendStatus struct {
	ID            string   `json:"id" binding:"required"`     // Email Message-ID header.
	Status        string   `json:"status" binding:"required"` // Request status.
	Failed        []string `json:"failed" binding:"required"` // List of addresses that couldn't be reached.
	TimeCompleted int64    `json:"time_completed"`            // Timestamp at SendRequest completion.
}

// FindHost looks for an email address' SMTP host by looking up DNS records and
// Returns empty string if nothing was found.
func FindHost(addr string) string {
	if _, err := mail.ParseAddress(addr); err != nil {
		return ""
	}

	domain := strings.Split(addr, "@")
	if len(domain) != 2 {
		return ""
	}

	// Try hitting cache before checking DNS.
	if val, found := hostCache[domain[1]]; found {
		return val
	}

	// Retrieve MX DNS records.
	mxrecords, err := net.LookupMX(domain[1])
	if err != nil || len(mxrecords) == 0 {
		return ""
	}

	res := mxrecords[0].Host + ":25"
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	hostCache[domain[1]] = res
	return res
}

// SendMail fulfills the SendRequest. This might take a while to complete if
// you specify too many email addresses. If you need immediate responses, use
// Sender struct instead.
func SendMail(id string, req SendRequest) SendStatus {
	failed := []string{}
	for _, dest := range req.To {
		host := FindHost(dest)
		if host == "" {
			failed = append(failed, dest)
			continue
		}
		// Use a separate, random Message-ID so that a recipient cannot
		// reverse-lookup the SendRequest that sent him the email.
		eid := uuid.New().String()
		if err := smtp.SendMail(host, nil, req.From, []string{dest},
			getData(eid, dest, req.From, req)); err != nil {
			failed = append(failed, dest)
		}
	}
	return SendStatus{id, "completed", failed, time.Now().Unix()}
}

const dataFmt = "Message-Id: <%s@refid.com>\r\nContent-Type: text/html\r\nTo: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n %s"

func getData(id, dest, src string, req SendRequest) io.Reader {
	return strings.NewReader(fmt.Sprintf(dataFmt, id, dest, src,
		req.Subject, req.Body))
}
