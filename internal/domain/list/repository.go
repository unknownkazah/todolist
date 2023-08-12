package list

import "context"

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
	Create(ctx context.Context, data Entity) (id string, err error)
	Get(ctx context.Context, id string) (dest Entity, err error)
	Status(ctx context.Context, id string, data Entity) (err error)
	Update(ctx context.Context, id string, data Entity) (err error)
	Delete(ctx context.Context, id string) (err error)
}
