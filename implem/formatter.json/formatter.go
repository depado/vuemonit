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
	return &Self{
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
	Count       int         `json:"count"`
	HealthCheck HealthCheck `json:"healthcheck"`
}

type Services []*Service

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

func (jsonFormatter) service(svc *models.Service) *Service {
	return &Service{
		ID:          svc.ID,
		Name:        svc.Name,
		Description: svc.Description,
		Count:       svc.Count,
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

func (j jsonFormatter) Service(svc *models.Service) interface{} {
	return j.service(svc)
}

func (j jsonFormatter) Services(svc []*models.Service) interface{} {
	svx := Services{}
	for _, s := range svc {
		svx = append(svx, j.service(s))
	}
	return svx
}

type TimedResponse struct {
	ID     string        `json:"id"`
	At     time.Time     `json:"at"`
	Server time.Duration `json:"server"`
	Total  time.Duration `json:"total"`
	Status int           `json:"status"`
}

type TimedResponses []*TimedResponse

func (jsonFormatter) timedResponse(tr *models.TimedResponse) *TimedResponse {
	return &TimedResponse{
		ID:     tr.ID.String(),
		At:     tr.At,
		Server: tr.Server,
		Total:  tr.Total,
		Status: tr.Status,
	}
}

func (j jsonFormatter) TimedResponses(tr []*models.TimedResponse) interface{} {
	trx := TimedResponses{}
	for _, t := range tr {
		trx = append(trx, j.timedResponse(t))
	}
	return trx
}
