package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/application/memory"
)

func testHandler() http.Handler {
	store := memory.NewStore()
	svc := &application.Service{
		Archers:  store.Archers(),
		Sessions: store.Sessions(),
		Clock:    application.FixedClock{T: time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)},
		IDs:      application.NewSeqIDs("id"),
	}
	return NewHandler(svc)
}

func TestHealthz(t *testing.T) {
	t.Parallel()
	rec := do(t, testHandler(), http.MethodGet, "/api/v1/healthz", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if got := jsonBody(t, rec)["status"]; got != "ok" {
		t.Fatalf("body = %v", rec.Body.String())
	}
}

func TestListArchers(t *testing.T) {
	t.Parallel()
	rec := do(t, testHandler(), http.MethodGet, "/api/v1/archers", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
	var archers []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &archers); err != nil {
		t.Fatal(err)
	}
	if len(archers) != 2 || archers[0]["displayName"] != "Becky" {
		t.Fatalf("archers = %v", archers)
	}
}

func TestStartGetRecordCompleteAndList(t *testing.T) {
	t.Parallel()
	h := testHandler()
	path := "/api/v1/archers/" + memory.TonyID + "/sessions"

	started := do(t, h, http.MethodPost, path, nil)
	if started.Code != http.StatusCreated {
		t.Fatalf("start = %d %s", started.Code, started.Body.String())
	}
	session := jsonBody(t, started)
	id, _ := session["id"].(string)
	if session["status"] != "in_progress" || session["arrowsPerEnd"].(float64) != 6 {
		t.Fatalf("session = %v", session)
	}

	empty := do(t, h, http.MethodPost, "/api/v1/sessions/"+id+"/complete", nil)
	assertError(t, empty, http.StatusConflict, "empty_complete")

	missingArcher := do(t, h, http.MethodPost, "/api/v1/archers/00000000-0000-4000-8000-000000000000/sessions", nil)
	assertError(t, missingArcher, http.StatusNotFound, "archer_not_found")

	arrow := do(t, h, http.MethodPost, "/api/v1/sessions/"+id+"/arrows", map[string]string{"scoreCode": "X"})
	if arrow.Code != http.StatusCreated {
		t.Fatalf("arrow = %d %s", arrow.Code, arrow.Body.String())
	}
	arrowBody := jsonBody(t, arrow)
	if arrowBody["endNumber"].(float64) != 1 {
		t.Fatalf("arrow body = %v", arrowBody)
	}

	bad := do(t, h, http.MethodPost, "/api/v1/sessions/"+id+"/arrows", map[string]string{"scoreCode": "11"})
	assertError(t, bad, http.StatusBadRequest, "invalid_score_code")

	got := do(t, h, http.MethodGet, "/api/v1/sessions/"+id, nil)
	if got.Code != http.StatusOK {
		t.Fatalf("get = %d %s", got.Code, got.Body.String())
	}

	missingSession := do(t, testHandler(), http.MethodPost, "/api/v1/sessions/id-1/complete", nil)
	assertError(t, missingSession, http.StatusNotFound, "session_not_found")

	done := do(t, h, http.MethodPost, "/api/v1/sessions/"+id+"/complete", nil)
	if done.Code != http.StatusOK {
		t.Fatalf("complete = %d %s", done.Code, done.Body.String())
	}
	if jsonBody(t, done)["status"] != "completed" {
		t.Fatalf("complete body = %s", done.Body.String())
	}

	conflict := do(t, h, http.MethodPost, "/api/v1/sessions/"+id+"/arrows", map[string]string{"scoreCode": "9"})
	assertError(t, conflict, http.StatusConflict, "session_completed")

	list := do(t, h, http.MethodGet, path, nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list = %d %s", list.Code, list.Body.String())
	}
	var items []map[string]any
	if err := json.Unmarshal(list.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0]["id"] != id {
		t.Fatalf("list = %v", items)
	}

	bests := do(t, h, http.MethodGet, "/api/v1/archers/"+memory.TonyID+"/personal-best", nil)
	if bests.Code != http.StatusOK {
		t.Fatalf("bests = %d %s", bests.Code, bests.Body.String())
	}
	var bestList []map[string]any
	if err := json.Unmarshal(bests.Body.Bytes(), &bestList); err != nil {
		t.Fatal(err)
	}
	if len(bestList) != 1 || bestList[0]["sessionId"] != id {
		t.Fatalf("bests = %v", bestList)
	}

	one := do(t, h, http.MethodGet, "/api/v1/archers/"+memory.TonyID+"/personal-best?arrowCount=1", nil)
	if one.Code != http.StatusOK {
		t.Fatalf("one = %d %s", one.Code, one.Body.String())
	}
	none := do(t, h, http.MethodGet, "/api/v1/archers/"+memory.TonyID+"/personal-best?arrowCount=60", nil)
	assertError(t, none, http.StatusNotFound, "personal_best_not_found")
}

func do(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		rdr = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func jsonBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("json %s: %v", rec.Body.String(), err)
	}
	return out
}

func assertError(t *testing.T, rec *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d want %d body = %s", rec.Code, status, rec.Body.String())
	}
	got := jsonBody(t, rec)
	if got["code"] != code {
		t.Fatalf("code = %v want %s", got["code"], code)
	}
	if _, ok := got["message"].(string); !ok {
		t.Fatalf("missing message: %v", got)
	}
}
