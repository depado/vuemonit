package interactor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Depado/vuemonit/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
)

// Interactor contains the dependencies of the service
type Interactor struct {
	Store     StorageProvider
	Auth      AuthProvider
	Formatter Formatter
	Logger    *zerolog.Logger
	Scheduler Scheduler
}

// NewInteractor will return a new Interactor struct with all the required
// dependecies injected. This interactor implements the LogicHandler interface
// and thus should be provided to the router.
func NewInteractor(s StorageProvider, a AuthProvider, f Formatter, l *zerolog.Logger, sch Scheduler) LogicHandler {
	return &Interactor{
		Store:     s,
		Auth:      a,
		Formatter: f,
		Logger:    l,
		Scheduler: sch,
	}
}

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

// LogicHandler is the interface that must implement the struct holding the
// usecases. The interactor struct must implement this interface.
// This interface is the main entrypoint for the router.
type LogicHandler interface {
	Register(email, password string) error
	Login(email, password string) (string, string, error)
	Refresh(token string) (string, string, error)
	AuthCheck(r *http.Request) (*models.User, error)
	NewService(user *models.User, name, description, url string) (*models.Service, error)
	FormatSelf(user *models.User) interface{}
	FormatService(svc *models.Service) (interface{}, error)
	GetServiceByID(user *models.User, id string) (interface{}, error)
}

// AuthProvider is a simple auth provider interface
type AuthProvider interface {
	GenerateTokenPair(user *models.User) (string, string, error)
	Check(token string) (jwt.StandardClaims, error)
	Extract(r *http.Request) (string, error)
}

// Formatter is a simple interface in charge of formatting our documents
type Formatter interface {
	Self(user *models.User) interface{}
	Service(svc *models.Service, trcount int) interface{}
}

// StorageProvider is a storage interface
type StorageProvider interface {
	SaveService(user *models.User, svc *models.Service) error
	SaveRawService(svc *models.Service) error
	SaveUser(usr *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetTimedResponses(svc *models.Service) ([]*models.TimedResponse, error)
	CountTimedResponses(svc *models.Service) (int, error)
	SaveTimedResponse(tr *models.TimedResponse) error
	GetAllServices() ([]*models.Service, error)
	GetServiceByID(id string) (*models.Service, error)
}

type Scheduler interface {
	Start(svc *models.Service) error
	Restart(svc *models.Service) error
}
