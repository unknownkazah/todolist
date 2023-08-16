package http

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"todo/internal/domain/list"
	"todo/internal/service/todo"
	"todo/pkg/server/response"
)

type Handler struct {
	todoService *todo.Service
}

func NewHandler(s *todo.Service) *Handler {
	return &Handler{todoService: s}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/done", h.status)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// List of tasks from the database
//

//	@Summary	Status the tasks from the database
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"path param"
//	@Success	200	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/tasks/{id}/done [put]
func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := list.Request{
		Status: true,
	}

	err := h.todoService.GetStatus(r.Context(), id, req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
}

//	@Summary	List of tasks from the database
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Success	200		{array}		response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/tasks 	[get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.todoService.List(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Add a new task to the database
//
//	@Summary	Add a new task to the database
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		request	body		list.Request	true	"body param"
//	@Success	200		{object}	response.Object
//	@Failure	400		{object}	response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/tasks [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	req := list.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.todoService.Create(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Update the task in the database
//
//	@Summary	Update the task in the database
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string			true	"path param"
//	@Param		request	body	list.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/tasks/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := list.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	if err := h.todoService.Update(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}
}

// Delete the task from the database
//
//	@Summary	Delete the task from the database
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"path param"
//	@Success	200
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/tasks/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.todoService.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}
}
