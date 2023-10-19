package websmtp

import (
	"sync"
	"testing"
	"time"
)

func TestSenderClear(t *testing.T) {
	sender := NewSender()
	ref := sender.Enqueue(SendRequest{
		"source@example.com",
		[]string{"bar@example.com"},
		"Email subject",
		"Lorem ipsum dolor",
	})

	// If cache is cleared, then we shouldn't be able to get the
	// request's status.
	sender.Clear()
	if status := sender.GetStatus(ref); status.Status != "not started" {
		t.Fatal("cache is not being cleared")
	}

	// If cache is cleared and sender is running, then the request
	// should eventually show up in the cache (because its still in
	// queue).
	go sender.Start(1)
	time.Sleep(150 * time.Millisecond)
	sender.Stop()
	if status := sender.GetStatus(ref); status.Status == "not started" {
		t.Fatal("clearing status cache should not affect the queue")
	}
}

func TestSenderClearConcurrency(t *testing.T) {
	sender := NewSender()
	signalChan := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(dowork chan struct{}) {
			defer wg.Done()
			<-dowork
			sender.Clear()
		}(signalChan)
	}
	close(signalChan)
	wg.Wait()
}

func TestSenderEnqueue(t *testing.T) {
	sender := NewSender()
	if ref := sender.Enqueue(SendRequest{"bad_example",
		[]string{"myself@example.com"}, "foo bar",
		"hello world"}); ref != "" {
		t.Errorf("enqueue should not allow invalid req.From address")
	}

	ref := sender.Enqueue(SendRequest{"good@example.com",
		[]string{"dest@example.com"}, "foo bar", "hello world"})

	if ref == "" {
		t.Fatalf("enqueue rejected valid SendRequest")
	}

	if len(sender.sendQueue) != 1 {
		t.Fatalf("sendRequest not enqueued")
	}

	if sender.sendQueue[0].Req.From != "good@example.com" {
		t.Fatalf("SendRequest not enqueued in order")
	}
}

func TestSenderEnqueueConcurrency(t *testing.T) {
	sender := NewSender()
	signalChan := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(dowork chan struct{}) {
			defer wg.Done()
			<-dowork
			sender.Enqueue(SendRequest{})
		}(signalChan)
	}
	close(signalChan)
	wg.Wait()
}

func BenchmarkSenderEnqueue(b *testing.B) {
	sender := NewSender()
	req := SendRequest{"from@example.com", []string{"to@example.com"},
		"subject", "body"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sender.Enqueue(req)
	}
}

func TestSenderStartStop(t *testing.T) {
	sender := NewSender()
	err := sender.Start(3)
	if err != nil {
		t.Error(err)
	}
	if !sender.started {
		t.Fatal("sender 'started' attribute remained false")
	}

	// Cannot start whilst already started.
	err = sender.Start(1)
	if err == nil {
		t.Fatal("no error returned when calling Start() multiple times")
	}

	// Stop all workers.
	err = sender.Stop()
	if err != nil {
		t.Error(err)
	}
	if sender.started {
		t.Fatal("sender 'started' attribute remained true after Stop()")
	}

	// Stop should fail when sender is not started.
	err = sender.Stop()
	if err == nil {
		t.Fatal("no error returned when calling Stop() multiple times")
	}
}

func BenchmarkSenderStartStop(b *testing.B) {
	sender := NewSender()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sender.Start(5)
		sender.Stop()
	}
}
