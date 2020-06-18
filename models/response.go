package models

import (
	"time"

	"github.com/rs/xid"
)

type TimedResponse struct {
	ID        xid.ID
	At        time.Time
	Server    time.Duration
	Total     time.Duration
	Status    int
	ServiceID string
}
