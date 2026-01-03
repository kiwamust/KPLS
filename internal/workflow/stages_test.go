package workflow

import (
	"kpls/internal/model"
	"testing"
)

func TestGetNextStage(t *testing.T) {
	tests := []struct {
		current model.JobStatus
		next    model.JobStatus
	}{
		{model.StatusBacklog, model.StatusIQC},
		{model.StatusIQC, model.StatusSkeleton},
		{model.StatusSkeleton, model.StatusIPQC1},
		{model.StatusIPQC1, model.StatusDraft},
		{model.StatusDraft, model.StatusIPQC2},
		{model.StatusIPQC2, model.StatusPackaging},
		{model.StatusPackaging, model.StatusFQC},
		{model.StatusFQC, model.StatusDone},
		{model.StatusDone, model.StatusDone},
	}

	for _, tt := range tests {
		result := GetNextStage(tt.current)
		if result != tt.next {
			t.Errorf("GetNextStage(%s) = %s, want %s", tt.current, result, tt.next)
		}
	}
}

func TestIsGateStage(t *testing.T) {
	tests := []struct {
		stage model.JobStatus
		want  bool
	}{
		{model.StatusIQC, true},
		{model.StatusIPQC1, true},
		{model.StatusIPQC2, true},
		{model.StatusFQC, true},
		{model.StatusBacklog, false},
		{model.StatusSkeleton, false},
		{model.StatusDraft, false},
		{model.StatusPackaging, false},
		{model.StatusDone, false},
	}

	for _, tt := range tests {
		result := IsGateStage(tt.stage)
		if result != tt.want {
			t.Errorf("IsGateStage(%s) = %v, want %v", tt.stage, result, tt.want)
		}
	}
}
