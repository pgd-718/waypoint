package statetest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/stretchr/testify/require"
)

func init() {
	tests["events"] = []testFunc{
		TestEventListPagination,
	}
}

func TestEventListPagination(t *testing.T, factory Factory, rf RestartFactory) {
	ctx := context.Background()
	require := require.New(t)
	s := factory(t)
	defer s.Close()
	// a b c d e
	// f g h i j
	// k l m n o
	// p q r s t
	// u v w x y
	// z
	startChar := 'a'
	endChar := 'm'
	eventCount := endChar - startChar + 1
	var chars []string

	// Generate randomized events
	for char := startChar; char <= endChar; char++ {
		chars = append(chars, fmt.Sprintf("%c", char))
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(chars), func(i, j int) {
		chars[i], chars[j] = chars[j], chars[i]
	})
	for _, char := range chars {
		err := s.BuildPut(ctx, false, serverptypes.TestBuild(t, &pb.Build{
			Id:       char,
			Sequence: 1,
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "project",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))
		require.NoError(err)

		err = s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id: char,
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "project",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))
		require.NoError(err)

		err = s.ReleasePut(ctx, false, serverptypes.TestRelease(t, &pb.Release{
			Id: char,
			Application: &pb.Ref_Application{
				Application: "app",
				Project:     "project",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))
		require.NoError(err)

		//err = s.PipelineRunPut(ctx, &pb.PipelineRun{
		//	Id: char,
		//	Pipeline: &pb.Pipeline{
		//		Id:    "test",
		//		Name:  "test",
		//		Owner: nil,
		//		Steps: nil,
		//	},
		//})
		//require.NoError(err)

	}

	t.Run("EventList", func(t *testing.T) {
		t.Run(fmt.Sprintf("works with nil for compatibility and returns all %d results", eventCount), func(t *testing.T) {
			{
				resp, _, err := s.EventListBundles(ctx, &pb.UI_ListEventsRequest{
					Application: &pb.Ref_Application{
						Application: "app",
						Project:     "project",
					},
					Workspace:   &pb.Ref_Workspace{Workspace: "default"},
					Pagination:  &pb.PaginationRequest{
						PageSize:          5,
						NextPageToken:     "",
						PreviousPageToken: "",
					},
					//Sorting:     &pb.SortingRequest{OrderBy: []string{"name","created_at asc"}},
				})
				require.NoError(err)
				require.Len(resp, int(eventCount))

			}
		})

	})
}

//TODO: JUST test that the length returned from the eventlistbundle is correct since pagination is already
//thoroughly tested