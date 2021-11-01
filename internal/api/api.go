package api

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/bss-workplace-api/internal/repo"

	pb "github.com/ozonmp/bss-workplace-api/pkg/bss-workplace-api"
)

var (
	totalWorkplaceNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bss_workplace_api_workplace_not_found_total",
		Help: "Total number of workplaces that were not found",
	})
)

type workplaceAPI struct {
	pb.UnimplementedBssWorkplaceApiServiceServer
	repo repo.Repo
}

// NewWorkplaceAPI returns api of bss-workplace-api service
func NewWorkplaceAPI(r repo.Repo) pb.BssWorkplaceApiServiceServer {
	return &workplaceAPI{repo: r}
}

func (o *workplaceAPI) DescribeWorkplaceV1(
	ctx context.Context,
	req *pb.DescribeWorkplaceV1Request,
) (*pb.DescribeWorkplaceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeWorkplaceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	workplace, err := o.repo.DescribeWorkplace(ctx, req.WorkplaceId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeWorkplaceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if workplace == nil {
		log.Debug().Uint64("workplaceId", req.WorkplaceId).Msg("workplace not found")
		totalWorkplaceNotFound.Inc()

		return nil, status.Error(codes.NotFound, "workplace not found")
	}

	log.Debug().Msg("DescribeWorkplaceV1 - success")

	return &pb.DescribeWorkplaceV1Response{
		Value: &pb.Workplace{
			Id:  workplace.ID,
			Foo: workplace.Foo,
		},
	}, nil
}
