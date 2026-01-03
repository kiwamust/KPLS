package store

import (
	"os"
	"testing"
	"time"

	"kpls/internal/model"
)

func TestFileStore_SaveAndLoadJob(t *testing.T) {
	// Setup
	store, err := NewFileStore()
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer os.RemoveAll("data")

	// Create test job
	job := &model.Job{
		ID:              "TEST-001",
		Title:           "Test Job",
		Owner:           "Test Owner",
		Created:         time.Now(),
		Due:             time.Now().Add(7 * 24 * time.Hour),
		Priority:        model.PriorityP1,
		Status:          model.StatusBacklog,
		OutputType:      model.OutputType1Pager,
		Confidentiality: model.ConfInternal,
		Audience:        []string{"DecisionMaker"},
		SuccessCriteria: []string{"Test criteria"},
		Constraints:     []string{"Test constraint"},
		Assumptions:     []string{"Test assumption"},
		ScopeIn:         []string{"In scope"},
		ScopeOut:        []string{"Out of scope"},
		Materials:       []model.Material{},
		OpenQuestions:   []string{},
		DefinitionOfDone: []string{},
		StageRuns:       []model.StageRun{},
		QualityChecks:   []model.QualityCheck{},
		Artifacts:       []model.Artifact{},
	}

	// Save
	err = store.SaveJob(job)
	if err != nil {
		t.Fatalf("Failed to save job: %v", err)
	}

	// Load
	loaded, err := store.LoadJob("TEST-001")
	if err != nil {
		t.Fatalf("Failed to load job: %v", err)
	}

	// Verify
	if loaded.ID != job.ID {
		t.Errorf("Expected ID %s, got %s", job.ID, loaded.ID)
	}
	if loaded.Title != job.Title {
		t.Errorf("Expected Title %s, got %s", job.Title, loaded.Title)
	}
}

func TestFileStore_ListJobsByStatus(t *testing.T) {
	// Setup
	store, err := NewFileStore()
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer os.RemoveAll("data")

	// Create test jobs
	job1 := &model.Job{
		ID:     "TEST-001",
		Title:  "Test Job 1",
		Status: model.StatusBacklog,
	}
	job2 := &model.Job{
		ID:     "TEST-002",
		Title:  "Test Job 2",
		Status: model.StatusIQC,
	}

	store.SaveJob(job1)
	store.SaveJob(job2)

	// List by status
	backlogJobs, err := store.ListJobsByStatus(model.StatusBacklog)
	if err != nil {
		t.Fatalf("Failed to list jobs: %v", err)
	}

	if len(backlogJobs) < 1 {
		t.Error("Expected at least one backlog job")
	}
}
