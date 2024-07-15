package postgres

import (
	"database/sql"
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/models"
	"time"
)

type userRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func (r *userRepository) Create(u *models.User) error {
	var id int
	query := `
		INSER INTO users(name, surname, address, passport, patronomic)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	if err := r.db.QueryRow(
		query,
		u.Name,
		u.Surname,
		u.Address,
		u.Passport,
		u.Patr,
	).Scan(&id); err != nil {
		return err
	}

	u.Id = id

	return nil
}

// filter params: name, surname, address, passport
func (r *userRepository) Get(lm int, u models.User) ([]models.User, error) {
	var l []models.User
	var query string

	query = `                                             
		SELECT * FROM users WHERE name ILIKE '%' || $2 || '%' 
			AND surname ILIKE '%' || $3 || '%'                  
			AND address ILIKE '%' || $4 || '%'
			AND passport ILIKE '%' || $5 || '%'
			AND patronomic ILIKE '%' || $6 || '%'
		LIMIT $1                                            
	`
	rows, err := r.db.Query(
		query,
		lm,
		u.Name,
		u.Surname,
		u.Address,
		u.Passport,
		u.Patr,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var us models.User
		if err := rows.Scan(
			&us.Id,
			&us.Name,
			&us.Passport,
			&us.Surname,
			&us.Address,
		); err != nil {
			return nil, err
		}
		l = append(l, us)
	}

	return l, nil
}

func (r *userRepository) GetTasks(id int, start, end time.Time) ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM users_tasks 
		WHERE user_id = $1 
	AND start_date >= $2 AND end_date <= $3
		ORDER BY duration`
	rows, err := r.db.Query(query, id, start, end)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.Id,
			&t.UserId,
			&t.StartDate,
			&t.EndDate,
			&t.LastStart,
			&t.Duration,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *userRepository) Update(u models.User) error {
	query := `UPDATE users 
		SET	name = $2, surname = $3, address = $4, passport = $5,patronomic = $6 
		WHERE id = $1
	`

	if err := r.db.QueryRow(
		query,
		u.Id,
		u.Name,
		u.Surname,
		u.Address,
		u.Passport,
		u.Patr,
	).Err(); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	if err := r.db.QueryRow(query, id).Err(); err != nil {
		return err
	}
	return nil
}
