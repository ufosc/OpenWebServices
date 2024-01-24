package authapi

import "github.com/ufosc/OpenWebServices/pkg/websmtp"

// MailSender encapsulates a websmtp.Sender object.
type MailSender struct {
	sender *websmtp.Sender
	addr   string
}

// NewMailSender creates a new MailSender with the given config.
func NewMailSender(addr string) MailSender {
	return MailSender{
		websmtp.NewSender(25),
		addr,
	}
}

// Start the mail sender project with 'n' threads.
func (ms MailSender) Start(n int) error {
	return ms.sender.Start(n)
}

// Stop the mail sender process.
func (ms MailSender) Stop() error {
	return ms.sender.Stop()
}

// SendVerification sends the signup verification email, where id is the
// MongoDB object ID of the pending user and email is their email address.
func (ms MailSender) SendVerification(id, email string) bool {
	ok := ms.sender.Enqueue(websmtp.SendRequest{
		From:    ms.addr,
		To:      []string{email},
		Subject: "UF Open Source Club: Verify Your Email Address",
		Body:    "go to api.ufosc.org/auth/verify/" + id,
	})

	return ok != ""
}
