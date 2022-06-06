package service

import (
	"context"
	"github.com/mariaiu/yandex-eda-app/internal/models"
)

type Repository interface {
	GetAllRestaurants(ctx context.Context) ([]models.Restaurant, error)
	GetPositions(ctx context.Context, id int) ([]models.Position, error)
	AddPosition(ctx context.Context, restaurant *models.Restaurant, positions []models.Position) error
	AddRestaurant(ctx context.Context, restaurant *models.Restaurant) error
}