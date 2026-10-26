package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/service"
)

type TerminologyHandler interface {
	Insert(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindByID(http.ResponseWriter, *http.Request)
}

type terminologyHandler struct {
	svc service.TerminologyService
}

func NewTerminologyHandler(svc service.TerminologyService) TerminologyHandler {
	return &terminologyHandler{svc: svc}
}

func (h *terminologyHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTerminologyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.svc.Insert(&req)
	if err != nil {
		if isValidationErr(err) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *terminologyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid path value", http.StatusBadRequest)
		return
	}

	var req models.UpdateTerminologyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.svc.Update(id, &req); err != nil {
		var status int
		switch {
		case strings.HasPrefix(err.Error(), "not_found:"):
			status = http.StatusNotFound
		case strings.HasPrefix(err.Error(), "forbidden:"):
			status = http.StatusForbidden
		case isValidationErr(err):
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}
		http.Error(w, err.Error(), status)
		return
	}

	if err = json.NewEncoder(w).Encode(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *terminologyHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	page, err := parsePagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	f := &models.TerminologyFilter{Term: strings.TrimSpace(query.Get("term")), Pagination: page}

	t, total, err := h.svc.FindAll(f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeJson(w, http.StatusOK, models.NewPaginated(t, page, total)); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *terminologyHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid path value", http.StatusBadRequest)
		return
	}
	t, err := h.svc.FindByID(id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err = writeJson(w, http.StatusOK, t); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
