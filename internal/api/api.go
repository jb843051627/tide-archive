package api

import (
	"encoding/json"
	"github.com/jb843051627/tide-archive/internal/service"
	"net/http"
	"strings"
)

type API struct{ svc *service.Service }

func New(svc *service.Service) *API { return &API{svc: svc} }
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.health)
	mux.HandleFunc("/sessions", a.sessions)
	mux.HandleFunc("/sessions/", a.sessionAction)
	mux.Handle("/", http.FileServer(http.Dir("web")))
	return mux
}
func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (a *API) sessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		ID   string `json:"id"`
		Site string `json:"site"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, err)
		return
	}
	if err := a.svc.CreateSession(r.Context(), body.ID, body.Site); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, body)
}
func (a *API) sessionAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/sessions/"), "/")
	if len(parts) != 2 || r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var err error
	switch parts[1] {
	case "start":
		err = a.svc.StartSession(r.Context(), parts[0])
	case "close":
		err = a.svc.CloseSession(r.Context(), parts[0])
	case "archive":
		err = a.svc.Archive(r.Context(), parts[0])
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"session": parts[0], "action": parts[1]})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}
