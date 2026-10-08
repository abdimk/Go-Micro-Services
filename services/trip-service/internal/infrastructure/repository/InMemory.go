package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)



type InMemoryRepository struct {
	trips map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel

}


func NewInMemoryRepository()(*InMemoryRepository){
	return &InMemoryRepository{
		trips: make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}


func (r *InMemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error){
	r.trips[trip.ID.Hex()] = trip
	
	return trip, nil
}
