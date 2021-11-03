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

func (o *workplaceAPI) CreateWorkplaceV1(
	ctx context.Context,
	req *pb.CreateWorkplaceV1Request,
) (*pb.CreateWorkplaceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateWorkplaceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	workplaceId, err := o.repo.CreateWorkplace(ctx, req.Foo)
	if err != nil {
		log.Error().Err(err).Msg("CreateWorkplaceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("workplaceId", workplaceId).Msg("Workplace was created")

	return &pb.CreateWorkplaceV1Response{
		WorkplaceId: workplaceId,
	}, nil
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

func (o *workplaceAPI) ListWorkplacesV1(
	ctx context.Context,
	req *pb.ListWorkplacesV1Request,
) (*pb.ListWorkplacesV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListWorkplaceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	workplaces, err := o.repo.ListWorkplaces(ctx)
	if err != nil {
		log.Error().Err(err).Msg("ListWorkplacesV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if workplaces == nil {
		log.Debug().Msg("Workplaces not found")
		totalWorkplaceNotFound.Inc()

		return nil, status.Error(codes.NotFound, "workplaces not found")
	}

	log.Debug().Msg("ListWorkplacesV1 - success")

	return &pb.ListWorkplacesV1Response{
		Items: []*pb.Workplace{},
	}, nil
}

func (o *workplaceAPI) RemoveWorkplaceV1(
	ctx context.Context,
	req *pb.RemoveWorkplaceV1Request,
) (*pb.RemoveWorkplaceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveWorkplaceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err := o.repo.RemoveWorkplace(ctx, req.WorkplaceId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeWorkplaceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if ok == false {
		log.Debug().Uint64("workplaceId", req.WorkplaceId).Msg("workplace not removed")
		totalWorkplaceNotFound.Inc()

		return nil, status.Error(codes.NotFound, "workplace not removed")
	}

	log.Debug().Msg("RemoveWorkplaceV1 - success")

	return &pb.RemoveWorkplaceV1Response{
		Found: ok,
	}, nil
}