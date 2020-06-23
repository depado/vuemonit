package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

type gocronScheduler struct {
	sch *gocron.Scheduler
	sp  interactor.StorageProvider
	log *zerolog.Logger
}

// fetch is a simple method to fetch a TimedResponse for a service and saving it
// in database, as well as updating the service itself with the latest data
func (gs gocronScheduler) fetch(s *models.Service) error {
	tr, err := s.Fetch()
	if err != nil {
		return fmt.Errorf("unable to fetch: %w", err)
	}
	if err := gs.sp.SaveTimedResponse(tr); err != nil {
		return fmt.Errorf("unable to save timed response: %w", err)
	}
	if err := gs.sp.SaveRawService(s); err != nil {
		return fmt.Errorf("unable to save service: %w", err)
	}
	return nil
}

// fetchAsync is a simple wrapper around fetch that will log any error that
// occurs. This is useful when ran using gocron as there is no way of handling
// the errors returned by the tasks
func (gs gocronScheduler) fetchAsync(s *models.Service) {
	clog := gs.log.With().Str("id", s.ID).Str("url", s.HealthCheck.URL).Logger()
	if err := gs.fetch(s); err != nil {
		clog.Err(err).Msg("unable to fetch")
	} else {
		clog.Debug().Msg("fetched")
	}
}

// NewGocronScheduler implements the interactor.Scheduler interface so it can be
// used directly in the interactor. It is using the gocron library to create
// periodic tasks
func NewGocronScheduler(sp interactor.StorageProvider, log *zerolog.Logger) (interactor.Scheduler, error) {
	gs := &gocronScheduler{
		sch: gocron.NewScheduler(time.Local),
		sp:  sp,
		log: log,
	}

	svx, err := sp.GetAllServices()
	if err != nil {
		return nil, fmt.Errorf("unable to start scheduler: %w", err)
	}

	for _, s := range svx {
		_, err := gs.sch.Every(10).Seconds().SetTag([]string{s.ID}).Do(gs.fetchAsync, s)
		if err != nil {
			gs.log.Err(err).Str("id", s.ID).Msg("unable to start routine")
		} else {
			gs.log.Debug().Str("id", s.ID).Msg("started routine")
		}
	}

	gs.sch.StartAsync()
	gs.log.Info().Int("services", len(svx)).Msg("started background routine")

	return gs, nil
}

// Restart is used to stop and restart a periodic task for a specific service
// that changed during run time, for example if its duration or URL changes
func (gs gocronScheduler) Restart(s *models.Service) error {
	panic("not implemented")
}

// Start is used to start the periodic task associated to a service
// It will first execute a fetch and return an error as quickly as possible as
// to handle the case of misconfigured URLs for example
func (gs gocronScheduler) Start(s *models.Service) error {
	if err := gs.fetch(s); err != nil {
		return fmt.Errorf("unable to fetch: %w", err)
	}
	_, err := gs.sch.Every(1).Minute().SetTag([]string{s.ID}).Do(gs.fetchAsync, s)
	if err != nil {
		return fmt.Errorf("unable to start: %w", err)
	} else {
		gs.log.Debug().Str("id", s.ID).Msg("started routine")
	}
	return nil
}
