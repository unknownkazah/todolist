package todo

import (
	"context"
	"todo/internal/domain/list"
	"todo/pkg/log"

	"go.uber.org/zap"
)

func (s *Service) List(ctx context.Context) (res []list.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("List")

	data, err := s.listRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}
	res = list.ParseFromEntities(data)

	return
}

func (s *Service) GetStatus(ctx context.Context, id string, req list.Request) (err error) {
	data := list.Entity{
		Status: &req.Status,
	}

	err = s.listRepository.Status(ctx, id, data)
	if err != nil {
		return
	}

	return
}

func (s *Service) Create(ctx context.Context, req list.Request) (res list.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("Create")

	data := list.Entity{
		Title:    &req.Title,
		ActiveAt: &req.ActiveAt,
	}

	data.ID, err = s.listRepository.Create(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}
	res = list.ParseFromEntity(data)

	return
}

func (s *Service) Get(ctx context.Context, id string) (res list.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("Get").With(zap.String("id", id))

	data, err := s.listRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}
	res = list.ParseFromEntity(data)

	return
}

func (s *Service) Update(ctx context.Context, id string, req list.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("Update").With(zap.String("id", id))

	data := list.Entity{
		Title:    &req.Title,
		ActiveAt: &req.ActiveAt,
	}

	err = s.listRepository.Update(ctx, id, data)
	if err != nil {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) Delete(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("Delete").With(zap.String("id", id))

	err = s.listRepository.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}
