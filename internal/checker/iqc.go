package checker

import (
	"time"

	"kpls/internal/model"
	"kpls/internal/store"
)

// IQCChecker handles IQC (Incoming Quality Control) checks
type IQCChecker struct {
	store *store.FileStore
}

// NewIQCChecker creates a new IQC checker
func NewIQCChecker(s *store.FileStore) *IQCChecker {
	return &IQCChecker{store: s}
}

// Check performs IQC check on a job
func (c *IQCChecker) Check(job *model.Job) (*model.QualityCheck, error) {
	check := &model.QualityCheck{
		ID:        time.Now().Format("20060102150405"),
		JobID:     job.ID,
		GateType:  "IQC",
		CheckedAt: time.Now(),
		MaxScore:  10,
	}

	score := 0
	var defectCodes []string

	// Check 1: Purpose is clear
	if len(job.SuccessCriteria) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D01")
	}

	// Check 2: Audience is clear
	if len(job.Audience) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D02")
	}

	// Check 3: Success criteria exists
	if len(job.SuccessCriteria) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D03")
	}

	// Check 4: Scope in/out exists
	if len(job.ScopeIn) > 0 || len(job.ScopeOut) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D07")
	}

	// Check 5: Constraints are specified
	if len(job.Constraints) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D08")
	}

	// Check 6: Required materials are present
	if len(job.Materials) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D04")
	}

	// Check 7: Sources are traceable
	hasTraceableSources := false
	for _, m := range job.Materials {
		if m.Kind == model.MaterialURL || m.Kind == model.MaterialFile {
			if m.Ref != "" {
				hasTraceableSources = true
				break
			}
		}
	}
	if hasTraceableSources {
		score++
	} else {
		defectCodes = append(defectCodes, "D05")
	}

	// Check 8: Freshness is verifiable
	hasFreshness := false
	for _, m := range job.Materials {
		if m.Freshness != "" {
			hasFreshness = true
			break
		}
	}
	if hasFreshness {
		score++
	} else {
		defectCodes = append(defectCodes, "D06")
	}

	// Check 9: Terms/assumptions are present
	if len(job.Assumptions) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D15")
	}

	// Check 10: Open questions are listed
	if len(job.OpenQuestions) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D04")
	}

	check.Score = score
	check.DefectCodes = defectCodes
	
	// Pass if score >= 8 and no critical defects (D05, D08)
	hasCriticalDefect := false
	for _, code := range defectCodes {
		if code == "D05" || code == "D08" {
			hasCriticalDefect = true
			break
		}
	}
	check.Passed = score >= 8 && !hasCriticalDefect

	return check, nil
}
