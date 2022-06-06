package handler

import (
	"context"
	"github.com/mariaiu/yandex-eda-app/internal/models"
)

type Service interface {
	GetRestaurants(ctx context.Context) ([]models.Restaurant, error)
	GetPositions(ctx context.Context, id int) ([]models.Position, error)
	ParseRestaurants(ctx context.Context, req *models.ParseRequest) (int, error)
}
