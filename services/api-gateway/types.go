package main

import (	
	"ride-sharing/shared/types"
	pb "ride-sharing/shared/proto/trip"
)


type previewTripRequest struct {
	UserID       string          `json:"userID"`  
	Pickup       types.Coordinate `json:"pickup"`
	Destination  types.Coordinate  `json:"destination"`
}

func (p *previewTripRequest) ToProto() (*pb.PreviewTripRequest){
	return &pb.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: &pb.Coordinate{
			Latitude: p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		
		EndLocation: &pb.Coordinate{
			Latitude: p.Destination.Latitude,
			Longitude: p.Destination.Longitude,
		},
	}
}