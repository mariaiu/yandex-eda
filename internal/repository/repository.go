package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	m "github.com/mariaiu/yandex-eda-app/internal/models"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) GetAllRestaurants(ctx context.Context) ([]m.Restaurant, error) {
	var restaurants []m.Restaurant

	query := fmt.Sprint(`SELECT r.* FROM restaurant AS r 
                                 JOIN position AS p 
                                     ON r.id = p.restaurant_id
		                     GROUP BY r.id, rating 
                             ORDER BY rating DESC,
                                      AVG((p.price + r.minimal_delivery_cost) / p.weight)`)
	if err := repo.db.SelectContext(ctx, &restaurants, query); err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (repo *Repository) GetPositions(ctx context.Context, id int) ([]m.Position, error) {
	var positions []m.Position

	query := fmt.Sprintf(`SELECT name, description, price, weight, date_of_parsing FROM
                								(SELECT name,
                                                        description,
                                                        price,
                                                        weight,
                                                        date_of_parsing,
                                                        restaurant_id,
                                                        max(date_of_parsing) OVER (PARTITION BY restaurant_id, name) as max_date_of_parsing
                                                FROM position) AS p
                                 WHERE restaurant_id=$1 AND p.max_date_of_parsing=p.date_of_parsing`)

	if err := repo.db.SelectContext(ctx, &positions, query, id); err != nil {
		return nil, err
	}
	return positions, nil
}

func (repo *Repository) AddRestaurant(ctx context.Context, r *m.Restaurant) error {

	query := fmt.Sprintf(`INSERT INTO restaurant (name, 
                                                         slug, 
                                                         minimal_delivery_cost, 
                                                         rating)
                                 VALUES ($1, $2, $3, $4)
                                 ON CONFLICT(slug) DO UPDATE SET 
                                                                name = $1, 
                                                                minimal_delivery_cost = $3, 
                                                                rating = $4
                                 RETURNING id`)

	if err := repo.db.QueryRowxContext(ctx, query, r.Name, r.Slug, r.MinimalDeliveryCost, r.Rating).Scan(&r.Id); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) AddPosition(ctx context.Context, r *m.Restaurant, p []m.Position) error {

	query := fmt.Sprintf(`INSERT INTO position (name, 
                                                       price, 
                                                       description, 
                                                       weight, 
                                                       date_of_parsing, 
                                                       restaurant_id) 
                                 SELECT CAST($1 AS VARCHAR), $2, $3, $4, CAST($5 AS DATE), CAST($6 AS INTEGER)
                                 WHERE NOT EXISTS 
                                     (SELECT true FROM position WHERE 
                                                                     name=$1 AND 
                                                                     date_of_parsing=$5 AND 
                                                                     restaurant_id=$6)`)

	tx := repo.db.MustBegin()
	for _, position := range p {
		tx.MustExecContext(ctx, query, position.Name, position.Price, position.Description, position.Weight, time.Now().Format("2006-01-02"), &r.Id)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
