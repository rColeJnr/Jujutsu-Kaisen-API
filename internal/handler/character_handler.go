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
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.Insert(&req)
	if err != nil {
		if isValidationErr(err) {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *characterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
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
		writeError(w, status, err)
		return
	}

	if err = json.NewEncoder(w).Encode(req); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *characterHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.FindAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err = writeJson(w, http.StatusOK, c); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *characterHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
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
		writeError(w, http.StatusBadRequest, err)
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
		writeError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
