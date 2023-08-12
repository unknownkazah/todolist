package list

import "context"

type Cache interface {
	Get(ctx context.Context, id string) (dest Entity, err error)
	Status(ctx context.Context, id string, data Entity) (err error)
}
