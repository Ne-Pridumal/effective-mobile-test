package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/config"

	sl "ne-pridumal/effective-mobile-test/lib/logger/slog"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db              *sql.DB
	logger          *slog.Logger
	usersRepository *userRepository
	tasksRepository *tasksRepository
}

func New(conf config.Postgres, l *slog.Logger) (*Storage, error) {
	psqlConf := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
		conf.Address,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Db,
		conf.Ssl,
	)
	const op = "storage.postgres.New"
	sqldb, err := sql.Open("postgres", psqlConf)
	if err != nil {
		return nil, sl.ErrWrapper(err, op)
	}

	return &Storage{
		Db:     sqldb,
		logger: l,
	}, nil
}

func (s *Storage) Close() error {
	return s.Db.Close()
}

func (s *Storage) Users() *userRepository {
	if s.usersRepository != nil {
		return s.usersRepository
	}
	s.usersRepository = &userRepository{
		db:     s.Db,
		logger: s.logger,
	}
	return s.usersRepository
}

func (s *Storage) Tasks() *tasksRepository {
	if s.tasksRepository != nil {
		return s.tasksRepository
	}
	s.tasksRepository = &tasksRepository{
		db:     s.Db,
		logger: s.logger,
	}
	return s.tasksRepository
}
