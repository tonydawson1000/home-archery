package memory

import "github.com/tonydawson1000/home-archery/api/internal/domain"

// Seeded household archers — ids match db/seed/001_archers.sql.
const (
	TonyID  = "8f0c1a2e-0b3d-4a7c-9e11-01a1a1a1a1a1"
	BeckyID = "8f0c1a2e-0b3d-4a7c-9e11-02b2b2b2b2b2"
)

func SeedArchers() []domain.Archer {
	return []domain.Archer{
		{ID: TonyID, DisplayName: "Tony"},
		{ID: BeckyID, DisplayName: "Becky"},
	}
}
