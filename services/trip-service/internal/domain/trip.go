package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserId            string
	PackageSlug       string   //ex: van, luxury, sedan
	TotalPriceInCents float64
}

type TripModel struct {
	ID primitive.ObjectID
	UserID string
	Status string
	RideFare RideFareModel
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip TripModel) (*TripModel, error)
	
}


type TripService interface {
	CreateTrip(ctx context.Context, fare RideFareModel) (*RideFareModel, error)
}