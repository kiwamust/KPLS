package model

import "time"

// QualityCheck represents a quality gate check result
type QualityCheck struct {
	ID           string    `json:"id"`
	JobID        string    `json:"job_id"`
	GateType     string    `json:"gate_type"` // IQC, IPQC-1, IPQC-2, FQC
	CheckedAt    time.Time `json:"checked_at"`
	Checker      string    `json:"checker"`
	Score        int       `json:"score"`
	MaxScore     int       `json:"max_score"`
	Passed       bool      `json:"passed"`
	DefectCodes  []string  `json:"defect_codes"`
	Notes        string    `json:"notes"`
	RejectReason string    `json:"reject_reason,omitempty"`
}

// Defect represents a defect record
type Defect struct {
	ID          string     `json:"id"`
	JobID       string     `json:"job_id"`
	Code        string     `json:"code"` // D01, D02, etc.
	GateType    string     `json:"gate_type"`
	Description string     `json:"description"`
	OccurredAt  time.Time  `json:"occurred_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	Resolution  string     `json:"resolution,omitempty"`
}
