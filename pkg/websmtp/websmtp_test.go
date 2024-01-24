package websmtp

import (
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"

// Its a good idea to test @ufl.edu as well, but sending TCP requests in rapid
// succession causes the server to block/timeout.
var hosts = []string{"gmail.com", "yahoo.com", "outlook.com", "hotmail.com",
	"aol.com", "msn.com", "yandex.ru", "yahoo.com.au", "voila.fr"}

func genRandEmail() string {
	a := []byte(alphabet)
	name := make([]byte, 8)
	for i := range name {
		name[i] = a[rand.Intn(len(a))]
	}
	domain := hosts[rand.Intn(len(hosts))]
	return string(name) + "@" + domain
}

func TestFindHost(t *testing.T) {
	for i := 0; i < 100; i++ {
		email := genRandEmail()
		host := FindHost(email) + ":25"
		if host == "" {
			t.Fatalf("failed to find host %s", email)
		}
		conn, err := net.DialTimeout("tcp", host, 2*time.Second)
		if err != nil {
			t.Fatalf("failed to reach %s at %s (%s)", email, host, err)
		}
		conn.Close()
	}
}

func TestFindHostBadAddr(t *testing.T) {
	if host := FindHost("foobar"); host != "" {
		t.Fatalf("found host for invalid email address")
	}
	if host := FindHost("foobar@baz"); host != "" {
		t.Fatalf("found host for invalid email address")
	}
}

func TestFindHostConcurrency(t *testing.T) {
	signalChan := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(dowork chan struct{}) {
			defer wg.Done()
			email := genRandEmail()
			<-dowork
			if str := FindHost(email); str == "" {
				t.Errorf("could not find host for %s", email)
			}

		}(signalChan)
	}
	close(signalChan)
	wg.Wait()
}

func BenchmarkFindHost(b *testing.B) {
	emails := []string{}
	for i := 0; i < 20; i++ {
		emails = append(emails, genRandEmail())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// All hosts exist, none should be nil.
		addr := emails[rand.Intn(len(emails))]
		if str := FindHost(addr); str == "" {
			b.Errorf("failed to find host %s", addr)
		}
	}
}

func TestSendMail(t *testing.T) {
	// TODO: need to set up a mock SMTP server but the smtp package needs
	// the server to be configured for TLS.
}
