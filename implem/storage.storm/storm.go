package storage

import (
	"context"
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/rs/zerolog"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/fx"

	"github.com/Depado/vuemonit/cmd"
	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

const bucketName = "History"

type StormStorage struct {
	db  *storm.DB
	log *zerolog.Logger
}

func NewStormStorage(lc fx.Lifecycle, conf *cmd.Conf, log *zerolog.Logger) (interactor.StorageProvider, error) {
	db, err := storm.Open(conf.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to init db: %w", err)
	}
	log.Debug().Str("path", conf.Database.Path).Msg("database opened")

	// User bucket initialization
	if err = db.Init(&models.User{}); err != nil {
		return nil, fmt.Errorf("unable to init user bucket: %w", err)
	}
	log.Debug().Str("bucket", "user").Msg("bucket initialized")

	// Service bucket initialization
	if err = db.Init(&models.Service{}); err != nil {
		return nil, fmt.Errorf("unable to init service bucket: %w", err)
	}
	log.Debug().Str("bucket", "service").Msg("bucket initialized")

	// TimedResponse bucket
	err = db.Bolt.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("create timed response bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create bucket: %w", err)
	}
	log.Debug().Str("bucket", "history").Msg("bucket initialized")

	st := &StormStorage{db: db, log: log}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			st.log.Debug().Str("path", conf.Database.Path).Msg("closing database")
			return st.db.Close()
		},
	})

	return st, nil
}
