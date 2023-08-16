package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"
	"todo/internal/domain/list"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) List(ctx context.Context) (dest []list.Entity, err error) {
	query := `
		SELECT title, active_at
		FROM items
		WHERE DATE(active_at) = DATE($1)
		ORDER BY id`

	currentDate := time.Now().Format("2006-01-02")

	err = s.db.SelectContext(ctx, &dest, query, currentDate)

	return
}

func (s *Repository) Status(ctx context.Context, id string, data list.Entity) (err error) {
	query := `
		UPDATE items
		SET status = true
		WHERE id = $1
	`

	_, err = s.db.ExecContext(ctx, query, id)
	trueValue := true
	data.Status = &trueValue
	return err
}

func (s *Repository) Create(ctx context.Context, data list.Entity) (id string, err error) {
	query := `
		INSERT INTO items (title, active_at)
		VALUES ($1, $2)
		RETURNING id`

	args := []any{data.Title, data.ActiveAt}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *Repository) Get(ctx context.Context, id string) (dest list.Entity, err error) {
	query := `
		SELECT title, active_at
		FROM items
		WHERE id=$1`

	args := []any{id}

	err = s.db.GetContext(ctx, &dest, query, args...)

	return
}

func (s *Repository) Update(ctx context.Context, id string, data list.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE items SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = s.db.ExecContext(ctx, query, args...)
	}

	return
}

func (s *Repository) prepareArgs(data list.Entity) (sets []string, args []any) {
	if data.Title != nil {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("title=$%d", len(args)))
	}

	if data.ActiveAt != nil {
		args = append(args, data.ActiveAt)
		sets = append(sets, fmt.Sprintf("active_at=$%d", len(args)))
	}

	return
}

func (s *Repository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM items
		WHERE id=$1`

	args := []any{id}

	_, err = s.db.ExecContext(ctx, query, args...)

	return
}
