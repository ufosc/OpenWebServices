package websmtp

import (
	"fmt"
	"github.com/google/uuid"
	"net/mail"
	"sync"
	"time"
)

// sendObj wraps a SendRequest with a unique ID.
type sendObj struct {
	ID  string
	Req SendRequest
}

// Sender maintains a queue of send requests and their status, and attempts to
// process them.
type Sender struct {
	sendQueue   []sendObj
	status      map[string]SendStatus
	port        int
	relayConfig *RelayConfig
	started     bool
	mutex       sync.Mutex
	wg          sync.WaitGroup
	ch          chan struct{}
}

// NewSender initializes a new sender.
func NewSender(port int) *Sender {
	sender := new(Sender)
	sender.sendQueue = []sendObj{}
	sender.status = map[string]SendStatus{}
	sender.started = false
	sender.relayConfig = nil
	sender.port = port
	return sender
}

// NewSenderRelay initializes a new sender that sends via SMTP relay.
func NewSenderRelay(port int, relay *RelayConfig) *Sender {
	sender := NewSender(port)
	sender.relayConfig = relay
	return sender
}

// Clear the SendRequest status cache. Queued requests are not affected, and will
// show up in the status queue once they're finished processing.
func (s *Sender) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.status = map[string]SendStatus{}
}

// Enqueue a request to the processing queue and return a reference ID.
// Returns an empty string if req.From is not a valid email address.
func (s *Sender) Enqueue(req SendRequest) string {
	id := uuid.New().String()
	if _, err := mail.ParseAddress(req.From); err != nil {
		return ""
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sendQueue = append(s.sendQueue, sendObj{id, req})
	return id
}

// GetStatus returns the status of a request by looking up its reference ID.
func (s *Sender) GetStatus(ref string) SendStatus {
	if val, found := s.status[ref]; found {
		return val
	}
	return SendStatus{ref, "not started", []string{}, 0}
}

// sendMail fulfills the sendRequest in the given entry and updates the
// request's status.
func (s *Sender) sendMail(entry sendObj) {
	var status SendStatus
	if s.relayConfig != nil {
		status = SendMailWithRelay(entry.ID, entry.Req, *s.relayConfig)
	} else {
		status = SendMail(entry.ID, entry.Req, s.port)
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.status[entry.ID] = status
}

// doStart is a single thread processing the SendQueue.
func (s *Sender) doStart() {
	s.wg.Add(1)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if len(s.sendQueue) <= 0 {
				continue
			}

			s.mutex.Lock()
			req := s.sendQueue[0]
			s.sendQueue = append([]sendObj{}, s.sendQueue[1:]...)
			s.status[req.ID] = SendStatus{req.ID, "processing",
				[]string{}, 0}

			s.mutex.Unlock()
			s.sendMail(req)
		case <-s.ch:
			s.wg.Done()
			return
		}
	}
}

// Start spawns n threads that process the SendRequest queue. Each thread will
// process one request at a time. Call Stop() to kill the workers.
func (s *Sender) Start(n int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.started {
		return fmt.Errorf("sender already started")
	}
	s.started = true
	s.ch = make(chan struct{})
	for i := 0; i < n; i++ {
		go s.doStart()
	}
	return nil
}

// Stop any threads currently processing the queue. Returns an error if the
// sender was never started.
func (s *Sender) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.started {
		return fmt.Errorf("sender has not started")
	}
	s.started = false
	close(s.ch)
	s.wg.Wait()
	return nil
}
