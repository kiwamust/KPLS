package workflow

import (
	"time"

	"kpls/internal/model"
)

// AdvanceJob advances a job to the next stage
func AdvanceJob(job *model.Job) error {
	if !CanAdvance(job.Status) {
		return nil // already at final stage
	}

	// Complete current stage run
	if len(job.StageRuns) > 0 {
		lastRun := &job.StageRuns[len(job.StageRuns)-1]
		if lastRun.CompletedAt == nil {
			now := time.Now()
			lastRun.CompletedAt = &now
		}
	}

	// Move to next stage
	nextStage := GetNextStage(job.Status)
	job.Status = nextStage

	// Create new stage run
	job.StageRuns = append(job.StageRuns, model.StageRun{
		Stage:     nextStage,
		StartedAt: time.Now(),
	})

	return nil
}

// RejectJob rejects a job and moves it back to previous stage
func RejectJob(job *model.Job, reason string, defectCodes []string) error {
	previousStage := GetPreviousStage(job.Status)
	if previousStage == job.Status {
		return nil // cannot go back further
	}

	// Record rejection in quality check
	check := model.QualityCheck{
		ID:           time.Now().Format("20060102150405"),
		JobID:        job.ID,
		GateType:     string(job.Status),
		CheckedAt:    time.Now(),
		Passed:       false,
		DefectCodes:  defectCodes,
		RejectReason: reason,
	}
	job.QualityChecks = append(job.QualityChecks, check)

	// Move back to previous stage
	job.Status = previousStage

	// Create new stage run for the previous stage
	job.StageRuns = append(job.StageRuns, model.StageRun{
		Stage:     previousStage,
		StartedAt: time.Now(),
		Notes:     "Rejected: " + reason,
	})

	return nil
}
