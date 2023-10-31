package main

import (
	"github.com/ufosc/OpenWebServices/pkg/websmtp"
)

type MailSender struct {
	sender *websmtp.Sender
	config Config
}

func NewMailSender(config Config) MailSender {
	return MailSender{
		websmtp.NewSender(),
		config,
	}
}

func (ms MailSender) Start(n int) error {
	return ms.sender.Start(n)
}

func (ms MailSender) Stop() error {
	return ms.sender.Stop()
}

func (ms MailSender) SendSignupVerification(id string, usr PendingUserModel) bool {
	ok := ms.sender.Enqueue(websmtp.SendRequest{
		From: ms.config.NOTIF_EMAIL_ADDR,
		To: []string{usr.User.Email},
		Subject: "Verify Email Address: Welcome to the Open Source Club @ UF",
		Body: "go to api.ufosc.org/auth/verify?type=signup&ref=" + id,
	})

	if ok == "" {
		return false
	}

	return true
}
