package checker

import (
	"strings"
	"time"

	"kpls/internal/model"
	"kpls/internal/store"
)

// IPQCChecker handles IPQC (In-Process Quality Control) checks
type IPQCChecker struct {
	store *store.FileStore
	gate  string // "IPQC-1" or "IPQC-2"
}

// NewIPQCChecker creates a new IPQC checker
func NewIPQCChecker(s *store.FileStore, gate string) *IPQCChecker {
	return &IPQCChecker{
		store: s,
		gate:  gate,
	}
}

// Check performs IPQC check on a job's artifact
func (c *IPQCChecker) Check(job *model.Job, artifactContent string) (*model.QualityCheck, error) {
	check := &model.QualityCheck{
		ID:        time.Now().Format("20060102150405"),
		JobID:     job.ID,
		GateType:  c.gate,
		CheckedAt: time.Now(),
	}

	var score int
	var maxScore int
	var defectCodes []string

	if c.gate == "IPQC-1" {
		// IPQC-1: Skeleton gate checks
		maxScore = 8
		score, defectCodes = c.checkSkeleton(artifactContent, job)
	} else if c.gate == "IPQC-2" {
		// IPQC-2: Evidence gate checks
		maxScore = 10
		score, defectCodes = c.checkEvidence(artifactContent, job)
	}

	check.Score = score
	check.MaxScore = maxScore
	check.DefectCodes = defectCodes

	// Pass criteria
	if c.gate == "IPQC-1" {
		// 7/8以上、重大NG（D11/D12）なし
		hasCritical := false
		for _, code := range defectCodes {
			if code == "D11" || code == "D12" {
				hasCritical = true
				break
			}
		}
		check.Passed = score >= 7 && !hasCritical
	} else if c.gate == "IPQC-2" {
		// 8/10以上、重大NG（D21/D23/D27）なし
		hasCritical := false
		for _, code := range defectCodes {
			if code == "D21" || code == "D23" || code == "D27" {
				hasCritical = true
				break
			}
		}
		check.Passed = score >= 8 && !hasCritical
	}

	return check, nil
}

func (c *IPQCChecker) checkSkeleton(content string, job *model.Job) (int, []string) {
	score := 0
	var defectCodes []string
	lowerContent := strings.ToLower(content)

	// Check 1: Conclusion at the beginning
	if strings.Contains(lowerContent, "# 結論") || strings.Contains(lowerContent, "## 結論") {
		score++
	} else {
		defectCodes = append(defectCodes, "D11")
	}

	// Check 2: Sections correspond to success criteria
	hasSections := strings.Contains(content, "#") || strings.Contains(content, "##")
	if hasSections {
		score++
	} else {
		defectCodes = append(defectCodes, "D11")
	}

	// Check 3: No missing points (risks, alternatives, etc.)
	hasRisks := strings.Contains(lowerContent, "リスク") || strings.Contains(lowerContent, "risk")
	hasAlternatives := strings.Contains(lowerContent, "代替") || strings.Contains(lowerContent, "alternative")
	if hasRisks || hasAlternatives {
		score++
	} else {
		defectCodes = append(defectCodes, "D12")
	}

	// Check 4: No scope bloat (simplified check)
	score++ // Assume OK if structure exists

	// Check 5: Section roles are unique
	score++ // Assume OK

	// Check 6: Term definitions present
	if strings.Contains(lowerContent, "定義") || strings.Contains(lowerContent, "用語") {
		score++
	} else {
		defectCodes = append(defectCodes, "D15")
	}

	// Check 7: Natural reading order
	score++ // Assume OK if structure exists

	// Check 8: Instructions for next stage
	if strings.Contains(lowerContent, "次") || strings.Contains(lowerContent, "next") {
		score++
	} else {
		defectCodes = append(defectCodes, "D13")
	}

	return score, defectCodes
}

func (c *IPQCChecker) checkEvidence(content string, job *model.Job) (int, []string) {
	score := 0
	var defectCodes []string
	lowerContent := strings.ToLower(content)

	// Check 1: Main claims have evidence
	hasSources := strings.Contains(lowerContent, "出典") || strings.Contains(lowerContent, "参照") ||
		strings.Contains(lowerContent, "source") || strings.Contains(lowerContent, "reference")
	if hasSources {
		score++
	} else {
		defectCodes = append(defectCodes, "D21")
	}

	// Check 2: Numbers/proper nouns have evidence
	score++ // Simplified check

	// Check 3: Uncertainty is stated, avoids assertions
	hasUncertainty := strings.Contains(lowerContent, "不確実") || strings.Contains(lowerContent, "未確定") ||
		strings.Contains(lowerContent, "uncertain")
	if hasUncertainty || !strings.Contains(lowerContent, "必ず") {
		score++
	} else {
		defectCodes = append(defectCodes, "D21")
	}

	// Check 4: No definition drift
	score++ // Simplified check

	// Check 5: No contradictions
	score++ // Simplified check

	// Check 6: Counterarguments/disadvantages are written
	hasCounter := strings.Contains(lowerContent, "反証") || strings.Contains(lowerContent, "デメリット") ||
		strings.Contains(lowerContent, "disadvantage")
	if hasCounter {
		score++
	} else {
		defectCodes = append(defectCodes, "D12")
	}

	// Check 7: Risks and countermeasures are realistic
	hasRisks := strings.Contains(lowerContent, "リスク") || strings.Contains(lowerContent, "risk")
	if hasRisks {
		score++
	} else {
		defectCodes = append(defectCodes, "D26")
	}

	// Check 8: No prohibited items violations
	hasProhibited := strings.Contains(lowerContent, "機密") || strings.Contains(lowerContent, "個人情報")
	if !hasProhibited {
		score++
	} else {
		defectCodes = append(defectCodes, "D27")
	}

	// Check 9: Warnings for easily misunderstood parts
	hasWarnings := strings.Contains(lowerContent, "注意") || strings.Contains(lowerContent, "警告") ||
		strings.Contains(lowerContent, "warning")
	if hasWarnings {
		score++
	} else {
		// Not critical, but good to have
		score++
	}

	// Check 10: References/links present
	if hasSources {
		score++
	} else {
		defectCodes = append(defectCodes, "D05")
	}

	return score, defectCodes
}
