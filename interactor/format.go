package interactor

import (
	"github.com/Depado/vuemonit/models"
)

func (i Interactor) FormatSelf(user *models.User) interface{} {
	return i.Formatter.Self(user)
}

func (i Interactor) FormatService(svc *models.Service) interface{} {
	return i.Formatter.Service(svc)
}

func (i Interactor) FormatServices(svx []*models.Service) interface{} {
	return i.Formatter.Services(svx)
}

func (i Interactor) FormatTimedResponses(tr []*models.TimedResponse) interface{} {
	return i.Formatter.TimedResponses(tr)
}
