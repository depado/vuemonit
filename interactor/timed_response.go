package interactor

import (
	"fmt"

	"github.com/Depado/vuemonit/models"
)

func (i Interactor) GetTimedResponsesByServiceID(user *models.User, id string) ([]*models.TimedResponse, error) {
	s, err := i.GetServiceByID(user, id)
	if err != nil {
		return nil, fmt.Errorf("get service by id: %w", err)
	}

	tr, err := i.Store.GetTimedResponses(s)
	if err != nil {
		return nil, fmt.Errorf("get timed responses: %w", err)
	}

	return tr, nil
}
