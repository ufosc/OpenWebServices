package main

import (
	"github.com/ufosc/OpenWebServices/pkg/websmtp"
)

// MailSender encapsulates a websmtp.Sender object.
type MailSender struct {
	sender *websmtp.Sender
	config Config
}

// NewMailSender creates a new MailSender with the given config.
func NewMailSender(config Config) MailSender {
	return MailSender{
		websmtp.NewSender(),
		config,
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

// SendSignupVerification sends the signup verification email for usr, where id
// is the MongoDB object ID of usr.
func (ms MailSender) SendSignupVerification(id string, usr PendingUserModel) bool {
	ok := ms.sender.Enqueue(websmtp.SendRequest{
		From:    ms.config.NOTIF_EMAIL_ADDR,
		To:      []string{usr.User.Email},
		Subject: "Verify Email Address: Welcome to the Open Source Club @ UF",
		Body:    "go to api.ufosc.org/auth/verify?type=signup&ref=" + id,
	})

	if ok == "" {
		return false
	}

	return true
}

// SendSigninVerification sends the signin verification email to the specified
// email address. 'id' is the MongoDB object ID for the corresponding
// VerifySigninModel object.
func (ms MailSender) SendSigninVerification(id, email string) bool {
	ok := ms.sender.Enqueue(websmtp.SendRequest{
		From:    ms.config.NOTIF_EMAIL_ADDR,
		To:      []string{email},
		Subject: "Verify Email Address: Sign in to Open Source Club",
		Body:    "go to api.ufosc.org/auth/verify?type=signin&ref=" + id,
	})

	if ok == "" {
		return false
	}

	return true
}
