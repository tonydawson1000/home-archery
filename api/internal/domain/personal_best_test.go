package domain

import (
	"testing"
	"time"
)

func TestPersonalBestPolicyIgnoresIncompleteAndGroupsByArrowCount(t *testing.T) {
	t.Parallel()

	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	incomplete := scoredSession("s-open", "a1", t0, nil, []string{"X", "X"})
	lowSix := completedSession("s-low", "a1", t0, t0.Add(time.Hour), sixNines())
	highSix := completedSession("s-high", "a1", t0.Add(2*time.Hour), t0.Add(3*time.Hour), sixXs())
	sixty := completedSession("s-60", "a1", t0, t0.Add(time.Hour), sixtyArrows())

	gots := PersonalBests([]Session{incomplete, lowSix, highSix, sixty})
	if len(gots) != 2 {
		t.Fatalf("bests = %d, want 2 buckets", len(gots))
	}

	byCount := map[int]PersonalBest{}
	for _, b := range gots {
		byCount[b.ArrowCount] = b
	}

	six := byCount[6]
	if six.SessionID != "s-high" || six.Total != 60 {
		t.Fatalf("6-arrow best = %+v", six)
	}
	if six.ArcherID != "a1" {
		t.Fatalf("archer = %s", six.ArcherID)
	}

	long := byCount[60]
	if long.SessionID != "s-60" || long.ArrowCount != 60 {
		t.Fatalf("60-arrow best = %+v", long)
	}
}

func TestPersonalBestPolicyTieBreaksMostRecent(t *testing.T) {
	t.Parallel()

	t0 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	older := completedSession("older", "a1", t0, t0.Add(time.Hour), sixXs())
	newer := completedSession("newer", "a1", t0.Add(2*time.Hour), t0.Add(3*time.Hour), sixXs())

	gots := PersonalBests([]Session{older, newer})
	if len(gots) != 1 || gots[0].SessionID != "newer" {
		t.Fatalf("tie best = %+v", gots)
	}
	if !gots[0].AchievedAt.Equal(*newer.CompletedAt) {
		t.Fatalf("achievedAt = %s", gots[0].AchievedAt)
	}
}

func TestPersonalBestForArrowCount(t *testing.T) {
	t.Parallel()

	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	six := completedSession("s6", "a1", t0, t0.Add(time.Hour), sixXs())
	got, ok := PersonalBestForArrowCount([]Session{six}, 6)
	if !ok || got.SessionID != "s6" {
		t.Fatalf("got=%+v ok=%v", got, ok)
	}
	if _, ok := PersonalBestForArrowCount([]Session{six}, 60); ok {
		t.Fatal("expected no 60-arrow best")
	}
}

func completedSession(id, archer string, started, completed time.Time, codes []string) Session {
	s := scoredSession(id, archer, started, nil, codes)
	if err := s.Complete(completed); err != nil {
		panic(err)
	}
	return s
}

func scoredSession(id, archer string, started time.Time, _ *time.Time, codes []string) Session {
	s := NewSession(id, archer, started)
	for i, code := range codes {
		if _, _, err := s.RecordArrow(code, idFor(i), started.Add(time.Duration(i)*time.Minute)); err != nil {
			panic(err)
		}
	}
	return s
}

func sixNines() []string {
	return []string{"9", "9", "9", "9", "9", "9"}
}

func sixXs() []string {
	return []string{"X", "X", "X", "X", "X", "X"}
}

func sixtyArrows() []string {
	codes := make([]string, 60)
	for i := range codes {
		codes[i] = "1"
	}
	return codes
}
