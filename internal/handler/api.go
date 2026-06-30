package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/stefanSmk/privaquest/internal/i18n"
	"github.com/stefanSmk/privaquest/internal/service"
	"github.com/stefanSmk/privaquest/internal/store"
)

type API struct {
	svc *service.Service
}

func New(svc *service.Service) *API {
	return &API{svc: svc}
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "product": "privaquest"})
}

type submitBody struct {
	Type        string `json:"type"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Locale      string `json:"locale"`
}

func (a *API) SubmitRequest(w http.ResponseWriter, r *http.Request) {
	var body submitBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if body.Locale == "" {
		body.Locale = r.Header.Get("Accept-Language")
	}

	result, err := a.svc.Submit(r.Context(), service.SubmitInput{
		Type:        store.RequestType(body.Type),
		Email:       body.Email,
		FullName:    body.FullName,
		Description: body.Description,
		Locale:      body.Locale,
	})
	if err != nil {
		status, msg := mapError(err)
		writeJSON(w, status, map[string]string{"error": msg})
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (a *API) CheckStatus(w http.ResponseWriter, r *http.Request) {
	ref := r.URL.Query().Get("reference")
	email := r.URL.Query().Get("email")
	if ref == "" || email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "reference and email required"})
		return
	}

	req, trail, err := a.svc.GetPublicStatus(r.Context(), ref, email)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"reference":  req.Reference,
		"type":       req.Type,
		"status":     req.Status,
		"status_label": i18n.T(req.Locale, "status_"+string(req.Status)),
		"due_at":     req.DueAt,
		"created_at": req.CreatedAt,
		"resolved_at": req.ResolvedAt,
		"audit_count": len(trail),
	})
}

func (a *API) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	d, err := a.svc.Dashboard(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (a *API) AdminList(w http.ResponseWriter, r *http.Request) {
	status := store.RequestStatus(r.URL.Query().Get("status"))
	list, err := a.svc.List(r.Context(), status)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"requests": list, "count": len(list)})
}

func (a *API) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var body struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	req, err := a.svc.UpdateStatus(r.Context(), id, service.StatusUpdate{
		Status: store.RequestStatus(body.Status),
		Notes:  body.Notes,
		Actor:  "admin",
	})
	if err != nil {
		status, msg := mapError(err)
		writeJSON(w, status, map[string]string{"error": msg})
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func (a *API) AdminAudit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	trail, err := a.svc.AuditTrail(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"audit": trail})
}

func (a *API) Locales(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"locales": i18n.SupportedLocales(),
		"types": map[string]string{
			"access":  i18n.T("en", "type_access"),
			"delete":  i18n.T("en", "type_delete"),
			"rectify": i18n.T("en", "type_rectify"),
			"object":  i18n.T("en", "type_object"),
		},
	})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidEmail):
		return http.StatusBadRequest, "invalid email"
	case errors.Is(err, service.ErrInvalidType):
		return http.StatusBadRequest, "invalid request type"
	case errors.Is(err, service.ErrInvalidLocale):
		return http.StatusBadRequest, "unsupported locale (use en, de, fr)"
	case errors.Is(err, store.ErrNotFound):
		return http.StatusNotFound, "not found"
	default:
		return http.StatusInternalServerError, "internal error"
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
