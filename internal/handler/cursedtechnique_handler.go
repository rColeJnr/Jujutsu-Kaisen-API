package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/service"
)

type CursedTechniqueHandler interface {
	Insert(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindByID(http.ResponseWriter, *http.Request)
}

type cursedTechniqueHandler struct {
	svc service.CursedTechniqueService
}

func NewCursedTechniqueHandler(svc service.CursedTechniqueService) CursedTechniqueHandler {
	return &cursedTechniqueHandler{svc: svc}
}

func (h *cursedTechniqueHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCursedTechniqueRequest
	if !decodeJSON(w, r, &req) {
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
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *cursedTechniqueHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req models.UpdateCursedTechniqueRequest
	if !decodeJSON(w, r, &req) {
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

func (h *cursedTechniqueHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	page, err := parsePagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	f := &models.CursedTechniqueFilter{Name: strings.TrimSpace(query.Get("name")), Pagination: page}

	ct, total, err := h.svc.FindAll(f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeJson(w, http.StatusOK, models.NewPaginated(ct, page, total)); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *cursedTechniqueHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ct, err := h.svc.FindByID(id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err = writeJson(w, http.StatusOK, ct); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}
