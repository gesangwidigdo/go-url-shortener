package url

import (
	"errors"
	"net/http"

	"github.com/gesangwidigdo/url-shortener/internal/platform/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/urls", h.create)
	mux.HandleFunc("GET /api/urls/{code}", h.get)
	mux.HandleFunc("GET /api/{code}", h.redirect)
	mux.HandleFunc("GET /{code}", h.redirect)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var form CreateForm
	if err := response.Decode(w, r, &form); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.svc.Create(r.Context(), form)
	switch {
	case errors.Is(err, ErrInvalidURL):
		response.Error(w, http.StatusUnprocessableEntity, "url must be a valid http or https address")
	case errors.Is(err, ErrInvalidAlias):
		response.Error(w, http.StatusUnprocessableEntity, "alias must be 3-32 letters, digits, - or _")
	case errors.Is(err, ErrInvalidExpiry):
		response.Error(w, http.StatusUnprocessableEntity, "expires_in_hours must be between 1 and 8760")
	case errors.Is(err, ErrCodeTaken):
		response.Error(w, http.StatusConflict, "alias already in use")
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	default:
		response.JSON(w, http.StatusCreated, u)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	u, err := h.svc.GetByCode(r.Context(), r.PathValue("code"))
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, "link not found")
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	default:
		response.JSON(w, http.StatusOK, u)
	}
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	target, err := h.svc.Resolve(r.Context(), r.PathValue("code"))
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, "link not found or expired")
	case err != nil:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	default:
		http.Redirect(w, r, target, http.StatusFound) // 302
	}
}