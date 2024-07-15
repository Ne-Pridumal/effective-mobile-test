package httpServer

import (
	"errors"
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/models"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	ErrBadRequest     = errors.New("Error bad request")
	ErrInternalServer = errors.New("Internal server error")
	ErrNoResponce     = errors.New("Error no responce")
)

type HttpServer struct {
	serveMux      *http.ServeMux
	usersHandlers *usersHandlers
	tasksHandlers *tasksHandlers
}

type apiCaller interface {
	GetPassportData(n int) (*models.ApiUser, error)
}

type usersRepository interface {
	Create(u *models.User) error
	//filter params: limit(int), name, surname, address, passport
	Get(l int, p models.User) ([]models.User, error)
	Delete(id int) error
	GetTasks(id int, start, end time.Time) ([]models.Task, error)
	Update(models.User) error
}

type tasksRepository interface {
	Get(id int) (*models.Task, error)
	Create(userId int) (*models.Task, error)
	StartTracking(id int) error
	StopTracking(id int) error
}

func New(uR usersRepository, tR tasksRepository, apiCaller apiCaller, logger *slog.Logger) *HttpServer {
	m := http.NewServeMux()
	s := &HttpServer{
		serveMux: m,
		usersHandlers: &usersHandlers{
			rep:       uR,
			logger:    logger,
			apiCaller: apiCaller,
		},
		tasksHandlers: &tasksHandlers{
			rep:    tR,
			logger: logger,
		},
	}
	s.setRoutes()
	return s
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *HttpServer) setRoutes() {
	s.serveMux.HandleFunc("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/docs/doc.json"),
	))

	s.serveMux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.usersHandlers.Get(w, r)
		case http.MethodDelete:
			s.usersHandlers.Delete(w, r)
		case http.MethodPut:
			s.usersHandlers.Update(w, r)
		default:
			return
		}
	})
	s.serveMux.HandleFunc("/users/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.usersHandlers.GetTasks(w, r)
		}
	})
	s.serveMux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.tasksHandlers.Get(w, r)
		case http.MethodPost:
			s.tasksHandlers.Create(w, r)
		default:
			return
		}
	})
	s.serveMux.HandleFunc("/tasks/tracking", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.tasksHandlers.Tracking(w, r)
		default:
			return
		}
	})
}
