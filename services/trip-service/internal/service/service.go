package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
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


func (s *TripUserService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmAPIResponse, error){
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
	
	var routeResp tripTypes.OsrmAPIResponse
	
	if err := json.Unmarshal(body, &routeResp); err != nil{
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	return &routeResp, nil
}

func (s *TripUserService) EstimatePackagesPriceWithRoute(route *tripTypes.OsrmAPIResponse) []*domain.RideFareModel{
	baseFares := getBaseFares()
	
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))
	
	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, route)
	}
	return estimatedFares

}

func (s *TripUserService) GenerateTripFares(ctx context.Context, ridefares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error){
	
	fares := make([]*domain.RideFareModel, len(ridefares))
	
	for i, f := range ridefares {
		id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			UserId: userID,
			ID: id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
		}
		
		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %w", err)
		}
		
		fares[i] = fare
		
	}
	return fares,nil
}

func estimateFareRoute(f *domain.RideFareModel, route *tripTypes.OsrmAPIResponse) *domain.RideFareModel {
	// distance
	// time
	// car price
	
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents
	
	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration
	
	
	// distance
	distanceFare := distanceKm + pricingCfg.PricePerUnitOfDistance
	
	//time
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	
	//car price
	totalPrice := carPackagePrice + distanceFare + timeFare
	
	return &domain.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug: f.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug: "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug: "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug: "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug: "luxury",
			TotalPriceInCents: 1000,
		},
	}
}