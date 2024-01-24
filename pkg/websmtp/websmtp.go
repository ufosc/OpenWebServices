package websmtp

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"io"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"time"
)

// hostCache caches MX records in memory to speed up looking for email hosts.
var hostCache = map[string]string{}
var cacheMutex sync.Mutex

// DefaultPort is the default outbound SMTP port.
var DefaultPort int = 25

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

// RelayConfig encapsulates the connection details for an SMTP relay node. This is
// for situations where a hosting provider blocks port 25, but you're still interested
// in using WebSMTP.
type RelayConfig struct {
	Host      string // The SMTP host address.
	Port      int    // The SMTP server port.
	Username  string // The server username.
	Password  string // The server password.
	UseSecure bool   // Whether to use TLS.
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

	res := mxrecords[0].Host
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	hostCache[domain[1]] = res
	return res
}

// SendMail fulfills the SendRequest. This might take a while to complete if
// you specify too many email addresses. If you need immediate responses, use
// Sender struct instead.
func SendMail(id string, req SendRequest, port int) SendStatus {
	failed := []string{}
	for _, dest := range req.To {
		host := FindHost(dest) + ":" + strconv.FormatInt(int64(port), 10)
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
			fmt.Println(err)
		}
	}
	return SendStatus{id, "completed", failed, time.Now().Unix()}
}

// SendMailWithRelay fulfills the SendRequest via a different SMTP relay.
func SendMailWithRelay(id string, req SendRequest, c RelayConfig) SendStatus {
	failed := []string{}
	host := c.Host + ":" + strconv.FormatInt(int64(c.Port), 10)
	eid := uuid.New().String()
	msg := getData(eid, req.From, req.From, req)

	var client sasl.Client
	if c.Username != "" {
		client = sasl.NewLoginClient(c.Username, c.Password)
	}

	if c.UseSecure {
		if err := smtp.SendMailTLS(host, client, req.From,
			req.To, msg); err != nil {
			failed = req.To
			fmt.Println(err)
		}
		return SendStatus{id, "completed", failed, time.Now().Unix()}
	}

	if err := smtp.SendMail(host, client, req.From, req.To, msg); err != nil {
		failed = req.To
		fmt.Println(err)
	}

	return SendStatus{id, "completed", failed, time.Now().Unix()}
}

const dataFmt = "Message-Id: <%s@refid.com>\r\nContent-Type: text/html\r\nTo: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n %s"

func getData(id, dest, src string, req SendRequest) io.Reader {
	return strings.NewReader(fmt.Sprintf(dataFmt, id, dest, src,
		req.Subject, req.Body))
}
