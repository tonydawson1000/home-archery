package domain

import (
	"testing"
	"time"
)

func TestNewSessionStartsInProgress(t *testing.T) {
	t.Parallel()

	started := time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC)
	s := NewSession("sess-1", "archer-1", started)

	if s.ID != "sess-1" || s.ArcherID != "archer-1" {
		t.Fatalf("ids = %s/%s", s.ID, s.ArcherID)
	}
	if s.Status != StatusInProgress {
		t.Fatalf("status = %s, want in_progress", s.Status)
	}
	if s.ArrowsPerEnd != 6 {
		t.Fatalf("arrowsPerEnd = %d", s.ArrowsPerEnd)
	}
	if !s.StartedAt.Equal(started) {
		t.Fatalf("startedAt = %s", s.StartedAt)
	}
	if s.CompletedAt != nil {
		t.Fatal("completedAt should be nil")
	}
	if len(s.Ends) != 0 {
		t.Fatalf("ends = %d, want 0", len(s.Ends))
	}
}

func TestRecordArrowRejectsInvalidCode(t *testing.T) {
	t.Parallel()

	s := NewSession("s", "a", time.Now().UTC())
	_, _, err := s.RecordArrow("11", "arrow-1", time.Now().UTC())
	if err != ErrInvalidScoreCode {
		t.Fatalf("err = %v, want ErrInvalidScoreCode", err)
	}
}

func TestRecordArrowOpensFirstEndAndSummarises(t *testing.T) {
	t.Parallel()

	s := NewSession("s", "a", time.Now().UTC())
	at := time.Date(2026, 10, 6, 10, 1, 0, 0, time.UTC)
	arrow, endNumber, err := s.RecordArrow("X", "arrow-1", at)
	if err != nil {
		t.Fatal(err)
	}
	if endNumber != 1 {
		t.Fatalf("endNumber = %d", endNumber)
	}
	if arrow.ID != "arrow-1" || arrow.ArrowNumber != 1 || arrow.Score.Code != "X" {
		t.Fatalf("arrow = %+v", arrow)
	}
	if !arrow.RecordedAt.Equal(at) {
		t.Fatalf("recordedAt = %s", arrow.RecordedAt)
	}

	sum := s.Summary()
	if sum != (ScoreSummary{Total: 10, Hits: 1, Golds: 1, XCount: 1, ArrowCount: 1}) {
		t.Fatalf("summary = %+v", sum)
	}
}

func TestRecordArrowRollsOverAfterSix(t *testing.T) {
	t.Parallel()

	s := NewSession("s", "a", time.Now().UTC())
	now := time.Now().UTC()
	for i := 0; i < 6; i++ {
		_, endNumber, err := s.RecordArrow("9", idFor(i), now)
		if err != nil {
			t.Fatal(err)
		}
		if endNumber != 1 {
			t.Fatalf("arrow %d end = %d, want 1", i+1, endNumber)
		}
	}

	arrow, endNumber, err := s.RecordArrow("M", "arrow-7", now)
	if err != nil {
		t.Fatal(err)
	}
	if endNumber != 2 || arrow.ArrowNumber != 1 {
		t.Fatalf("seventh arrow end=%d number=%d", endNumber, arrow.ArrowNumber)
	}
	if len(s.Ends) != 2 {
		t.Fatalf("ends = %d", len(s.Ends))
	}
}

func TestSummaryCountsHitsGoldsAndMisses(t *testing.T) {
	t.Parallel()

	s := NewSession("s", "a", time.Now().UTC())
	now := time.Now().UTC()
	codes := []string{"X", "10", "9", "8", "M", "1"}
	for i, code := range codes {
		if _, _, err := s.RecordArrow(code, idFor(i), now); err != nil {
			t.Fatal(err)
		}
	}

	sum := s.Summary()
	want := ScoreSummary{Total: 10 + 10 + 9 + 8 + 0 + 1, Hits: 5, Golds: 3, XCount: 1, ArrowCount: 6}
	if sum != want {
		t.Fatalf("summary = %+v, want %+v", sum, want)
	}
}

func TestCompleteRequiresAtLeastOneArrow(t *testing.T) {
	t.Parallel()

	s := NewSession("s", "a", time.Now().UTC())
	if err := s.Complete(time.Now().UTC()); err != ErrEmptyComplete {
		t.Fatalf("err = %v, want ErrEmptyComplete", err)
	}
}

func TestCompleteStopsFurtherArrows(t *testing.T) {
	t.Parallel()

	s := NewSession("s", "a", time.Now().UTC())
	if _, _, err := s.RecordArrow("9", "a1", time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	done := time.Date(2026, 10, 6, 11, 0, 0, 0, time.UTC)
	if err := s.Complete(done); err != nil {
		t.Fatal(err)
	}
	if s.Status != StatusCompleted {
		t.Fatalf("status = %s", s.Status)
	}
	if s.CompletedAt == nil || !s.CompletedAt.Equal(done) {
		t.Fatalf("completedAt = %v", s.CompletedAt)
	}
	if _, _, err := s.RecordArrow("9", "a2", time.Now().UTC()); err != ErrSessionCompleted {
		t.Fatalf("err = %v, want ErrSessionCompleted", err)
	}
	if err := s.Complete(time.Now().UTC()); err != ErrSessionCompleted {
		t.Fatalf("second complete err = %v", err)
	}
}

func idFor(i int) string {
	return string(rune('a' + i))
}
