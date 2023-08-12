package memory

import (
	"context"
	"database/sql"
	"sync"
	"todo/internal/domain/list"

	"github.com/google/uuid"
)

type ListRepository struct {
	db map[string]list.Entity
	sync.RWMutex
}

func NewListRepository() *ListRepository {
	return &ListRepository{
		db: make(map[string]list.Entity),
	}
}

func (r *ListRepository) List(ctx context.Context) (dest []list.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]list.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *ListRepository) Create(ctx context.Context, data list.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *ListRepository) Get(ctx context.Context, id string) (dest list.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = sql.ErrNoRows
		return
	}

	return
}

func (r *ListRepository) Update(ctx context.Context, id string, data list.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	r.db[id] = data

	return
}

func (r *ListRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	delete(r.db, id)

	return
}

func (r *ListRepository) generateID() string {
	return uuid.New().String()
}
