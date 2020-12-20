package dnshandler

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

var defaultType = map[string]string{
	"example.com.": "192.168.100.1",
}

// DefaultHandler ...
type DefaultHandler string

// ServeTypeA ...
func (d DefaultHandler) ServeTypeA(w dns.ResponseWriter, r *dns.Msg) {
	// just reply now
	m := dns.Msg{}
	m.SetReply(r)

	ans, ok := defaultType[r.Question[0].Name]
	if ok {
		m.Answer = []dns.RR{
			&dns.A{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(ans),
			},
		}
	}
	if err := w.WriteMsg(&m); err != nil {
		log.Println(err)
	}

}

// ServeDNS ...
func (d DefaultHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	switch r.Question[0].Qtype {
	case dns.TypeA:
		d.ServeTypeA(w, r)
	default:
		// no serve
	}
}
