package ptypes

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidatePipeline(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.Pipeline)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"no owner",
			func(v *pb.Pipeline) { v.Owner = nil },
			"Owner: cannot be blank",
		},

		{
			"project is blank",
			func(v *pb.Pipeline) {
				v.Owner = &pb.Pipeline_Project{
					Project: &pb.Ref_Project{Project: ""},
				}
			},
			"project: cannot be blank",
		},

		{
			"no steps",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{}
			},
			"steps: cannot be blank",
		},

		{
			"step name not set",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "",
					},
				}
			},
			"name: cannot be blank",
		},

		{
			"step name doesn't match key",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "bar",
					},
				}
			},
			`key "root" doesn't match`,
		},

		{
			"multiple root steps",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "root",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},

					"root2": {
						Name: "root2",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},
				}
			},
			`exactly one root`,
		},

		{
			"exec image required",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "root",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{},
						},
					},
				}
			},
			`image: cannot be blank`,
		},

		{
			"cycle",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "root",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},

					"A": {
						Name:      "A",
						DependsOn: []string{"B"},
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},

					"B": {
						Name:      "B",
						DependsOn: []string{"A"},
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},
				}
			},
			`one or more cycles`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			v := TestPipeline(t, nil)
			if f := tt.Modify; f != nil {
				f(v)
			}

			err := ValidatePipeline(v)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}

func TestUI_PipelineRunTreeFromJobs(t *testing.T) {
	cases := map[string]struct {
		Jobs []*pb.Job
		Tree *pb.UI_PipelineRunTreeNode
	}{
		"one queued exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-1",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State: pb.Job_QUEUED,
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State: pb.UI_PipelineRunTreeNode_QUEUED,
				Job: &pb.Ref_Job{
					Id: "job-1",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one running exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-1",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:   pb.Job_RUNNING,
					AckTime: quickTimestamp("2023-01-01T13:00:00Z"),
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:     pb.UI_PipelineRunTreeNode_RUNNING,
				StartTime: quickTimestamp("2023-01-01T13:00:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-1",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one successful exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-1",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:        pb.Job_SUCCESS,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_SUCCESS,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-1",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one errored exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-1",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:        pb.Job_ERROR,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_ERROR,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-1",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one cancelled exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-1",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:        pb.Job_ERROR,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CancelTime:   quickTimestamp("2023-01-01T13:08:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_CANCELLED,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-1",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one running step and one queued step": {
			Jobs: []*pb.Job{
				{
					Id: "job-1",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					AckTime: quickTimestamp("2023-01-01T13:00:00Z"),
					State:   pb.Job_RUNNING,
				},
				{
					Id: "job-2",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "bye",
								DependsOn: []string{"hello"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"bye"},
									},
								},
							},
						},
					},
					State: pb.Job_QUEUED,
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:     pb.UI_PipelineRunTreeNode_RUNNING,
				StartTime: quickTimestamp("2023-01-01T13:00:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-1",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode: pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{
						{
							Step: &pb.Pipeline_Step{
								Name: "bye",
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"bye"},
									},
								},
								DependsOn: []string{"hello"},
							},
							State: pb.UI_PipelineRunTreeNode_QUEUED,
							Job:   &pb.Ref_Job{Id: "job-2"},
							Children: &pb.UI_PipelineRunTreeNode_Children{
								Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
								Nodes: []*pb.UI_PipelineRunTreeNode{},
							},
						},
					},
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			result, err := UI_PipelineRunTreeFromJobs(tt.Jobs)

			require.NoError(err)

			if diff := cmp.Diff(tt.Tree, result, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected difference:\n%v", diff)
			}

		})
	}
}

// quickTimestamp parses an RFC3339-formatted string and returns the time it
// represents as a timestamppb.Timestamp.
//
// This is intended purely to make tests more readble and robust to daylight
// savings time etc.
func quickTimestamp(s string) *timestamppb.Timestamp {
	t, _ := time.Parse(time.RFC3339, s)
	return timestamppb.New(t)
}
