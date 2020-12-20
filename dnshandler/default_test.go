package dnshandler

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/phayes/freeport"
)

func TestQueryTypeAWithDefaultHandler(t *testing.T) {
	t.Log("start testing query type A with default handler")
	testcases := []struct {
		name string
		ans  string
	}{
		{
			name: "example.com.",
			ans:  "192.168.100.1",
		},
	}

	chErr := make(chan error)

	// port
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	var h DefaultHandler = "default handler"
	srv.Handler = h
	// server
	t.Log("server is running at:", srv.Addr)
	go func() {
		srv.ListenAndServe()
	}()

	<-time.After(1 * time.Second) // wait for server ok

	// client
	go func() {
		c := &dns.Client{}
		for _, testcase := range testcases {
			qa := &dns.Msg{} // question typeA
			qa.SetQuestion(testcase.name, dns.TypeA)
			reply, rtt, err := c.Exchange(qa, srv.Addr)
			if err != nil {
				chErr <- err
				break
			}
			v, ok := reply.Answer[0].(*dns.A)
			if !ok {
				chErr <- errors.New("wrong type from server:" + dns.TypeToString[v.Hdr.Rrtype])
				break
			}
			if v.A.String() != testcase.ans {
				chErr <- errors.New("wrong response of Type A query, it should be: " + testcase.ans + " bot got: " + v.A.String())
				break
			}
			t.Log("query type A", testcase.name, "done in", rtt)
		}
		close(chErr)
		return
	}()

	for e := range chErr {
		if e != nil {
			t.Fatal(e)
		}
	}

	t.Log("server is exiting...")
	if err := srv.Shutdown(); err != nil {
		t.Fatal(err)
	}
	t.Log("... Passed")
}
