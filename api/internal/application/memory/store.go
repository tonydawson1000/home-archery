package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

// Store is an in-memory fake of the application ports.
type Store struct {
	mu       sync.RWMutex
	archers  map[string]domain.Archer
	sessions map[string]domain.Session
}

func NewStore() *Store {
	archers := make(map[string]domain.Archer)
	for _, archer := range SeedArchers() {
		archers[archer.ID] = archer
	}
	return &Store{
		archers:  archers,
		sessions: make(map[string]domain.Session),
	}
}

func (s *Store) Archers() application.ArcherRepository   { return archerRepo{s} }
func (s *Store) Sessions() application.SessionRepository { return sessionRepo{s} }

type archerRepo struct{ store *Store }

func (r archerRepo) Get(_ context.Context, id string) (domain.Archer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	archer, ok := r.store.archers[id]
	if !ok {
		return domain.Archer{}, application.ErrArcherNotFound
	}
	return archer, nil
}

func (r archerRepo) List(_ context.Context) ([]domain.Archer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	out := make([]domain.Archer, 0, len(r.store.archers))
	for _, archer := range r.store.archers {
		out = append(out, archer)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DisplayName < out[j].DisplayName })
	return out, nil
}

type sessionRepo struct{ store *Store }

func (r sessionRepo) Save(_ context.Context, session domain.Session) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.sessions[session.ID] = cloneSession(session)
	return nil
}

func (r sessionRepo) Get(_ context.Context, id string) (domain.Session, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	session, ok := r.store.sessions[id]
	if !ok {
		return domain.Session{}, application.ErrSessionNotFound
	}
	return cloneSession(session), nil
}

func (r sessionRepo) ListByArcher(_ context.Context, archerID string) ([]domain.Session, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	out := make([]domain.Session, 0)
	for _, session := range r.store.sessions {
		if session.ArcherID == archerID {
			out = append(out, cloneSession(session))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartedAt.After(out[j].StartedAt) })
	return out, nil
}

func cloneSession(session domain.Session) domain.Session {
	copied := session
	if session.CompletedAt != nil {
		t := *session.CompletedAt
		copied.CompletedAt = &t
	}
	copied.Ends = make([]domain.End, len(session.Ends))
	for i, end := range session.Ends {
		arrows := make([]domain.Arrow, len(end.Arrows))
		copy(arrows, end.Arrows)
		copied.Ends[i] = domain.End{Number: end.Number, Arrows: arrows}
	}
	return copied
}
