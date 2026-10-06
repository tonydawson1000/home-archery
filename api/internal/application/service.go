package application

import (
	"context"

	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

func (s *Service) ListArchers(ctx context.Context) ([]domain.Archer, error) {
	return s.Archers.List(ctx)
}

func (s *Service) StartSession(ctx context.Context, archerID string) (domain.Session, error) {
	if _, err := s.Archers.Get(ctx, archerID); err != nil {
		return domain.Session{}, err
	}
	session := domain.NewSession(s.IDs.New(), archerID, s.Clock.Now())
	if err := s.Sessions.Save(ctx, session); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (s *Service) GetSession(ctx context.Context, sessionID string) (domain.Session, error) {
	return s.Sessions.Get(ctx, sessionID)
}

type RecordedArrow struct {
	Arrow     domain.Arrow
	EndNumber int
	Summary   domain.ScoreSummary
}

func (s *Service) RecordArrow(ctx context.Context, sessionID, scoreCode string) (RecordedArrow, error) {
	session, err := s.Sessions.Get(ctx, sessionID)
	if err != nil {
		return RecordedArrow{}, err
	}
	arrow, endNumber, err := session.RecordArrow(scoreCode, s.IDs.New(), s.Clock.Now())
	if err != nil {
		return RecordedArrow{}, err
	}
	if err := s.Sessions.Save(ctx, session); err != nil {
		return RecordedArrow{}, err
	}
	return RecordedArrow{Arrow: arrow, EndNumber: endNumber, Summary: session.Summary()}, nil
}

func (s *Service) CompleteSession(ctx context.Context, sessionID string) (domain.Session, error) {
	session, err := s.Sessions.Get(ctx, sessionID)
	if err != nil {
		return domain.Session{}, err
	}
	if err := session.Complete(s.Clock.Now()); err != nil {
		return domain.Session{}, err
	}
	if err := s.Sessions.Save(ctx, session); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (s *Service) ListSessions(ctx context.Context, archerID string) ([]domain.Session, error) {
	if _, err := s.Archers.Get(ctx, archerID); err != nil {
		return nil, err
	}
	return s.Sessions.ListByArcher(ctx, archerID)
}

func (s *Service) GetPersonalBests(ctx context.Context, archerID string) ([]domain.PersonalBest, error) {
	sessions, err := s.ListSessions(ctx, archerID)
	if err != nil {
		return nil, err
	}
	return domain.PersonalBests(sessions), nil
}

func (s *Service) GetPersonalBest(ctx context.Context, archerID string, arrowCount int) (domain.PersonalBest, error) {
	sessions, err := s.ListSessions(ctx, archerID)
	if err != nil {
		return domain.PersonalBest{}, err
	}
	best, ok := domain.PersonalBestForArrowCount(sessions, arrowCount)
	if !ok {
		return domain.PersonalBest{}, ErrPersonalBestNotFound
	}
	return best, nil
}
