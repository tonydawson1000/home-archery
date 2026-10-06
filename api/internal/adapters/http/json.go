package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

type scoreSummaryJSON struct {
	Total      int `json:"total"`
	Hits       int `json:"hits"`
	Golds      int `json:"golds"`
	XCount     int `json:"xCount"`
	ArrowCount int `json:"arrowCount"`
}

type arrowJSON struct {
	ID          string    `json:"id"`
	ArrowNumber int       `json:"arrowNumber"`
	ScoreCode   string    `json:"scoreCode"`
	ScoreValue  int       `json:"scoreValue"`
	RecordedAt  time.Time `json:"recordedAt"`
}

type endJSON struct {
	EndNumber int         `json:"endNumber"`
	Arrows    []arrowJSON `json:"arrows"`
}

type sessionListItemJSON struct {
	ID          string           `json:"id"`
	ArcherID    string           `json:"archerId"`
	Status      string           `json:"status"`
	StartedAt   time.Time        `json:"startedAt"`
	CompletedAt *time.Time       `json:"completedAt"`
	Summary     scoreSummaryJSON `json:"summary"`
}

type sessionJSON struct {
	sessionListItemJSON
	ArrowsPerEnd int       `json:"arrowsPerEnd"`
	Ends         []endJSON `json:"ends"`
}

type archerJSON struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type personalBestJSON struct {
	ArcherID   string    `json:"archerId"`
	SessionID  string    `json:"sessionId"`
	Total      int       `json:"total"`
	ArrowCount int       `json:"arrowCount"`
	AchievedAt time.Time `json:"achievedAt"`
}

type recordedArrowJSON struct {
	Arrow     arrowJSON        `json:"arrow"`
	EndNumber int              `json:"endNumber"`
	Summary   scoreSummaryJSON `json:"summary"`
}

type errorJSON struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func summaryJSON(s domain.ScoreSummary) scoreSummaryJSON {
	return scoreSummaryJSON{
		Total:      s.Total,
		Hits:       s.Hits,
		Golds:      s.Golds,
		XCount:     s.XCount,
		ArrowCount: s.ArrowCount,
	}
}

func arrowJSONFrom(a domain.Arrow) arrowJSON {
	return arrowJSON{
		ID:          a.ID,
		ArrowNumber: a.ArrowNumber,
		ScoreCode:   a.Score.Code,
		ScoreValue:  a.Score.NumericValue,
		RecordedAt:  a.RecordedAt,
	}
}

func sessionListJSON(s domain.Session) sessionListItemJSON {
	return sessionListItemJSON{
		ID:          s.ID,
		ArcherID:    s.ArcherID,
		Status:      string(s.Status),
		StartedAt:   s.StartedAt,
		CompletedAt: s.CompletedAt,
		Summary:     summaryJSON(s.Summary()),
	}
}

func sessionJSONFrom(s domain.Session) sessionJSON {
	ends := make([]endJSON, 0, len(s.Ends))
	for _, end := range s.Ends {
		arrows := make([]arrowJSON, 0, len(end.Arrows))
		for _, arrow := range end.Arrows {
			arrows = append(arrows, arrowJSONFrom(arrow))
		}
		ends = append(ends, endJSON{EndNumber: end.Number, Arrows: arrows})
	}
	return sessionJSON{
		sessionListItemJSON: sessionListJSON(s),
		ArrowsPerEnd:        s.ArrowsPerEnd,
		Ends:                ends,
	}
}

func personalBestJSONFrom(b domain.PersonalBest) personalBestJSON {
	return personalBestJSON{
		ArcherID:   b.ArcherID,
		SessionID:  b.SessionID,
		Total:      b.Total,
		ArrowCount: b.ArrowCount,
		AchievedAt: b.AchievedAt,
	}
}

func writeAppError(w http.ResponseWriter, err error) {
	status, code, message := mapError(err)
	writeJSON(w, status, errorJSON{Code: code, Message: message})
}

func mapError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidScoreCode):
		return http.StatusBadRequest, domain.ErrInvalidScoreCode.Error(), "Invalid arrow score code"
	case errors.Is(err, domain.ErrSessionCompleted):
		return http.StatusConflict, domain.ErrSessionCompleted.Error(), "Session is already completed"
	case errors.Is(err, domain.ErrEmptyComplete):
		return http.StatusConflict, domain.ErrEmptyComplete.Error(), "Cannot complete a session with no arrows"
	case errors.Is(err, application.ErrArcherNotFound):
		return http.StatusNotFound, application.ErrArcherNotFound.Error(), "Archer not found"
	case errors.Is(err, application.ErrSessionNotFound):
		return http.StatusNotFound, application.ErrSessionNotFound.Error(), "Session not found"
	case errors.Is(err, application.ErrPersonalBestNotFound):
		return http.StatusNotFound, application.ErrPersonalBestNotFound.Error(), "Personal best not found"
	default:
		return http.StatusInternalServerError, "internal_error", "Internal error"
	}
}
