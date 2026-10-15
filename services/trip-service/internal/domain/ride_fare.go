package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	pb "ride-sharing/shared/proto/trip"
)
	




type RideFareModel struct {
	ID               primitive.ObjectID
	UserId            string
	PackageSlug       string   //ex: van, luxury, sedan
	TotalPriceInCents float64
}


func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id: r.ID.Hex(),
		UserID: r.UserId,
		PackageSlug: r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	
	return protoFares
}