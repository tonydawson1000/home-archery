package application

import (
	"context"
	"time"

	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

type ArcherRepository interface {
	Get(ctx context.Context, id string) (domain.Archer, error)
	List(ctx context.Context) ([]domain.Archer, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session domain.Session) error
	Get(ctx context.Context, id string) (domain.Session, error)
	ListByArcher(ctx context.Context, archerID string) ([]domain.Session, error)
}

type Clock interface {
	Now() time.Time
}

type IDs interface {
	New() string
}

type Service struct {
	Archers  ArcherRepository
	Sessions SessionRepository
	Clock    Clock
	IDs      IDs
}
