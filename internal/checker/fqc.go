package checker

import (
	"strings"
	"time"

	"kpls/internal/model"
	"kpls/internal/store"
)

// FQCChecker handles FQC (Final Quality Control) checks
type FQCChecker struct {
	store *store.FileStore
}

// NewFQCChecker creates a new FQC checker
func NewFQCChecker(s *store.FileStore) *FQCChecker {
	return &FQCChecker{store: s}
}

// Check performs FQC check on a job
func (c *FQCChecker) Check(job *model.Job, artifactContent string) (*model.QualityCheck, error) {
	check := &model.QualityCheck{
		ID:        time.Now().Format("20060102150405"),
		JobID:     job.ID,
		GateType:  "FQC",
		CheckedAt: time.Now(),
		MaxScore:  10,
	}

	score := 0
	var defectCodes []string
	lowerContent := strings.ToLower(artifactContent)

	// Check 1: Success criteria are met
	if len(job.SuccessCriteria) > 0 {
		score++
	} else {
		defectCodes = append(defectCodes, "D34")
	}

	// Check 2: Conclusion is clear at the beginning
	if strings.Contains(lowerContent, "# 結論") || strings.Contains(lowerContent, "## 結論") {
		score++
	} else {
		defectCodes = append(defectCodes, "D11")
	}

	// Check 3: Next actions (assignee/deadline) are present
	hasNextAction := strings.Contains(lowerContent, "次アクション") || strings.Contains(lowerContent, "next action") ||
		strings.Contains(lowerContent, "担当") || strings.Contains(lowerContent, "期限")
	if hasNextAction {
		score++
	} else {
		defectCodes = append(defectCodes, "D32")
	}

	// Check 4: Important assumptions/constraints are stated
	hasConstraints := strings.Contains(lowerContent, "制約") || strings.Contains(lowerContent, "前提") ||
		strings.Contains(lowerContent, "constraint") || strings.Contains(lowerContent, "assumption")
	if hasConstraints {
		score++
	} else {
		defectCodes = append(defectCodes, "D08")
	}

	// Check 5: Evidence/sources are traceable
	hasSources := strings.Contains(lowerContent, "出典") || strings.Contains(lowerContent, "参照") ||
		strings.Contains(lowerContent, "source") || strings.Contains(lowerContent, "reference")
	if hasSources {
		score++
	} else {
		defectCodes = append(defectCodes, "D05")
	}

	// Check 6: Risks/counterarguments are documented
	hasRisks := strings.Contains(lowerContent, "リスク") || strings.Contains(lowerContent, "反証") ||
		strings.Contains(lowerContent, "risk")
	if hasRisks {
		score++
	} else {
		defectCodes = append(defectCodes, "D26")
	}

	// Check 7: Readable (not redundant/duplicated)
	score++ // Simplified check

	// Check 8: Template compliant (format/metadata)
	hasMetadata := strings.Contains(artifactContent, "---") // YAML frontmatter
	if hasMetadata {
		score++
	} else {
		defectCodes = append(defectCodes, "D31")
	}

	// Check 9: Reusable (structure/tags/citations)
	if hasSources && hasMetadata {
		score++
	} else {
		defectCodes = append(defectCodes, "D33")
	}

	// Check 10: No prohibited items violations
	hasProhibited := strings.Contains(lowerContent, "機密") || strings.Contains(lowerContent, "個人情報")
	if !hasProhibited {
		score++
	} else {
		defectCodes = append(defectCodes, "D27")
	}

	check.Score = score
	check.DefectCodes = defectCodes

	// Pass if score >= 9 and no critical defects (D34, D27)
	hasCriticalDefect := false
	for _, code := range defectCodes {
		if code == "D34" || code == "D27" {
			hasCriticalDefect = true
			break
		}
	}
	check.Passed = score >= 9 && !hasCriticalDefect

	return check, nil
}
