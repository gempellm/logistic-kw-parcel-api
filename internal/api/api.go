package api

import (
	"context"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gempellm/logistic-parcel-api/internal/repo"

	pb "github.com/gempellm/logistic-parcel-api/pkg/logistic_parcel_api"
)

var (
	totalparcelNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistic_parcel_api_parcel_not_found_total",
		Help: "Total number of parcels that were not found",
	})
)

type parcelAPI struct {
	pb.UnimplementedLogisticParcelApiServiceServer
	repo repo.Repo
}

// NewparcelAPI returns api of logistic-parcel-api service
func NewparcelAPI(r repo.Repo) pb.LogisticParcelApiServiceServer {
	return &parcelAPI{repo: r}
}

func (o *parcelAPI) DescribeParcelV1(
	ctx context.Context,
	req *pb.DescribeParcelV1Request,
) (*pb.DescribeParcelV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeParcelV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parcel, err := o.repo.DescribeParcel(ctx, req.ParcelId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeParcelV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if parcel == nil {
		log.Debug().Uint64("parcelId", req.ParcelId).Msg("parcel not found")
		totalparcelNotFound.Inc()

		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	log.Debug().Msg("DescribeParcelV1 - success")

	return &pb.DescribeParcelV1Response{
		Value: &pb.Parcel{
			Id: parcel.ID,
		},
	}, nil
}

func (o *parcelAPI) DescribeParcel(ctx context.Context, req *pb.DescribeParcelRequest) (*pb.DescribeParcelResponse, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeParcel - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parcel, err := o.repo.DescribeParcel(ctx, req.ParcelId)
	if err != nil {
		if errors.Is(err, repo.ErrParcelNotFound) {
			return nil, status.Error(codes.NotFound, "parcel not found")
		}

		log.Error().Err(err).Msg("DescribeParcel -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if parcel == nil {
		log.Debug().Uint64("parcelId", req.ParcelId).Msg("parcel not found")
		totalparcelNotFound.Inc()

		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	log.Debug().Msg("DescribeParcel - success")

	return &pb.DescribeParcelResponse{
		Value: &pb.Parcel{
			Id:   parcel.ID,
			Name: parcel.Name,
		},
	}, nil
}
