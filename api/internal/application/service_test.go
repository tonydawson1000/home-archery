package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/application/memory"
	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

func testService() (*application.Service, *memory.Store) {
	store := memory.NewStore()
	return &application.Service{
		Archers:  store.Archers(),
		Sessions: store.Sessions(),
		Clock:    application.FixedClock{T: time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)},
		IDs:      application.NewSeqIDs("id"),
	}, store
}

func TestStartSessionUnknownArcher(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	_, err := svc.StartSession(context.Background(), "missing")
	if err != application.ErrArcherNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestStartSessionCreatesInProgressSession(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	session, err := svc.StartSession(context.Background(), memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if session.ID != "id-1" || session.ArcherID != memory.TonyID {
		t.Fatalf("session = %+v", session)
	}
	if session.Status != domain.StatusInProgress || session.ArrowsPerEnd != 6 {
		t.Fatalf("session = %+v", session)
	}
	if len(session.Ends) != 0 {
		t.Fatal("expected empty ends")
	}
	if !session.StartedAt.Equal(time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("startedAt = %s", session.StartedAt)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	_, err := svc.GetSession(context.Background(), "missing")
	if err != application.ErrSessionNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestGetSessionReturnsScorecard(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	ctx := context.Background()
	started, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordArrow(ctx, started.ID, "X"); err != nil {
		t.Fatal(err)
	}
	got, err := svc.GetSession(ctx, started.ID)
	if err != nil {
		t.Fatal(err)
	}
	sum := got.Summary()
	if sum.Total != 10 || sum.ArrowCount != 1 {
		t.Fatalf("summary = %+v", sum)
	}
}

func TestRecordArrowUnknownSession(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	_, err := svc.RecordArrow(context.Background(), "missing", "9")
	if err != application.ErrSessionNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestRecordArrowAppendsAndReturnsRunningSummary(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	ctx := context.Background()
	started, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	got, err := svc.RecordArrow(ctx, started.ID, "9")
	if err != nil {
		t.Fatal(err)
	}
	if got.EndNumber != 1 || got.Arrow.Score.Code != "9" || got.Summary.Total != 9 {
		t.Fatalf("recorded = %+v", got)
	}
}

func TestRecordArrowMapsDomainErrors(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	ctx := context.Background()
	started, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.RecordArrow(ctx, started.ID, "11")
	if !errors.Is(err, domain.ErrInvalidScoreCode) {
		t.Fatalf("err = %v", err)
	}
	if _, err := svc.RecordArrow(ctx, started.ID, "9"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteSession(ctx, started.ID); err != nil {
		t.Fatal(err)
	}
	_, err = svc.RecordArrow(ctx, started.ID, "9")
	if !errors.Is(err, domain.ErrSessionCompleted) {
		t.Fatalf("err = %v", err)
	}
}

func TestCompleteSessionEmptyAndUnknown(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	ctx := context.Background()
	_, err := svc.CompleteSession(ctx, "missing")
	if err != application.ErrSessionNotFound {
		t.Fatalf("err = %v", err)
	}
	started, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.CompleteSession(ctx, started.ID)
	if !errors.Is(err, domain.ErrEmptyComplete) {
		t.Fatalf("err = %v", err)
	}
}

func TestCompleteSessionSetsCompletedAt(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	ctx := context.Background()
	started, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordArrow(ctx, started.ID, "10"); err != nil {
		t.Fatal(err)
	}
	done, err := svc.CompleteSession(ctx, started.ID)
	if err != nil {
		t.Fatal(err)
	}
	if done.Status != domain.StatusCompleted || done.CompletedAt == nil {
		t.Fatalf("done = %+v", done)
	}
}

func TestListSessionsUnknownArcher(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	_, err := svc.ListSessions(context.Background(), "missing")
	if err != application.ErrArcherNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestListSessionsNewestFirst(t *testing.T) {
	t.Parallel()
	store := memory.NewStore()
	clock := &stepClock{times: []time.Time{
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
	}}
	svc := &application.Service{
		Archers:  store.Archers(),
		Sessions: store.Sessions(),
		Clock:    clock,
		IDs:      application.NewSeqIDs("id"),
	}
	ctx := context.Background()
	first, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	list, err := svc.ListSessions(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != second.ID || list[1].ID != first.ID {
		t.Fatalf("list = %+v", list)
	}
}

func TestGetPersonalBestUnknownArcher(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	_, err := svc.GetPersonalBests(context.Background(), "missing")
	if err != application.ErrArcherNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestGetPersonalBestFiltersIncompleteAndOptionalArrowCount(t *testing.T) {
	t.Parallel()
	svc, _ := testService()
	ctx := context.Background()
	open, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordArrow(ctx, open.ID, "X"); err != nil {
		t.Fatal(err)
	}

	done, err := svc.StartSession(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordArrow(ctx, done.ID, "9"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteSession(ctx, done.ID); err != nil {
		t.Fatal(err)
	}

	bests, err := svc.GetPersonalBests(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if len(bests) != 1 || bests[0].SessionID != done.ID || bests[0].Total != 9 {
		t.Fatalf("bests = %+v", bests)
	}

	got, err := svc.GetPersonalBest(ctx, memory.TonyID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.SessionID != done.ID {
		t.Fatalf("single = %+v", got)
	}
	_, err = svc.GetPersonalBest(ctx, memory.TonyID, 60)
	if err != application.ErrPersonalBestNotFound {
		t.Fatalf("missing bucket err = %v", err)
	}
}

type stepClock struct {
	i     int
	times []time.Time
}

func (c *stepClock) Now() time.Time {
	t := c.times[c.i]
	if c.i < len(c.times)-1 {
		c.i++
	}
	return t
}
