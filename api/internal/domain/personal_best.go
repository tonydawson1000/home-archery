package domain

import "time"

type PersonalBest struct {
	ArcherID   string
	SessionID  string
	Total      int
	ArrowCount int
	AchievedAt time.Time
}

func PersonalBests(sessions []Session) []PersonalBest {
	bestByCount := map[int]PersonalBest{}
	order := make([]int, 0)

	for _, session := range sessions {
		best, ok := bestFromCompleted(session)
		if !ok {
			continue
		}
		existing, seen := bestByCount[best.ArrowCount]
		if !seen {
			bestByCount[best.ArrowCount] = best
			order = append(order, best.ArrowCount)
			continue
		}
		if better(best, existing) {
			bestByCount[best.ArrowCount] = best
		}
	}

	out := make([]PersonalBest, 0, len(order))
	for _, count := range order {
		out = append(out, bestByCount[count])
	}
	return out
}

func PersonalBestForArrowCount(sessions []Session, arrowCount int) (PersonalBest, bool) {
	for _, best := range PersonalBests(sessions) {
		if best.ArrowCount == arrowCount {
			return best, true
		}
	}
	return PersonalBest{}, false
}

func bestFromCompleted(session Session) (PersonalBest, bool) {
	if session.Status != StatusCompleted || session.CompletedAt == nil {
		return PersonalBest{}, false
	}
	sum := session.Summary()
	if sum.ArrowCount == 0 {
		return PersonalBest{}, false
	}
	return PersonalBest{
		ArcherID:   session.ArcherID,
		SessionID:  session.ID,
		Total:      sum.Total,
		ArrowCount: sum.ArrowCount,
		AchievedAt: *session.CompletedAt,
	}, true
}

func better(candidate, existing PersonalBest) bool {
	if candidate.Total > existing.Total {
		return true
	}
	if candidate.Total < existing.Total {
		return false
	}
	return candidate.AchievedAt.After(existing.AchievedAt)
}
