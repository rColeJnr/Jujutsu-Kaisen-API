package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/service"
)

type CharacterHandler interface {
	Insert(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindByID(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

// CharacterHandler handles character CRUD HTTP requests.
type characterHandler struct {
	svc service.CharacterService
}

// NewCharacterHandler creates a new CharacterHandler.
func NewCharacterHandler(svc service.CharacterService) CharacterHandler {
	return &characterHandler{svc: svc}
}

func (h *characterHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCharacterRequest
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

func (h *characterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid path value", http.StatusBadRequest)
		return
	}

	var req models.UpdateCharacterRequest
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

func (h *characterHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	page, err := parsePagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	filters := &models.CharacterFilter{
		Name:       strings.TrimSpace(query.Get("name")),
		Status:     strings.TrimSpace(query.Get("status")),
		Species:    strings.TrimSpace(query.Get("species")),
		Grade:      strings.TrimSpace(query.Get("grade")),
		Gender:     strings.TrimSpace(query.Get("gender")),
		Pagination: page,
	}

	c, total, err := h.svc.FindAll(filters)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeJson(w, http.StatusOK, models.NewPaginated(c, page, total)); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *characterHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid path value", http.StatusBadRequest)
		return
	}
	c, err := h.svc.FindByID(id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not_found:") {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err = writeJson(w, http.StatusOK, c); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *characterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid path value", http.StatusBadRequest)
		return
	}

	if err = h.svc.Delete(id); err != nil {
		var status int
		switch {
		case strings.HasPrefix(err.Error(), "not_found:"):
			status = http.StatusNotFound
		case strings.HasPrefix(err.Error(), "forbidden:"):
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
