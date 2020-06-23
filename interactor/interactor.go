package interactor

import (
	"net/http"

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

// LogicHandler is the interface that must implement the struct holding the
// usecases. The interactor struct must implement this interface.
// This interface is the main entrypoint for the router.
type LogicHandler interface {
	Register(email, password string) error
	Login(email, password string) (*models.TokenPair, *http.Cookie, error)
	Refresh(token string) (*models.TokenPair, error)
	Logout() *http.Cookie
	AuthCheck(w http.ResponseWriter, r *http.Request) (*models.User, error)
	NewService(user *models.User, name, description, url string) (*models.Service, error)

	FormatSelf(user *models.User) interface{}
	FormatService(svc *models.Service) interface{}
	FormatServices(svx []*models.Service) interface{}
	FormatTimedResponses(tr []*models.TimedResponse) interface{}

	GetServiceByID(user *models.User, id string) (*models.Service, error)
	GetServices(user *models.User) ([]*models.Service, error)
	GetTimedResponsesByServiceID(user *models.User, id string) ([]*models.TimedResponse, error)
}

// AuthProvider is a simple auth provider interface
type AuthProvider interface {
	GenerateTokenPair(user *models.User) (*models.TokenPair, error)
	CheckToken(token string) (*jwt.StandardClaims, error)
	ValidateBearerToken(r *http.Request) (*jwt.StandardClaims, error)

	ValidateCookie(r *http.Request) (*jwt.StandardClaims, bool, error)
	GenerateCookie(u *models.User, tp *models.TokenPair) (*http.Cookie, error)
	DropAccessCookie() *http.Cookie
}

// Formatter is a simple interface in charge of formatting our documents
type Formatter interface {
	Self(user *models.User) interface{}
	Service(svc *models.Service) interface{}
	Services(svx []*models.Service) interface{}
	TimedResponses(tr []*models.TimedResponse) interface{}
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
	GetServices(user *models.User) ([]*models.Service, error)
}

type Scheduler interface {
	Start(svc *models.Service) error
	Restart(svc *models.Service) error
}
