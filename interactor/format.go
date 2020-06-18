package interactor

import (
	"errors"
	"fmt"

	"github.com/Depado/vuemonit/models"
)

func (i Interactor) FormatSelf(user *models.User) interface{} {
	return i.Formatter.Self(user)
}

func (i Interactor) FormatService(svc *models.Service) (interface{}, error) {
	trcount, err := i.Store.CountTimedResponses(svc)
	if err != nil && !errors.Is(err, ErrNotFound) {
		i.Logger.Err(err).Msg("unable to count timed responses")
	}
	return i.Formatter.Service(svc, trcount), nil
}

func (i Interactor) GetServiceByID(user *models.User, id string) (interface{}, error) {
	s, err := i.Store.GetServiceByID(id)
	if err != nil {
		return nil, fmt.Errorf("unable to get service: %w", err)
	}
	if s.UserID != user.ID {
		return nil, ErrPermission
	}
	return i.FormatService(s)
}
