package storage

import (
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

func (s StormStorage) SaveUser(usr *models.User) error {
	n := time.Now()
	if usr.ID == "" {
		usr.ID = xid.New().String()
	}
	if err := s.db.Save(usr); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	s.log.Debug().Str("id", usr.ID).Dur("took", time.Since(n)).Msg("user saved")
	return nil
}

func (s StormStorage) GetUserByEmail(email string) (*models.User, error) {
	n := time.Now()
	u := &models.User{}
	if err := s.db.One("Email", email, u); err != nil {
		return nil, fmt.Errorf("find by email: %v: %w", err, interactor.ErrNotFound)
	}
	s.log.Debug().Str("id", u.ID).Dur("took", time.Since(n)).Msg("retrieved by email")
	return u, nil
}

func (s StormStorage) GetUserByID(id string) (*models.User, error) {
	n := time.Now()
	u := &models.User{}
	if err := s.db.One("ID", id, u); err != nil {
		return nil, fmt.Errorf("find by id: %v: %w", err, interactor.ErrNotFound)
	}
	s.log.Debug().Str("id", u.ID).Dur("took", time.Since(n)).Msg("retrieved by id")
	return u, nil
}
