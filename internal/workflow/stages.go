package workflow

import "kpls/internal/model"

// GetNextStage returns the next stage in the workflow
func GetNextStage(current model.JobStatus) model.JobStatus {
	switch current {
	case model.StatusBacklog:
		return model.StatusIQC
	case model.StatusIQC:
		return model.StatusSkeleton
	case model.StatusSkeleton:
		return model.StatusIPQC1
	case model.StatusIPQC1:
		return model.StatusDraft
	case model.StatusDraft:
		return model.StatusIPQC2
	case model.StatusIPQC2:
		return model.StatusPackaging
	case model.StatusPackaging:
		return model.StatusFQC
	case model.StatusFQC:
		return model.StatusDone
	default:
		return current // no next stage
	}
}

// GetPreviousStage returns the previous stage (for rejection)
func GetPreviousStage(current model.JobStatus) model.JobStatus {
	switch current {
	case model.StatusIQC:
		return model.StatusBacklog
	case model.StatusSkeleton:
		return model.StatusIQC
	case model.StatusIPQC1:
		return model.StatusSkeleton
	case model.StatusDraft:
		return model.StatusIPQC1
	case model.StatusIPQC2:
		return model.StatusDraft
	case model.StatusPackaging:
		return model.StatusIPQC2
	case model.StatusFQC:
		return model.StatusPackaging
	default:
		return current // no previous stage
	}
}

// CanAdvance checks if job can advance to next stage
func CanAdvance(status model.JobStatus) bool {
	return status != model.StatusDone
}

// IsGateStage checks if the stage is a quality gate
func IsGateStage(status model.JobStatus) bool {
	return status == model.StatusIQC ||
		status == model.StatusIPQC1 ||
		status == model.StatusIPQC2 ||
		status == model.StatusFQC
}

// GetAllStages returns all workflow stages in order
func GetAllStages() []model.JobStatus {
	return []model.JobStatus{
		model.StatusBacklog,
		model.StatusIQC,
		model.StatusSkeleton,
		model.StatusIPQC1,
		model.StatusDraft,
		model.StatusIPQC2,
		model.StatusPackaging,
		model.StatusFQC,
		model.StatusDone,
	}
}
