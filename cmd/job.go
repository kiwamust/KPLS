package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"kpls/internal/model"
	"kpls/internal/store"
	"kpls/internal/workflow"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manage jobs",
	Long:  `Create, list, show, advance, and reject jobs in the workflow.`,
}

var jobCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new job",
	Long:  `Create a new job ticket with required information.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		job := &model.Job{
			ID:               store.GenerateJobID(),
			Title:            "",
			Owner:            "",
			Created:          time.Now(),
			Due:              time.Now().Add(7 * 24 * time.Hour), // default 7 days
			Priority:         model.PriorityP2,
			Status:           model.StatusBacklog,
			OutputType:       model.OutputType1Pager,
			Confidentiality:  model.ConfInternal,
			Audience:         []string{},
			SuccessCriteria:  []string{},
			Constraints:      []string{},
			Assumptions:      []string{},
			ScopeIn:          []string{},
			ScopeOut:         []string{},
			Materials:        []model.Material{},
			OpenQuestions:    []string{},
			DefinitionOfDone: []string{},
			StageRuns:        []model.StageRun{},
			QualityChecks:    []model.QualityCheck{},
			Artifacts:        []model.Artifact{},
		}

		// Interactive input (simplified - in real implementation, use survey or similar)
		fmt.Print("Title: ")
		if _, err := fmt.Scanln(&job.Title); err != nil {
			return fmt.Errorf("failed to read title: %w", err)
		}
		if job.Title == "" {
			return fmt.Errorf("title is required")
		}

		fmt.Print("Owner: ")
		if _, err := fmt.Scanln(&job.Owner); err != nil {
			return fmt.Errorf("failed to read owner: %w", err)
		}

		// Initialize first stage run
		job.StageRuns = append(job.StageRuns, model.StageRun{
			Stage:     model.StatusBacklog,
			StartedAt: time.Now(),
		})

		if err := s.SaveJob(job); err != nil {
			return err
		}

		fmt.Printf("Created job: %s\n", job.ID)
		return nil
	},
}

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all jobs",
	Long:  `List all jobs in kanban format grouped by stage.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		jobs, err := s.ListJobs()
		if err != nil {
			return err
		}

		// Group by status
		byStatus := make(map[model.JobStatus][]*model.Job)
		for _, job := range jobs {
			byStatus[job.Status] = append(byStatus[job.Status], job)
		}

		// Display kanban
		stages := workflow.GetAllStages()
		for _, stage := range stages {
			stageJobs := byStatus[stage]
			if len(stageJobs) == 0 {
				continue
			}

			fmt.Printf("\n=== %s (%d) ===\n", stage, len(stageJobs))
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Title", "Owner", "Priority", "Due"})
			table.SetBorder(false)

			for _, job := range stageJobs {
				table.Append([]string{
					job.ID,
					job.Title,
					job.Owner,
					string(job.Priority),
					job.Due.Format("2006-01-02"),
				})
			}
			table.Render()
		}

		return nil
	},
}

var jobShowCmd = &cobra.Command{
	Use:   "show [job-id]",
	Short: "Show job details",
	Long:  `Show detailed information about a job.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		job, err := s.LoadJob(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Job ID: %s\n", job.ID)
		fmt.Printf("Title: %s\n", job.Title)
		fmt.Printf("Status: %s\n", job.Status)
		fmt.Printf("Owner: %s\n", job.Owner)
		fmt.Printf("Priority: %s\n", job.Priority)
		fmt.Printf("Output Type: %s\n", job.OutputType)
		fmt.Printf("Created: %s\n", job.Created.Format("2006-01-02 15:04"))
		fmt.Printf("Due: %s\n", job.Due.Format("2006-01-02"))

		if len(job.Materials) > 0 {
			fmt.Println("\nMaterials:")
			for i, m := range job.Materials {
				fmt.Printf("  %d. [%s] %s\n", i+1, m.Kind, m.Ref)
			}
		}

		if len(job.QualityChecks) > 0 {
			fmt.Println("\nQuality Checks:")
			for _, qc := range job.QualityChecks {
				status := "PASS"
				if !qc.Passed {
					status = "FAIL"
				}
				fmt.Printf("  %s: %s - Score: %d/%d - %s\n",
					qc.GateType, status, qc.Score, qc.MaxScore, qc.CheckedAt.Format("2006-01-02 15:04"))
				if len(qc.DefectCodes) > 0 {
					fmt.Printf("    Defects: %v\n", qc.DefectCodes)
				}
			}
		}

		return nil
	},
}

var jobAdvanceCmd = &cobra.Command{
	Use:   "advance [job-id]",
	Short: "Advance job to next stage",
	Long:  `Move a job to the next stage in the workflow.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		job, err := s.LoadJob(args[0])
		if err != nil {
			return err
		}

		if err := workflow.AdvanceJob(job); err != nil {
			return err
		}

		if err := s.SaveJob(job); err != nil {
			return err
		}

		fmt.Printf("Job %s advanced to %s\n", job.ID, job.Status)
		return nil
	},
}

var jobRejectCmd = &cobra.Command{
	Use:   "reject [job-id]",
	Short: "Reject job and move back",
	Long:  `Reject a job and move it back to the previous stage with defect codes.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		job, err := s.LoadJob(args[0])
		if err != nil {
			return err
		}

		reason, err := cmd.Flags().GetString("reason")
		if err != nil {
			return fmt.Errorf("failed to get reason flag: %w", err)
		}
		if reason == "" {
			return fmt.Errorf("rejection reason is required (--reason)")
		}

		defects, err := cmd.Flags().GetStringSlice("defects")
		if err != nil {
			return fmt.Errorf("failed to get defects flag: %w", err)
		}
		if len(defects) == 0 {
			return fmt.Errorf("at least one defect code is required (--defects)")
		}

		if err := workflow.RejectJob(job, reason, defects); err != nil {
			return err
		}

		if err := s.SaveJob(job); err != nil {
			return err
		}

		fmt.Printf("Job %s rejected and moved back to %s\n", job.ID, job.Status)
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobCreateCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobShowCmd)
	jobCmd.AddCommand(jobAdvanceCmd)
	jobCmd.AddCommand(jobRejectCmd)

	jobRejectCmd.Flags().String("reason", "", "Rejection reason")
	jobRejectCmd.Flags().StringSlice("defects", []string{}, "Defect codes (e.g., D01,D02)")
}
