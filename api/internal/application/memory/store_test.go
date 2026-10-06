package memory

import (
	"context"
	"testing"
	"time"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

func TestStoreSeedsTonyAndBecky(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := NewStore()

	archers, err := store.Archers().List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(archers) != 2 {
		t.Fatalf("archers = %d", len(archers))
	}
	if archers[0].DisplayName != "Becky" || archers[0].ID != BeckyID {
		t.Fatalf("first archer = %+v", archers[0])
	}
	if archers[1].DisplayName != "Tony" || archers[1].ID != TonyID {
		t.Fatalf("second archer = %+v", archers[1])
	}

	_, err = store.Archers().Get(ctx, "missing")
	if err != application.ErrArcherNotFound {
		t.Fatalf("err = %v", err)
	}
}

func TestSessionRepoSaveGetAndListNewestFirst(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := NewStore()
	sessions := store.Sessions()

	older := domain.NewSession("old", TonyID, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	newer := domain.NewSession("new", TonyID, time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC))
	other := domain.NewSession("other", BeckyID, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC))
	if err := sessions.Save(ctx, older); err != nil {
		t.Fatal(err)
	}
	if err := sessions.Save(ctx, newer); err != nil {
		t.Fatal(err)
	}
	if err := sessions.Save(ctx, other); err != nil {
		t.Fatal(err)
	}

	got, err := sessions.Get(ctx, "new")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "new" {
		t.Fatalf("got %+v", got)
	}

	list, err := sessions.ListByArcher(ctx, TonyID)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != "new" || list[1].ID != "old" {
		t.Fatalf("list = %+v", list)
	}

	_, err = sessions.Get(ctx, "missing")
	if err != application.ErrSessionNotFound {
		t.Fatalf("err = %v", err)
	}
}
