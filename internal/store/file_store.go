package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"kpls/internal/model"
)

const (
	dataDir     = "data"
	jobsDir     = "data/jobs"
	defectsDir  = "data/defects"
)

// FileStore handles file-based storage
type FileStore struct {
	jobsDir    string
	defectsDir string
}

// NewFileStore creates a new file store
func NewFileStore() (*FileStore, error) {
	if err := os.MkdirAll(jobsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create jobs directory: %w", err)
	}
	if err := os.MkdirAll(defectsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create defects directory: %w", err)
	}
	return &FileStore{
		jobsDir:    jobsDir,
		defectsDir: defectsDir,
	}, nil
}

// SaveJob saves a job to file
func (s *FileStore) SaveJob(job *model.Job) error {
	path := filepath.Join(s.jobsDir, job.ID+".json")
	data, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadJob loads a job by ID
func (s *FileStore) LoadJob(id string) (*model.Job, error) {
	path := filepath.Join(s.jobsDir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read job file: %w", err)
	}
	var job model.Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}
	return &job, nil
}

// ListJobs returns all jobs
func (s *FileStore) ListJobs() ([]*model.Job, error) {
	entries, err := os.ReadDir(s.jobsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read jobs directory: %w", err)
	}
	
	var jobs []*model.Job
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := entry.Name()[:len(entry.Name())-5] // remove .json
		job, err := s.LoadJob(id)
		if err != nil {
			continue // skip invalid files
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// ListJobsByStatus returns jobs filtered by status
func (s *FileStore) ListJobsByStatus(status model.JobStatus) ([]*model.Job, error) {
	allJobs, err := s.ListJobs()
	if err != nil {
		return nil, err
	}
	var filtered []*model.Job
	for _, job := range allJobs {
		if job.Status == status {
			filtered = append(filtered, job)
		}
	}
	return filtered, nil
}

// SaveDefect saves a defect record
func (s *FileStore) SaveDefect(defect *model.Defect) error {
	path := filepath.Join(s.defectsDir, defect.ID+".json")
	data, err := json.MarshalIndent(defect, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal defect: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ListDefects returns all defects
func (s *FileStore) ListDefects() ([]*model.Defect, error) {
	entries, err := os.ReadDir(s.defectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read defects directory: %w", err)
	}
	
	var defects []*model.Defect
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(s.defectsDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var defect model.Defect
		if err := json.Unmarshal(data, &defect); err != nil {
			continue
		}
		defects = append(defects, &defect)
	}
	return defects, nil
}

// GenerateJobID generates a unique job ID
func GenerateJobID() string {
	return fmt.Sprintf("JT-%s-%03d", time.Now().Format("20060102"), time.Now().Unix()%1000)
}

// GenerateDefectID generates a unique defect ID
func GenerateDefectID() string {
	return fmt.Sprintf("DEF-%s-%03d", time.Now().Format("20060102"), time.Now().Unix()%1000)
}
