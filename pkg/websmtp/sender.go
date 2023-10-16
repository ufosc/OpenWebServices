package websmtp

import (
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"net/mail"
	"sync"
)

// sendObj wraps a SendRequest with a unique ID.
type sendObj struct {
	ID  string
	Req SendRequest
}

// Sender maintains a queue of send requests and their status, and attempts to
// process them.
type Sender struct {
	// TODO: need to decide on restrictions. Clients are essentially
	// sending emails using your address.
	src       string
	sendQueue []sendObj

	// TODO: status needs to be cleared occasionally.
	status map[string]SendStatus

	mutex sync.Mutex
}

// NewSender initializes a new sender with the given src string. Returns nil if
// src is an invalid email address.
func NewSender(src string) *Sender {
	if _, err := mail.ParseAddress(src); err != nil {
		return nil
	}
	sender := new(Sender)
	sender.src = src
	sender.sendQueue = []sendObj{}
	sender.status = map[string]SendStatus{}
	return sender
}

// Enqueue a request to the processing queue and return a reference ID.
func (s *Sender) Enqueue(req SendRequest) string {
	id := uuid.New().String()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sendQueue = append(s.sendQueue, sendObj{id, req})
	s.status[id] = SendStatus{id, "pending", []string{}}
	return id
}

// GetStatus returns the status of a request by looking up its reference ID.
func (s *Sender) GetStatus(ref string) SendStatus {
	if val, found := s.status[ref]; found {
		return val
	}
	return SendStatus{ref, "not found", []string{}}
}

// sendMail fulfills the sendRequest in the given entry and updates the
// request's status.
func (s *Sender) sendMail(entry sendObj) {
	failed := []string{}
	for _, dest := range entry.Req.To {
		host := FindHost(dest)
		if host == "" {
			failed = append(failed, dest)
			continue
		}
		if err := smtp.SendMail(host, nil, s.src, []string{dest},
			getData(entry.ID, dest, s.src, entry.Req)); err != nil {
			failed = append(failed, dest)
		}
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.status[entry.ID] = SendStatus{
		ID:     entry.ID,
		Status: "completed",
		Failed: failed,
	}
}

// Run begins processing the underlying send queue, one request at a time.
// It is meant to be ran as a goroutine.
// TODO: should take context as parameter, allow stopping, multiple threads.
func (s *Sender) Run() {
	for true {
		if len(s.sendQueue) <= 0 {
			continue
		}
		s.mutex.Lock()
		req := s.sendQueue[0]
		s.sendQueue = append([]sendObj{}, s.sendQueue[1:]...)
		s.mutex.Unlock()
		s.sendMail(req)
	}
}