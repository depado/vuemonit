package storage

import (
	"errors"
	"fmt"

	"github.com/rs/xid"

	"github.com/Depado/vuemonit/models"
)

func (s StormStorage) SaveService(user *models.User, svc *models.Service) error {
	if user.ID == "" {
		return errors.New("user has no id")
	}
	if svc.ID == "" {
		svc.ID = xid.New().String()
	}
	svc.UserID = user.ID
	if err := s.db.Save(svc); err != nil {
		return fmt.Errorf("save service: %w", err)
	}
	s.log.Debug().Str("id", svc.ID).Msg("service saved")
	return nil
}

func (s StormStorage) SaveRawService(svc *models.Service) error {
	if svc.UserID == "" {
		return errors.New("user has no id")
	}
	if svc.ID == "" {
		return errors.New("service has no id")
	}
	if err := s.db.Save(svc); err != nil {
		return fmt.Errorf("save service: %w", err)
	}
	s.log.Debug().Str("id", svc.ID).Msg("service saved")
	return nil
}

func (s StormStorage) GetServiceByID(id string) (*models.Service, error) {
	svc := &models.Service{}
	if err := s.db.One("ID", id, svc); err != nil {
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return svc, nil
}

func (s StormStorage) GetAllServices() ([]*models.Service, error) {
	svx := []*models.Service{}
	if err := s.db.All(&svx); err != nil {
		return nil, fmt.Errorf("get all services: %w", err)
	}
	return svx, nil
}
