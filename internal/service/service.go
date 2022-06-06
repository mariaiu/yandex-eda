package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mariaiu/yandex-eda-app/internal/config"
	m "github.com/mariaiu/yandex-eda-app/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	repo Repository
	logger *logrus.Logger
	endpoint string
	minRating float64
}

func NewService(repo Repository, logger *logrus.Logger, app *config.App) *Service {
	return &Service{
		repo: repo,
		logger: logger,
		endpoint: app.Endpoint,
		minRating: app.MinRating,
	}
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 5 * time.Second}
	return client
}

func (svc *Service) sendRequest(ctx context.Context, client *http.Client, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil); if err != nil {
		svc.logger.Println(req.Method, req.URL, err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")

	res, err := client.Do(req); if err != nil {
		svc.logger.Println(req.Method, req.URL, err.Error())
		return err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		svc.logger.Println(req.Method, req.URL, res.StatusCode, err.Error())
		return err
	}

	svc.logger.Println(req.Method, req.URL, res.StatusCode)

	return nil
}

func (svc *Service) GetRestaurants(ctx context.Context) ([]m.Restaurant, error) {
	return svc.repo.GetAllRestaurants(ctx)
}

func (svc *Service) GetPositions(ctx context.Context, id int) ([]m.Position, error) {
	return svc.repo.GetPositions(ctx, id)
}

var mu sync.Mutex

func (svc *Service) ParseRestaurants(ctx context.Context, req *m.ParseRequest) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	type processedResult struct {
		count int
		err error
	}

	var result processedResult

	var processDone = make(chan struct{})
	go func() {
		var client = httpClient()

		chRestaurants, err := svc.requestRestaurants(ctxTimeout, client, req); if err != nil {
			result.err = err
			processDone <- struct{}{}
			return
		}
		chPositions := svc.requestPositions(ctxTimeout, chRestaurants, client, req)
		chCounter := svc.addToDB(ctxTimeout, chPositions)

		for range chCounter {
			result.count++
		}

		processDone <- struct{}{}
	}()

	select {
	case <-ctxTimeout.Done():
		return 0, errors.New("process timeout")
	case  <-processDone:
		return result.count, result.err
	}
}

func (svc *Service) requestRestaurants(ctx context.Context, client *http.Client, req *m.ParseRequest) (chan m.Restaurant, error) {
	var (
		yandexRestaurants m.YandexRestaurantsResponse
		urlRestaurants = fmt.Sprintf("%s?latitude=%f&longitude=%f", svc.endpoint, req.Latitude, req.Longitude)
	)

	if err := svc.sendRequest(ctx, client, urlRestaurants,  &yandexRestaurants); err != nil {
		return nil, err
	}

	output := make(chan m.Restaurant, req.Workers)

	go func() {
		for _, place := range yandexRestaurants.Payload.FoundPlaces {
			if place.Rating >= svc.minRating {
				output <- place.Restaurant
			}
		}

		close(output)

	}()

	return output, nil
}

func (svc *Service) requestPositions(ctx context.Context, input <- chan m.Restaurant, client *http.Client, req *m.ParseRequest) chan m.RestaurantPositions {
	output := make(chan m.RestaurantPositions)

	wg := new(sync.WaitGroup)
	wg.Add(req.Workers)

	go func() {
		for worker := 0; worker < req.Workers; worker++ {
			go func() {
				defer wg.Done()
				for restaurant := range input {
					select {
					case <-ctx.Done():
						return
					default:
						var (
							url = fmt.Sprintf("%s/%s/menu?latitude=%f&longitude=%f", svc.endpoint, restaurant.Slug, req.Latitude, req.Longitude)
							yandexPositions m.YandexPositionsResponse
						)

						if err := svc.sendRequest(ctx, client, url,  &yandexPositions); err != nil {
							continue
						}

						var positions []m.Position

						for _, category := range yandexPositions.Payload.Categories {
							for _, item := range category.Items {
								if !item.MatchFoodParameter() {
									continue
								}

								if err := item.ConvertWeightToInt(); err != nil {
									svc.logger.Printf("error converting weight for url:, weight: %v, %v",item.Weight, err.Error())
									continue
								}
								positions = append(positions, item)
							}
						}

						if len(positions) > 0 {
							output <- m.RestaurantPositions{Restaurant: restaurant, Positions: positions}
						}
					}
				}
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func (svc *Service) addToDB(ctx context.Context, input <- chan m.RestaurantPositions) chan struct{} {
	output := make(chan struct{})
	wg := new(sync.WaitGroup)

	go func(){
		for value := range input {
			wg.Add(1)
			go func(value m.RestaurantPositions) {
				defer wg.Done()

				if err := svc.repo.AddRestaurant(ctx, &value.Restaurant); err != nil {
					svc.logger.Printf("error saving restaurant to db: %s", err.Error())
					return
				}

				if err := svc.repo.AddPosition(ctx, &value.Restaurant, value.Positions); err != nil {
					svc.logger.Printf("error saving positions to db: %s", err.Error())
					return
				}

				output <- struct{}{}
			}(value)
		}

		wg.Wait()

		close(output)
	}()

	return output
}

