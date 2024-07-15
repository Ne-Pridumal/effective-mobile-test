package postgres

import (
	"database/sql"
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/models"
	sl "ne-pridumal/effective-mobile-test/lib/logger/slog"
	"time"
)

type tasksRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func (r *tasksRepository) Get(id int) (*models.Task, error) {
	t := &models.Task{}
	query := `
		SELECT * FROM users_tasks WHERE id = $1
	`
	if err := r.db.QueryRow(
		query,
		id,
	).Scan(
		&t.Id,
		&t.UserId,
		&t.StartDate,
		&t.EndDate,
		&t.LastStart,
		&t.Duration,
	); err != nil {
		return nil, err
	}
	return t, nil
}

func (r *tasksRepository) Create(userId int) (*models.Task, error) {
	var id int
	startTime := time.Now().UTC()
	query := `
		INSERT INTO users_tasks(user_id, start_date, end_date, last_start, duration)
		VALUES($1,$2,$3,$4,$5)
		RETURNING id
	`
	if err := r.db.QueryRow(
		query,
		userId,
		startTime,
		startTime,
		startTime,
		0,
	).Scan(&id); err != nil {
		return nil, err
	}
	return &models.Task{
		Id:        id,
		UserId:    userId,
		StartDate: startTime,
		EndDate:   startTime,
		LastStart: startTime,
		Duration:  0,
	}, nil
}

func (r *tasksRepository) StartTracking(id int) error {
	var err error
	newTime := time.Now()
	query := `
		UPDATE users_tasks SET last_start=$2 WHERE id = $1
	`
	err = r.db.QueryRow(
		query,
		id,
		newTime,
	).Err()

	return err
}

// update end_date and duration fields, duration += time.Now() - last_start,
// duration - amount of minutes
func (r *tasksRepository) StopTracking(id int) error {
	const op = "storage.postgres.tasksRepository.StopTracking"
	query := `
		UPDATE users_tasks SET end_date = $2, duration = duration + $3 WHERE id = $1
	`
	t, err := r.Get(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	d := int(now.Sub(t.LastStart).Minutes())
	r.logger.Debug(op+" time calc", sl.Debug(int(now.Sub(t.LastStart).Minutes())))
	if err := r.db.QueryRow(
		query,
		id,
		now,
		d,
	).Err(); err != nil {
		return err
	}

	return nil
}
