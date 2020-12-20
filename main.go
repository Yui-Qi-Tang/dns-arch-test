package main

import (
	"log"
	"net"

	"simple-dns.yuki.org/answer"

	"github.com/miekg/dns"
)

type myDNS struct {
	storage    string
	dnskeyPriv string
	dnsKeyPub  string
}

func setup() {
	answer.AddTypeA("www.google.com.", "1.2.3.4")
}

const addr = ":8553"

type handle string

func (h *handle) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r) // set question fields and set success reply code
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := answer.GetTypeA(domain)
		if ok { // ok
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		} // else failed
	}
	w.WriteMsg(&msg)
}

func main() {
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion("miek.nl.", dns.TypeMX)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		return
	}
	if r.Rcode != dns.RcodeSuccess {
		return
	}
	for _, a := range r.Answer {
		if mx, ok := a.(*dns.MX); ok {
			log.Printf("%s\n", mx.String())
		}
	}
}
