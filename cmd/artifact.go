package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"kpls/internal/model"
	"kpls/internal/store"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage artifacts",
	Long:  `Add and list artifacts for jobs.`,
}

var artifactAddCmd = &cobra.Command{
	Use:   "add <job-id>",
	Short: "Add an artifact to a job",
	Long:  `Add an artifact file to a job with automatic naming.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID := args[0]
		artifactType, err := cmd.Flags().GetString("type")
		if err != nil {
			return fmt.Errorf("failed to get type flag: %w", err)
		}
		if artifactType == "" {
			return fmt.Errorf("--type is required (skeleton, draft, final, etc.)")
		}

		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("failed to get file flag: %w", err)
		}
		if filePath == "" {
			return fmt.Errorf("--file is required")
		}

		// Load job
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		job, err := s.LoadJob(jobID)
		if err != nil {
			return fmt.Errorf("failed to load job: %w", err)
		}

		// Read source file
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read artifact file: %w", err)
		}

		// Determine version number
		version := 1
		for _, art := range job.Artifacts {
			if art.Type == artifactType {
				// Extract version number from existing artifacts
				if strings.HasPrefix(art.Content, jobID+"-"+artifactType+"-v") {
					// Try to parse version
					parts := strings.Split(art.Content, "-v")
					if len(parts) > 1 {
						var v int
						if _, err := fmt.Sscanf(parts[1], "%d", &v); err == nil {
							if v >= version {
								version = v + 1
							}
						}
					}
				}
			}
		}

		// Generate artifact filename
		artifactsDir := "data/artifacts"
		if err := os.MkdirAll(artifactsDir, 0755); err != nil {
			return fmt.Errorf("failed to create artifacts directory: %w", err)
		}

		artifactFileName := fmt.Sprintf("%s-%s-v%d.md", jobID, artifactType, version)
		artifactPath := filepath.Join(artifactsDir, artifactFileName)

		// Write artifact file
		if err := os.WriteFile(artifactPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write artifact file: %w", err)
		}

		// Add to job artifacts
		artifact := model.Artifact{
			Type:      artifactType,
			Content:   artifactFileName, // Store relative path
			CreatedAt: time.Now(),
			Version:   fmt.Sprintf("v%d", version),
		}
		job.Artifacts = append(job.Artifacts, artifact)

		// Save job
		if err := s.SaveJob(job); err != nil {
			return fmt.Errorf("failed to save job: %w", err)
		}

		fmt.Printf("Artifact added: %s\n", artifactFileName)
		return nil
	},
}

var artifactListCmd = &cobra.Command{
	Use:   "list <job-id>",
	Short: "List artifacts for a job",
	Long:  `List all artifacts associated with a job.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		jobID := args[0]

		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		job, err := s.LoadJob(jobID)
		if err != nil {
			return fmt.Errorf("failed to load job: %w", err)
		}

		if len(job.Artifacts) == 0 {
			fmt.Printf("No artifacts found for job %s\n", jobID)
			return nil
		}

		fmt.Printf("Artifacts for job %s:\n\n", jobID)
		for i, art := range job.Artifacts {
			fmt.Printf("%d. Type: %s\n", i+1, art.Type)
			fmt.Printf("   File: %s\n", art.Content)
			fmt.Printf("   Version: %s\n", art.Version)
			fmt.Printf("   Created: %s\n", art.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	artifactCmd.AddCommand(artifactAddCmd)
	artifactCmd.AddCommand(artifactListCmd)

	artifactAddCmd.Flags().String("type", "", "Artifact type (skeleton, draft, final, etc.)")
	artifactAddCmd.Flags().String("file", "", "Path to artifact file")
}
