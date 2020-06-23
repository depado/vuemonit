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

func (i Interactor) GetServices(user *models.User) ([]*models.Service, error) {
	svx, err := i.Store.GetServices(user)
	if err != nil {
		return nil, fmt.Errorf("unable to get services: %w", err)
	}
	return svx, nil
}

func (i Interactor) GetServiceByID(user *models.User, id string) (*models.Service, error) {
	s, err := i.Store.GetServiceByID(id)
	if err != nil {
		return nil, fmt.Errorf("unable to get service: %w", err)
	}
	if s.UserID != user.ID {
		return nil, ErrPermission
	}
	return s, nil
}
