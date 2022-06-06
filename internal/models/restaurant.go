package models

import (
	"strconv"
	"strings"
	"time"
)

type ParseRequest struct {
	Latitude  float64 `schema:"latitude"  validate:"latitude"`
	Longitude float64 `schema:"longitude" validate:"longitude"`
	Workers   int     `schema:"workers"   validate:"min=1,worker"`
}

type Restaurant struct {
	Id 					int       `json:"id"                  db:"id"`
	Name                string    `json:"name"                db:"name"`
	Slug                string    `json:"slug"                db:"slug"`
	Rating              float64   `json:"rating"              db:"rating"`
	MinimalDeliveryCost float64   `json:"minimalDeliveryCost" db:"minimal_delivery_cost"`
}

type Position struct {
	Id            int         `json:"-"               db:"id"`
	RestaurantId  int         `json:"-"               db:"restaurant_id"`
	Name          string      `json:"name"            db:"name"`
	Description   string      `json:"description"     db:"description"`
	Price         int         `json:"price"           db:"price"`
	Weight	      interface{} `json:"weight"          db:"weight"`
	DateOfParsing time.Time   `json:"DateOfParsing"   db:"date_of_parsing"`
}

type RestaurantPositions struct {
	Restaurant Restaurant
	Positions []Position
}

type YandexPositionsResponse struct {
	Payload struct {
		Categories []struct {
			Items      []Position `json:"items"`
		} `json:"categories"`
	} `json:"payload"`
}

type YandexRestaurantsResponse struct {
	Payload struct {
		FoundPlaces []struct {
			Restaurant `json:"place"`
		} `json:"foundPlaces"`
	} `json:"payload"`
}

func (p *Position) MatchFoodParameter() bool {
	if strings.Contains(strings.ToLower(p.Name), "филадельфи") &&
		strings.Contains(strings.ToLower(p.Name), "ролл") &&
		strings.Contains(strings.ToLower(p.Description), "лосос") {

		return true
	}
	return false
}

func (p *Position) ConvertWeightToInt() error {
	switch p.Weight.(type) {
	case string:
		weight, err := strconv.Atoi(strings.TrimSuffix(p.Weight.(string), "\u00A0г")); if err != nil {
		return err
	}
		p.Weight = weight
	}

	return nil
}
