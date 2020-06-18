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

func (g gocronScheduler) fetchAndSave(s *models.Service) {
	tr, err := s.Fetch()
	if err != nil {
		g.log.Err(err).Str("id", s.ID).Msg("unable to fetch")
		return
	}
	if err := g.sp.SaveTimedResponse(tr); err != nil {
		g.log.Err(err).Str("id", s.ID).Msg("unable to store timed response")
		return
	}
	if err := g.sp.SaveRawService(s); err != nil {
		g.log.Err(err).Str("id", s.ID).Msg("unable to store timed response")
		return
	}
}

func NewGocronScheduler(sp interactor.StorageProvider, log *zerolog.Logger) (interactor.Scheduler, error) {
	svx, err := sp.GetAllServices()
	if err != nil {
		return nil, fmt.Errorf("unable to start scheduler: %w", err)
	}
	gs := &gocronScheduler{
		sch: gocron.NewScheduler(time.Local),
		sp:  sp,
		log: log,
	}

	for _, s := range svx {
		_, err := gs.sch.Every(5).Minutes().SetTag([]string{s.ID}).Do(gs.fetchAndSave, s)
		if err != nil {
			gs.log.Err(err).Str("id", s.ID).Msg("unable to start routine")
		} else {
			gs.log.Debug().Str("id", s.ID).Msg("started routine")
		}
	}
	gs.sch.StartAsync()

	return gs, nil
}

func (gs gocronScheduler) Restart(s *models.Service) error {
	panic("not implemented")
}

func (gs gocronScheduler) Start(s *models.Service) error {
	gs.fetchAndSave(s)
	_, err := gs.sch.Every(5).Minutes().SetTag([]string{s.ID}).Do(gs.fetchAndSave, s)
	if err != nil {
		return fmt.Errorf("unable to start: %w", err)
	} else {
		gs.log.Debug().Str("id", s.ID).Msg("started routine")
	}
	return nil
}
