package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/murkh/gig-mobile-backend/internal/core/ports"
	"github.com/murkh/gig-mobile-backend/internal/platform/errors"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get("/users/{id}", h.GetUser)
	r.Post("/users", h.CreateUser)
	r.Get("/users", h.ListUsers)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	opts := parseQueryOptions(r)
	user, err := h.service.GetUser(r.Context(), id, opts...)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, errors.BadRequest("invalid request body"))
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	opts := parseQueryOptions(r)
	users, err := h.service.ListUsers(r.Context(), opts...)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func parseQueryOptions(r *http.Request) []ports.QueryOptions {
	includes := r.URL.Query()["include"]
	if len(includes) > 0 {
		return []ports.QueryOptions{{Includes: includes}}
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		appErr = errors.InternalServerError(err)
	}

	respondWithJSON(w, appErr.Code, appErr)
}
