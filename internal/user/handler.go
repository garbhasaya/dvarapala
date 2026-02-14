package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"dvarapala/internal/platform/auth"
	"dvarapala/internal/platform/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler handles HTTP requests for users.
type Handler struct {
	svc      Service
	validate *validator.Validate
}

// NewHandler creates a new user handler.
func NewHandler(svc Service) *Handler {
	return &Handler{
		svc:      svc,
		validate: validator.New(),
	}
}

// Routes returns the chi router for user endpoints.
func (h *Handler) Routes(jwtManager *auth.JWTManager) chi.Router {
	r := chi.NewRouter()

	// Public routes
	r.Post("/auth", h.Authenticate)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(jwtManager))
		
		r.Post("/", h.Create)
		r.Get("/", h.List)
		
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetByID)
			r.Post("/", h.Update) // README said POST for update
			r.Delete("/", h.Delete)
		})
	})

	return r
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} render.Response{data=User}
// @Failure 400 {object} render.Response
// @Failure 401 {object} render.Response
// @Failure 500 {object} render.Response
// @Security Bearer
// @Router /users [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode create user request", "error", err)
		render.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		slog.Warn("invalid create user request", "error", err)
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.svc.Create(r.Context(), req)
	if err != nil {
		// slog.Error is already called in service
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, u)
}

// List godoc
// @Summary List all users
// @Description Get a list of all registered users
// @Tags users
// @Produce json
// @Success 200 {object} render.Response{data=[]User}
// @Failure 401 {object} render.Response
// @Failure 500 {object} render.Response
// @Security Bearer
// @Router /users [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, users)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get a single user by their unique ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} render.Response{data=User}
// @Failure 400 {object} render.Response
// @Failure 401 {object} render.Response
// @Failure 404 {object} render.Response
// @Security Bearer
// @Router /users/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("invalid user id in request", "id", idStr)
		render.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		slog.Warn("user not found", "id", id)
		render.Error(w, http.StatusNotFound, "user not found")
		return
	}

	render.JSON(w, http.StatusOK, u)
}

// Update godoc
// @Summary Update user
// @Description Update an existing user's details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "Updated user details"
// @Success 200 {object} render.Response{data=User}
// @Failure 400 {object} render.Response
// @Failure 401 {object} render.Response
// @Failure 500 {object} render.Response
// @Security Bearer
// @Router /users/{id} [post]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("invalid user id in update request", "id", idStr)
		render.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode update user request", "id", id, "error", err)
		render.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		slog.Warn("invalid update user request", "id", id, "error", err)
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, u)
}

// Delete godoc
// @Summary Delete user
// @Description Remove a user from the system by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} render.Response
// @Failure 401 {object} render.Response
// @Failure 500 {object} render.Response
// @Security Bearer
// @Router /users/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("invalid user id in delete request", "id", idStr)
		render.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusNoContent, nil)
}

// Authenticate godoc
// @Summary Authenticate user
// @Description Login with email and password to receive a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body AuthRequest true "Login credentials"
// @Success 200 {object} render.Response{data=AuthResponse}
// @Failure 400 {object} render.Response
// @Failure 401 {object} render.Response
// @Router /users/auth [post]
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode auth request", "error", err)
		render.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		slog.Warn("invalid auth request", "error", err)
		render.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Authenticate(r.Context(), req)
	if err != nil {
		render.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, res)
}
