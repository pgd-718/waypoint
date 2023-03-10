package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

// UI_ListPipelines returns pipelines for a given project. While paginating is
// part of the request, this doesn't yet support pagination and will return
// everything every time.
func (s *Service) UI_ListPipelines(
	ctx context.Context,
	req *pb.UI_ListPipelinesRequest,
) (*pb.UI_ListPipelinesResponse, error) {
	log := hclog.FromContext(ctx)

	if err := serverptypes.ValidateUIListPipelinesRequest(req); err != nil {
		return nil, err
	}

	// Create uninitialized array of pipeline bundles
	var allPipelines []*pb.UI_PipelineBundle

	// Get list of all pipelines
	pipelineListResponse, err := s.state(ctx).PipelineList(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error listing piplines",
		)
	}

	// Create bundles
	for _, pipeline := range pipelineListResponse {
		ref := &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		}
		// Get total runs for the pipeline
		pipelineRunListResponse, err := s.state(ctx).PipelineRunList(ctx, ref)
		if err != nil {
			return nil, hcerr.Externalize(
				log,
				err,
				"failed to count pipeline runs",
			)
		}
		// Get the last run
		pipelineLastRun, err := s.state(ctx).PipelineRunGetLatest(ctx, pipeline.Id)
		if err != nil && status.Code(err) != codes.NotFound {
			return nil, hcerr.Externalize(
				log,
				err,
				"failed to find latest pipeline run",
			)
		}
		var lastRunBundle *pb.UI_PipelineRunBundle
		if pipelineLastRun != nil {
			lastRunBundle = &pb.UI_PipelineRunBundle{PipelineRun: pipelineLastRun}
		}

		pipelineBundle := &pb.UI_PipelineBundle{
			Pipeline:  pipeline,
			TotalRuns: uint64(len(pipelineRunListResponse)),
			LastRun:   lastRunBundle,
		}
		// Add pipeline bundle to uninitialized array
		allPipelines = append(allPipelines, pipelineBundle)

	}

	// Return the array
	return &pb.UI_ListPipelinesResponse{
		Pipelines:  allPipelines,
		Pagination: &pb.PaginationResponse{},
	}, nil
}
