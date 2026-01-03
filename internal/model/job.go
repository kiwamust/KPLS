package model

import "time"

// JobStatus represents the workflow stage
type JobStatus string

const (
	StatusBacklog   JobStatus = "Backlog"
	StatusIQC       JobStatus = "IQC"
	StatusSkeleton  JobStatus = "Skeleton"
	StatusIPQC1     JobStatus = "IPQC-1"
	StatusDraft     JobStatus = "Draft"
	StatusIPQC2     JobStatus = "IPQC-2"
	StatusPackaging JobStatus = "Packaging"
	StatusFQC       JobStatus = "FQC"
	StatusDone      JobStatus = "Done"
)

// Priority represents job priority
type Priority string

const (
	PriorityP1 Priority = "P1"
	PriorityP2 Priority = "P2"
	PriorityP3 Priority = "P3"
)

// OutputType represents the output template type
type OutputType string

const (
	OutputType1Pager    OutputType = "1pager"
	OutputTypeComparison OutputType = "comparison"
	OutputTypePRD       OutputType = "prd"
)

// Confidentiality represents confidentiality level
type Confidentiality string

const (
	ConfPublic       Confidentiality = "Public"
	ConfInternal     Confidentiality = "Internal"
	Confidential     Confidentiality = "Confidential"
)

// MaterialKind represents the type of material
type MaterialKind string

const (
	MaterialURL  MaterialKind = "url"
	MaterialFile MaterialKind = "file"
	MaterialNote MaterialKind = "note"
)

// Reliability represents material reliability
type Reliability string

const (
	RelHigh    Reliability = "High"
	RelMid     Reliability = "Mid"
	RelLow     Reliability = "Low"
	RelUnknown Reliability = "Unknown"
)

// Material represents input source
type Material struct {
	Kind       MaterialKind `json:"kind"`
	Ref        string      `json:"ref"`
	Freshness  string      `json:"freshness"`
	Reliability Reliability `json:"reliability"`
}

// Job represents a job ticket
type Job struct {
	ID              string          `json:"id"`
	Title           string          `json:"title"`
	Owner           string          `json:"owner"`
	Created         time.Time       `json:"created"`
	Due             time.Time       `json:"due"`
	Priority        Priority        `json:"priority"`
	Status          JobStatus       `json:"status"`
	OutputType      OutputType      `json:"output_type"`
	Confidentiality Confidentiality `json:"confidentiality"`
	Audience        []string        `json:"audience"`
	SuccessCriteria []string        `json:"success_criteria"`
	Constraints     []string        `json:"constraints"`
	Assumptions     []string        `json:"assumptions"`
	ScopeIn         []string        `json:"scope_in"`
	ScopeOut        []string        `json:"scope_out"`
	Materials       []Material      `json:"materials"`
	OpenQuestions   []string        `json:"open_questions"`
	DefinitionOfDone []string        `json:"definition_of_done"`
	
	// Workflow tracking
	StageRuns       []StageRun       `json:"stage_runs"`
	QualityChecks   []QualityCheck  `json:"quality_checks"`
	Artifacts       []Artifact      `json:"artifacts"`
}

// StageRun represents execution log of a stage
type StageRun struct {
	Stage     JobStatus  `json:"stage"`
	StartedAt time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes     string     `json:"notes"`
}

// Artifact represents generated output
type Artifact struct {
	Type      string    `json:"type"` // skeleton, draft, final
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Version   string    `json:"version"`
}
