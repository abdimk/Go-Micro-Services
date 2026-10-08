package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type TripUserService struct{
	repo domain.TripRepository
}


func NewTripUserService(repo domain.TripRepository)(*TripUserService){
	return &TripUserService{repo: repo}
}


func (s *TripUserService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error){
	t := &domain.TripModel{
		ID: primitive.NewObjectID(),
		UserID: fare.UserId,
		Status: "pending",
		RideFare: fare,
	}
	
	return s.repo.CreateTrip(ctx, t)
}