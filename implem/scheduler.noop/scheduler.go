package scheduler

import (
	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

type noopSched struct{}

func NewNoopScheduler() interactor.Scheduler {
	return &noopSched{}
}
func (noopSched) Start(svc *models.Service) error {
	return nil
}

func (noopSched) Restart(svc *models.Service) error {
	return nil
}
