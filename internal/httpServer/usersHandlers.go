package httpServer

import (
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/models"
	resp "ne-pridumal/effective-mobile-test/lib/api/response"
	"ne-pridumal/effective-mobile-test/lib/json"
	sl "ne-pridumal/effective-mobile-test/lib/logger/slog"
	"net/http"
	"strconv"
	"strings"
)

type usersHandlers struct {
	rep       usersRepository
	logger    *slog.Logger
	apiCaller apiCaller
}

// GetUser godoc
// @Summary      Get users
// @Description  get filtered users by params
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        passport query string false "user passport"
// @Param				 name query string false "user name"
// @Param				 address query string false "user address"
// @Param				 surname query string false "user surname"
// @Param				 patronomic query string false "user partranomic"
// @Success      200  {object}  []models.User
// @Router       /users [get]
func (h *usersHandlers) Get(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.usersHandler.Get"
	type response struct {
		resp.Response
		Data []models.User `json:"data"`
	}

	passp, surn, name, address := r.URL.Query().Get("passport"),
		r.URL.Query().Get("surname"),
		r.URL.Query().Get("name"),
		r.URL.Query().Get("address")
	r.URL.Query().Get("patronomic")

	stLim := r.URL.Query().Get("limit")

	var lim int
	var err error

	if stLim != "" {
		lim, err = strconv.Atoi(stLim)

		if err != nil {
			resp.Encode(w, r, 400, resp.Error("wrong limit"))
			h.logger.Error("error with converting limit", sl.Err(sl.ErrWrapper(err, op)))
			return
		}
	} else {
		lim = 10
	}

	h.logger.Debug("params", sl.Debug(passp, surn, name, address, lim))

	res, err := h.rep.Get(lim, models.User{
		Name:     name,
		Address:  address,
		Passport: passp,
		Surname:  surn,
	})
	if err != nil {
		h.logger.Error("db query error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}

	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
		Data:     res,
	})
}

// GetTasks godoc
// @Summary      Get user's tasks
// @Description  get user's tasks
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  body int true "user id"
// @Param				 start body time.Time true "period start"
// @Param				 end body time.Time true "period end"
// @Success      200  {object} []models.Task
// @Router       /users/tasks [post]
func (h *usersHandlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.usersHandler.GetTasks"
	type data struct {
		Sum   int           `json:"sum"`
		Tasks []models.Task `json:"tasks"`
	}

	type response struct {
		resp.Response
		Data data `json:"data"`
	}

	type request struct {
		Id    int             `json:"user-id"`
		Start json.CustomDate `json:"start"`
		End   json.CustomDate `json:"end"`
	}

	req, err := resp.Decode[request](r)
	if err != nil {
		resp.Encode(w, r, 400, resp.Error(err.Error()))
		h.logger.Error("error with converting data", sl.Err(sl.ErrWrapper(err, op)))
		return
	}
	res, err := h.rep.GetTasks(req.Id, req.Start.Time, req.End.Time)
	if err != nil {
		h.logger.Error("db query error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}
	sum := 0
	for _, v := range res {
		sum += v.Duration
	}
	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
		Data: data{
			Sum:   sum,
			Tasks: res,
		},
	})
}

// Create godoc
// @Summary      Create user
// @Description  create new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        passportNumber  body string true "user's passport number"
// @Success      200  {object} []models.Task
// @Router       /users [post]
func (h *usersHandlers) Create(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.usersHandler.Create"
	type response struct {
		resp.Response
		Data *models.User
	}

	type request struct {
		PasNum string `json:"passportNumber"`
	}

	req, err := resp.Decode[request](r)
	if err != nil {
		resp.Encode(w, r, 400, resp.Error("wrong passport number"))
		h.logger.Error("wrong passport number", sl.Err(sl.ErrWrapper(err, op)))
		return
	}
	ps, err := strconv.Atoi(strings.ReplaceAll(req.PasNum, " ", ""))
	if err != nil {
		resp.Encode(w, r, 400, resp.Error("wrong passport number"))
		h.logger.Error("error with converting passport number", sl.Err(sl.ErrWrapper(err, op)))
		return
	}

	res, err := h.apiCaller.GetPassportData(ps)
	if err != nil {
		resp.Encode(w, r, 500, resp.Error("smth went wrong"))
		h.logger.Error("error with external api call", sl.Err(sl.ErrWrapper(err, op)))
		return
	}
	u := &models.User{
		Passport: req.PasNum,
		Address:  res.Address,
		Name:     res.Name,
		Surname:  res.Surname,
		Patr:     res.Patr,
	}
	if err := h.rep.Create(u); err != nil {
		h.logger.Error("db query error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}

	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
		Data:     u,
	})
}

// Update godoc
// @Summary      Update user
// @Description  update user data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body models.User true "user's new data"
// @Router       /users [put]
func (h *usersHandlers) Update(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.usersHandler.Update"
	type response struct {
		resp.Response
	}
	type request struct {
		models.User
	}
	req, err := resp.Decode[request](r)
	if err != nil {
		resp.Encode(w, r, 400, resp.Error("wrong data"))
		h.logger.Error("error with converting data", sl.Err(sl.ErrWrapper(err, op)))
		return
	}

	err = h.rep.Update(req.User)
	if err != nil {
		h.logger.Error("db query error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}

	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
	})
}

// Delete godoc
// @Summary      Delete user
// @Description  delete user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  body int true "user's id"
// @Router       /users [put]
func (h *usersHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "httpServer.usersHandler.Delete"
	type response struct {
		resp.Response
	}
	type request struct {
		Id int `json:"id"`
	}

	req, err := resp.Decode[request](r)
	if err != nil {
		resp.Encode(w, r, 400, resp.Error("no id"))
		h.logger.Error("error with converting id", sl.Err(sl.ErrWrapper(err, op)))
		return
	}

	err = h.rep.Delete(req.Id)
	if err != nil {
		h.logger.Error("db query error", sl.Err(sl.ErrWrapper(err, op)))
		resp.Encode(w, r, 500, resp.Error("error during db query"))
		return
	}
	resp.Encode(w, r, 200, response{
		Response: resp.OK(),
	})
}
