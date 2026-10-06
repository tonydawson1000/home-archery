package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tonydawson1000/home-archery/api/internal/application"
)

type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) http.Handler {
	h := &Handler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/healthz", h.healthz)
	mux.HandleFunc("GET /api/v1/archers", h.listArchers)
	mux.HandleFunc("GET /api/v1/archers/{archerId}/sessions", h.listSessions)
	mux.HandleFunc("POST /api/v1/archers/{archerId}/sessions", h.startSession)
	mux.HandleFunc("GET /api/v1/archers/{archerId}/personal-best", h.personalBest)
	mux.HandleFunc("GET /api/v1/sessions/{sessionId}", h.getSession)
	mux.HandleFunc("POST /api/v1/sessions/{sessionId}/arrows", h.recordArrow)
	mux.HandleFunc("POST /api/v1/sessions/{sessionId}/complete", h.completeSession)
	return mux
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listArchers(w http.ResponseWriter, r *http.Request) {
	archers, err := h.svc.ListArchers(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]archerJSON, 0, len(archers))
	for _, archer := range archers {
		out = append(out, archerJSON{ID: archer.ID, DisplayName: archer.DisplayName})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) listSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.svc.ListSessions(r.Context(), r.PathValue("archerId"))
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]sessionListItemJSON, 0, len(sessions))
	for _, session := range sessions {
		out = append(out, sessionListJSON(session))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) startSession(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.StartSession(r.Context(), r.PathValue("archerId"))
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, sessionJSONFrom(session))
}

func (h *Handler) personalBest(w http.ResponseWriter, r *http.Request) {
	archerID := r.PathValue("archerId")
	raw := r.URL.Query().Get("arrowCount")
	if raw == "" {
		bests, err := h.svc.GetPersonalBests(r.Context(), archerID)
		if err != nil {
			writeAppError(w, err)
			return
		}
		out := make([]personalBestJSON, 0, len(bests))
		for _, best := range bests {
			out = append(out, personalBestJSONFrom(best))
		}
		writeJSON(w, http.StatusOK, out)
		return
	}
	arrowCount, err := strconv.Atoi(raw)
	if err != nil || arrowCount < 1 {
		writeJSON(w, http.StatusBadRequest, errorJSON{Code: "invalid_arrow_count", Message: "arrowCount must be an integer of at least 1"})
		return
	}
	best, err := h.svc.GetPersonalBest(r.Context(), archerID, arrowCount)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, personalBestJSONFrom(best))
}

func (h *Handler) getSession(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.GetSession(r.Context(), r.PathValue("sessionId"))
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sessionJSONFrom(session))
}

func (h *Handler) recordArrow(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ScoreCode string `json:"scoreCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ScoreCode == "" {
		writeJSON(w, http.StatusBadRequest, errorJSON{Code: "invalid_request", Message: "scoreCode is required"})
		return
	}
	got, err := h.svc.RecordArrow(r.Context(), r.PathValue("sessionId"), body.ScoreCode)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, recordedArrowJSON{
		Arrow:     arrowJSONFrom(got.Arrow),
		EndNumber: got.EndNumber,
		Summary:   summaryJSON(got.Summary),
	})
}

func (h *Handler) completeSession(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.CompleteSession(r.Context(), r.PathValue("sessionId"))
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sessionJSONFrom(session))
}
