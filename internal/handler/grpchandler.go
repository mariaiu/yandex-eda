package handler

import (
	"context"
	"github.com/mariaiu/yandex-eda-app/internal/config"
	m "github.com/mariaiu/yandex-eda-app/internal/models"
	"github.com/mariaiu/yandex-eda-app/internal/validator"
	pb "github.com/mariaiu/yandex-eda-app/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	pb.UnimplementedYandexEdaServer
	srv Service
	logger *logrus.Logger
	validator *validating.Validator
	latitude float64
	longitude float64
	maxWorkers int
}

func NewGRPCHandler(srv Service, logger *logrus.Logger, app *config.App, validator *validating.Validator) *GRPCHandler {
	return &GRPCHandler{
		srv: srv,
		logger: logger,
		validator: validator,
		latitude: app.Latitude,
		longitude: app.Longitude,
		maxWorkers: app.MaxWorkers,

	}
}

func (g *GRPCHandler) GetRestaurants(ctx context.Context, _ *emptypb.Empty) (*pb.GetRestaurantsResponse, error) {
	restaurants, err := g.srv.GetRestaurants(ctx); if err != nil {
		return nil, status.Error(500, err.Error())
	}

	if len(restaurants) == 0 {
		return nil, status.Error(404,"no restaurants found")
	}

	restaurantsResp := make([]*pb.GetRestaurantsResponse_Restaurant, 0, len(restaurants))

	for _, restaurant := range restaurants {
		restaurantsResp = append(restaurantsResp, &pb.GetRestaurantsResponse_Restaurant{
			Id: int64(restaurant.Id),
			Name: restaurant.Name,
			Slug: restaurant.Slug,
			DeliveryPrice: restaurant.MinimalDeliveryCost,
			Rating: restaurant.Rating,
		})
	}

	return &pb.GetRestaurantsResponse{Restaurants: restaurantsResp}, nil
}

func (g *GRPCHandler) GetRestaurant(ctx context.Context, r *pb.GetRestaurantRequest) (*pb.GetRestaurantResponse, error) {
	restaurantId := int(r.GetId())
	positions, err := g.srv.GetPositions(ctx, restaurantId); if err != nil {
		return nil, status.Error(500, err.Error())
	}

	if len(positions) == 0 {
		return nil, status.Error(404,"restaurant not found")
	}

	positionsResp :=  make([]*pb.GetRestaurantResponse_Position, 0, len(positions))

	for _, position := range positions {
		positionsResp = append(positionsResp, &pb.GetRestaurantResponse_Position{
			Name: position.Name,
			Description: position.Description,
			Price: int64(position.Price),
			Weight: position.Weight.(int64),
		})
	}

	return &pb.GetRestaurantResponse{Positions: positionsResp}, nil
}

func (g *GRPCHandler) ParseRestaurants(ctx context.Context, r *pb.ParseRestaurantsRequest) (*pb.ParseRestaurantsResponse, error) {
	req := m.ParseRequest{
		Longitude: g.longitude,
		Latitude: g.latitude,
		Workers: g.maxWorkers,
	}

	if r.Longitude != nil {
		req.Longitude = r.GetLongitude()
	}

	if r.Latitude != nil {
		req.Latitude = r.GetLatitude()
	}

	if r.Workers != nil {
		req.Workers = int(r.GetWorkers())
	}

	if err := g.validator.ValidateStruct(req); err != nil {
		return nil, status.Error(422, err.Error())
	}

	count, err := g.srv.ParseRestaurants(ctx, &req); if err != nil {
		switch err.Error() {
		case "process timeout":
			return nil, status.Error(504, err.Error())
		default:
			return nil, status.Error(500, err.Error())
		}
	}

	return &pb.ParseRestaurantsResponse{WereProcessed: int32(count)}, nil
}