package interactor

import (
	"fmt"
	"time"

	"github.com/Depado/vuemonit/models"
)

func (i Interactor) NewService(user *models.User, name, description, url string) (*models.Service, error) {
	s, err := models.NewService(user, url, name, description, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("create new service: %w", err)
	}
	if err = i.Store.SaveService(user, s); err != nil {
		return nil, fmt.Errorf("save new service: %w", err)
	}
	if err = i.Scheduler.Start(s); err != nil {
		return nil, fmt.Errorf("unable to start routine: %w", err)
	}
	return s, nil
}
