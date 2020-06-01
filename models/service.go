package models

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

type Services []*Service

type ResponseInfo struct {
	DNS            time.Duration `json:"dns"`
	Handshake      time.Duration `json:"tls_handshake"`
	Connect        time.Duration `json:"connect"`
	TotalResponse  time.Duration `json:"total"`
	ServerResponse time.Duration `json:"server"`
}

type Service struct {
	URL      *url.URL     `json:"-"`
	Status   int          `json:"status"`
	Response ResponseInfo `json:"response"`
}

func NewService(u string) (*Service, error) {
	su, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("new service: %w", err)
	}

	return &Service{URL: su}, nil
}

func (s *Service) Fetch() error {
	var start, connect, dns, tlsHandshake time.Time

	req, err := http.NewRequest("GET", s.URL.String(), nil)
	if err != nil {
		return fmt.Errorf("fetch service: %w", err)
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			s.Response.DNS = time.Since(dns)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			s.Response.Handshake = time.Since(tlsHandshake)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			s.Response.Connect = time.Since(connect)
		},

		GotFirstResponseByte: func() {
			s.Response.TotalResponse = time.Since(start)
			s.Response.ServerResponse = s.Response.TotalResponse - s.Response.DNS - s.Response.Handshake - s.Response.Connect
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(ioutil.Discard, resp.Body) // nolint: errcheck
	defer resp.Body.Close()
	s.Status = resp.StatusCode

	return nil
}
