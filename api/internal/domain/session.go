package domain

import "time"

const ArrowsPerEnd = 6

type SessionStatus string

const (
	StatusInProgress SessionStatus = "in_progress"
	StatusCompleted  SessionStatus = "completed"
)

type ScoreSummary struct {
	Total      int
	Hits       int
	Golds      int
	XCount     int
	ArrowCount int
}

type Arrow struct {
	ID          string
	ArrowNumber int
	Score       ArrowScore
	RecordedAt  time.Time
}

type End struct {
	Number int
	Arrows []Arrow
}

type Session struct {
	ID           string
	ArcherID     string
	Status       SessionStatus
	ArrowsPerEnd int
	StartedAt    time.Time
	CompletedAt  *time.Time
	Ends         []End
}

func NewSession(id, archerID string, startedAt time.Time) Session {
	return Session{
		ID:           id,
		ArcherID:     archerID,
		Status:       StatusInProgress,
		ArrowsPerEnd: ArrowsPerEnd,
		StartedAt:    startedAt,
	}
}

func (s *Session) RecordArrow(code, arrowID string, at time.Time) (Arrow, int, error) {
	if s.Status == StatusCompleted {
		return Arrow{}, 0, ErrSessionCompleted
	}
	score, err := ParseArrowScore(code)
	if err != nil {
		return Arrow{}, 0, err
	}

	if len(s.Ends) == 0 || len(s.Ends[len(s.Ends)-1].Arrows) >= s.ArrowsPerEnd {
		s.Ends = append(s.Ends, End{Number: len(s.Ends) + 1})
	}
	end := &s.Ends[len(s.Ends)-1]
	arrow := Arrow{
		ID:          arrowID,
		ArrowNumber: len(end.Arrows) + 1,
		Score:       score,
		RecordedAt:  at,
	}
	end.Arrows = append(end.Arrows, arrow)
	return arrow, end.Number, nil
}

func (s *Session) Complete(at time.Time) error {
	if s.Status == StatusCompleted {
		return ErrSessionCompleted
	}
	if s.Summary().ArrowCount == 0 {
		return ErrEmptyComplete
	}
	s.Status = StatusCompleted
	completed := at
	s.CompletedAt = &completed
	return nil
}

func (s Session) Summary() ScoreSummary {
	var sum ScoreSummary
	for _, end := range s.Ends {
		for _, arrow := range end.Arrows {
			sum.ArrowCount++
			sum.Total += arrow.Score.NumericValue
			if arrow.Score.IsHit() {
				sum.Hits++
			}
			if arrow.Score.IsGold() {
				sum.Golds++
			}
			if arrow.Score.Code == "X" {
				sum.XCount++
			}
		}
	}
	return sum
}
