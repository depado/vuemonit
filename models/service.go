package models

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"github.com/rs/xid"
)

type HealthCheck struct {
	At             time.Time     `json:"at"`
	URL            string        `json:"url"`
	Every          time.Duration `json:"every"`
	DNS            time.Duration `json:"dns"`
	Handshake      time.Duration `json:"handshake"`
	Connect        time.Duration `json:"connect"`
	TotalResponse  time.Duration `json:"total"`
	ServerResponse time.Duration `json:"server"`
	Status         int           `json:"status"`
}

type Service struct {
	ID          string      `json:"id" storm:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	HealthCheck HealthCheck `json:"healthcheck"`
	UserID      string      `json:"user_id"`
}

func NewService(user *User, hurl, name, description string, every time.Duration) (*Service, error) {
	if user.ID == "" {
		return nil, fmt.Errorf("user has no id")
	}
	su, err := url.Parse(hurl)
	if err != nil {
		return nil, fmt.Errorf("new service: %w", err)
	}

	return &Service{
		ID:          xid.New().String(),
		Name:        name,
		Description: description,
		HealthCheck: HealthCheck{
			URL:   su.String(),
			Every: every,
		},
		UserID: user.ID,
	}, nil
}

func (s *Service) Fetch() (*TimedResponse, error) {
	var start, c, d, t time.Time
	var co, dns, handshake bool

	su, err := url.Parse(s.HealthCheck.URL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	req, err := http.NewRequest("GET", su.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("fetch service: %w", err)
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { d = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dns = true
			s.HealthCheck.DNS = time.Since(d)
		},

		TLSHandshakeStart: func() { t = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			handshake = true
			s.HealthCheck.Handshake = time.Since(t)
		},

		ConnectStart: func(network, addr string) { c = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			co = true
			s.HealthCheck.Connect = time.Since(c)
		},

		GotFirstResponseByte: func() {
			took := time.Since(start)
			s.HealthCheck.ServerResponse = took
			s.HealthCheck.TotalResponse = took
			if dns {
				s.HealthCheck.ServerResponse -= s.HealthCheck.DNS
			} else {
				s.HealthCheck.TotalResponse += s.HealthCheck.DNS
			}
			if handshake {
				s.HealthCheck.ServerResponse -= s.HealthCheck.Handshake
			} else {
				s.HealthCheck.TotalResponse += s.HealthCheck.Handshake
			}
			if co {
				s.HealthCheck.ServerResponse -= s.HealthCheck.Connect
			} else {
				s.HealthCheck.TotalResponse += s.HealthCheck.Connect
			}
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		defer resp.Body.Close()
		return nil, fmt.Errorf("unable to roundtrip for service %v: %w", s.Name, err)
	}
	io.Copy(ioutil.Discard, resp.Body) // nolint: errcheck
	s.HealthCheck.Status = resp.StatusCode
	s.HealthCheck.At = time.Now()
	return &TimedResponse{
		ID:        xid.New(),
		Server:    s.HealthCheck.ServerResponse,
		Total:     s.HealthCheck.TotalResponse,
		Status:    s.HealthCheck.Status,
		ServiceID: s.ID,
	}, nil
}
