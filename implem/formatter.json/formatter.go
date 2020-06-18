package formatter

import (
	"time"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

type jsonFormatter struct {
}

func NewJSONFormatter() interactor.Formatter {
	return &jsonFormatter{}
}

type Self struct {
	ID    string `json:"id"`
	Email string `json:"email"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}

func (jsonFormatter) Self(u *models.User) interface{} {
	return Self{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		LastLogin: u.LastLogin,
	}
}

type Service struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	TRCount     int         `json:"timed_response_count"`
	HealthCheck HealthCheck `json:"healthcheck"`
}

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

func (jsonFormatter) Service(svc *models.Service, trcount int) interface{} {
	return Service{
		ID:          svc.ID,
		Name:        svc.Name,
		Description: svc.Description,
		TRCount:     trcount,
		HealthCheck: HealthCheck{
			At:             svc.HealthCheck.At,
			URL:            svc.HealthCheck.URL,
			Every:          svc.HealthCheck.Every,
			DNS:            svc.HealthCheck.DNS,
			Handshake:      svc.HealthCheck.Handshake,
			TotalResponse:  svc.HealthCheck.TotalResponse,
			ServerResponse: svc.HealthCheck.ServerResponse,
			Connect:        svc.HealthCheck.Connect,
			Status:         svc.HealthCheck.Status,
		},
	}
}
