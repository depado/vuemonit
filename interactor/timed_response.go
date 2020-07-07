package interactor

import (
	"fmt"
	"time"

	"github.com/Depado/vuemonit/models"
)

func (i Interactor) GetTimedResponsesByServiceID(user *models.User, id string, limit int, reverse bool) ([]*models.TimedResponse, error) {
	s, err := i.GetServiceByID(user, id)
	if err != nil {
		return nil, fmt.Errorf("get service by id: %w", err)
	}

	tr, err := i.Store.GetTimedResponses(s, limit, reverse)
	if err != nil {
		return nil, fmt.Errorf("get timed responses: %w", err)
	}

	return tr, nil
}

func (i Interactor) GetTimedResponseRange(user *models.User, id string, from, to time.Time) ([]*models.TimedResponse, error) {
	s, err := i.GetServiceByID(user, id)
	if err != nil {
		return nil, fmt.Errorf("get service by id: %w", err)
	}

	tr, err := i.Store.GetTimedResponseRange(s, from, to)
	if err != nil {
		return nil, fmt.Errorf("get timed responses range: %w", err)
	}

	return tr, nil
}
