package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/config"
	"ne-pridumal/effective-mobile-test/internal/models"
	"net/http"
)

type ApiCaller struct {
	conf   config.Api
	logger *slog.Logger
}

var (
	ErrBadRequest     = errors.New("Error bad request")
	ErrInternalServer = errors.New("Internal server error")
	ErrNoResponce     = errors.New("Error no responce")
)

func New(c config.Api, l *slog.Logger) *ApiCaller {
	return &ApiCaller{
		conf:   c,
		logger: l,
	}
}

func (c *ApiCaller) GetPassportData(num int) (*models.ApiUser, error) {
	var u models.ApiUser
	r, err := http.Get(fmt.Sprintf("%s/info?name=%d", c.conf.Passport, num))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		switch r.StatusCode {
		case http.StatusBadRequest:
			return nil, ErrBadRequest
		case http.StatusInternalServerError:
			return nil, ErrInternalServer
		default:
			return nil, ErrNoResponce
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}
