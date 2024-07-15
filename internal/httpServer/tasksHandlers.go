package httpServer

import (
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/models"
	resp "ne-pridumal/effective-mobile-test/lib/api/response"
	sl "ne-pridumal/effective-mobile-test/lib/logger/slog"
	"net/http"
	"strconv"
	"strings"
)

type tasksHandlers struct {
	rep    tasksRepository
	logger *slog.Logger
}

// Get godoc
// @Summary      Get task
// @Description  get task by id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id query string false "task's id"
// @Success      200  {object}  models.Task
// @Router       /tasks [get]
func (h *tasksHandlers) Get(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.tasksHandlers.Get"
	type response struct {
		resp.Response
		Data *models.Task
	}
	sId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(sId)
	if err != nil {
		h.logger.Error("wrong id", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 400, resp.Error("wrong id"))
		return
	}

	res, err := h.rep.Get(id)
	if err != nil {
		h.logger.Error("bd error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}

	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
		Data:     res,
	})
}

// Create godoc
// @Summary      create task
// @Description  create task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        user-id body int true "user's id"
// @Success      200  {object}  models.Task
// @Router       /tasks [post]
func (h *tasksHandlers) Create(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.tasksHandlers.Create"
	type response struct {
		resp.Response
		Data *models.Task
	}
	type request struct {
		UserId int `json:"user-id"`
	}
	req, err := resp.Decode[request](r)
	if err != nil {
		h.logger.Error("wrong id", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 400, resp.Error("wrong id"))
		return
	}

	res, err := h.rep.Create(req.UserId)
	if err != nil {
		h.logger.Error("bd error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}

	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
		Data:     res,
	})
}

// Tracking godoc
// @Summary      task's tracking
// @Description  start and stop task's tracking
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id body int true "task's id"
// @Param        command body string true "either start or stop"
// @Success      200  {object}  models.Task
// @Router       /tasks/tracking [post]
func (h *tasksHandlers) Tracking(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.tasksHandlers.StartTracking"
	type response struct {
		resp.Response
	}
	type request struct {
		Id   int    `json:"id"`
		Comm string `json:"command"`
	}
	req, err := resp.Decode[request](r)
	if err != nil {
		h.logger.Error("wrong id", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 400, resp.Error("wrong id"))
		return
	}

	switch strings.ToLower(req.Comm) {
	case "start":
		err = h.rep.StartTracking(req.Id)
		if err != nil {
			h.logger.Error("bd error", sl.Err(sl.ErrWrapper(err, op)))
			resp.Encode(w, r, 500, resp.Error("error during db query"))
			return
		}
	case "stop":
		err = h.rep.StopTracking(req.Id)
		if err != nil {
			h.logger.Error("bd error", sl.Err(sl.ErrWrapper(err, op)))
			resp.Encode(w, r, 500, resp.Error("error during db query"))
			return
		}
	}

	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
	})
}
