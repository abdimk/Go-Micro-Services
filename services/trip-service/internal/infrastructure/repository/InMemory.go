package repository

import (
	"context"
	"fmt"
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


func (r *InMemoryRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error){
	fare, exist := r.rideFares[id]
	
	if !exist {
		return nil, fmt.Errorf("fare does not exist with ID: %s", id)
	}
	
	return fare,nil
}


func (r *InMemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error){
	r.trips[trip.ID.Hex()] = trip
	
	return trip, nil
}

func (r *InMemoryRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.rideFares[f.ID.Hex()] = f
	return nil
}