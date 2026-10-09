package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const baseURL= "https://router.project-osrm.org"

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


func (s *TripUserService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmAPIResponse, error){
	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		baseURL,
		pickup.Longitude,pickup.Latitude,
		destination.Longitude,destination.Latitude,
	)
	
	resp, err := http.Get(url)
	
	if err != nil{
		return nil,fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}
	
	var routeResp types.OsrmAPIResponse
	
	if err := json.Unmarshal(body, &routeResp); err != nil{
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	return &routeResp, nil
	

}