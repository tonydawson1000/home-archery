package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/application/memory"
	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

func TestArcherListAndGet(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()

	archers, err := store.Archers().List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(archers) != 2 {
		t.Fatalf("archers = %d", len(archers))
	}
	if archers[0].DisplayName != "Becky" || archers[0].ID != memory.BeckyID {
		t.Fatalf("first = %+v", archers[0])
	}
	if archers[1].DisplayName != "Tony" || archers[1].ID != memory.TonyID {
		t.Fatalf("second = %+v", archers[1])
	}

	got, err := store.Archers().Get(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if got.DisplayName != "Tony" {
		t.Fatalf("got %+v", got)
	}

	_, err = store.Archers().Get(ctx, "00000000-0000-4000-8000-000000000000")
	if err != application.ErrArcherNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestSessionSaveGetAndListNewestFirst(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	sessions := store.Sessions()

	olderID := "11111111-1111-4111-8111-111111111111"
	newerID := "22222222-2222-4222-8222-222222222222"
	otherID := "33333333-3333-4333-8333-333333333333"

	older := domain.NewSession(olderID, memory.TonyID, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	newer := domain.NewSession(newerID, memory.TonyID, time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC))
	other := domain.NewSession(otherID, memory.BeckyID, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC))
	for _, session := range []domain.Session{older, newer, other} {
		if err := sessions.Save(ctx, session); err != nil {
			t.Fatal(err)
		}
	}

	got, err := sessions.Get(ctx, newerID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != newerID || got.Status != domain.StatusInProgress || got.ArrowsPerEnd != 6 {
		t.Fatalf("got %+v", got)
	}
	if len(got.Ends) != 0 {
		t.Fatalf("ends = %d", len(got.Ends))
	}

	list, err := sessions.ListByArcher(ctx, memory.TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != newerID || list[1].ID != olderID {
		t.Fatalf("list = %+v", ids(list))
	}

	_, err = sessions.Get(ctx, "00000000-0000-4000-8000-000000000099")
	if err != application.ErrSessionNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestSessionSaveReloadsArrowsAndProjections(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	sessions := store.Sessions()

	id := "44444444-4444-4444-8444-444444444444"
	session := domain.NewSession(id, memory.TonyID, time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC))
	at := time.Date(2026, 10, 6, 12, 1, 0, 0, time.UTC)
	if _, _, err := session.RecordArrow("X", "55555555-5555-4555-8555-555555555555", at); err != nil {
		t.Fatal(err)
	}
	if _, _, err := session.RecordArrow("M", "66666666-6666-4666-8666-666666666666", at.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := sessions.Save(ctx, session); err != nil {
		t.Fatal(err)
	}

	got, err := sessions.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	sum := got.Summary()
	if sum != (domain.ScoreSummary{Total: 10, Hits: 1, Golds: 1, XCount: 1, ArrowCount: 2}) {
		t.Fatalf("summary = %+v", sum)
	}
	if len(got.Ends) != 1 || len(got.Ends[0].Arrows) != 2 {
		t.Fatalf("ends = %+v", got.Ends)
	}
	if got.Ends[0].Arrows[0].Score.Code != "X" || got.Ends[0].Arrows[1].Score.Code != "M" {
		t.Fatalf("arrows = %+v", got.Ends[0].Arrows)
	}

	if _, _, err := got.RecordArrow("9", "77777777-7777-4777-8777-777777777777", at.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := sessions.Save(ctx, got); err != nil {
		t.Fatal(err)
	}
	reloaded, err := sessions.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Summary().ArrowCount != 3 || reloaded.Summary().Total != 19 {
		t.Fatalf("reloaded summary = %+v", reloaded.Summary())
	}
}

func TestSessionCompleteRoundTrip(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	sessions := store.Sessions()

	id := "88888888-8888-4888-8888-888888888888"
	session := domain.NewSession(id, memory.TonyID, time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC))
	if _, _, err := session.RecordArrow("10", "99999999-9999-4999-8999-999999999999", time.Date(2026, 10, 6, 12, 2, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	done := time.Date(2026, 10, 6, 13, 0, 0, 0, time.UTC)
	if err := session.Complete(done); err != nil {
		t.Fatal(err)
	}
	if err := sessions.Save(ctx, session); err != nil {
		t.Fatal(err)
	}

	got, err := sessions.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.StatusCompleted || got.CompletedAt == nil || !got.CompletedAt.UTC().Equal(done) {
		t.Fatalf("got %+v completedAt=%v", got, got.CompletedAt)
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	store, err := Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	if err := store.truncateScoring(ctx); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return store
}

func ids(sessions []domain.Session) []string {
	out := make([]string, len(sessions))
	for i, session := range sessions {
		out[i] = session.ID
	}
	return out
}
