package cmd

import (
	"encoding/json"
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
	RunE: func(cmd *cobra.Command, _ []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		var job *model.Job

		// Check for --from-file option
		fromFile, err := cmd.Flags().GetString("from-file")
		if err != nil {
			return fmt.Errorf("failed to get from-file flag: %w", err)
		}
		if fromFile != "" {
			// Load from JSON file
			data, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("failed to read job file: %w", err)
			}
			job = &model.Job{}
			if err := json.Unmarshal(data, job); err != nil {
				return fmt.Errorf("failed to parse job JSON: %w", err)
			}
			// Override ID and timestamps
			job.ID = store.GenerateJobID()
			job.Created = time.Now()
			if job.Due.IsZero() {
				job.Due = time.Now().Add(7 * 24 * time.Hour)
			}
		} else {
			// Check for flag-based creation
			title, err := cmd.Flags().GetString("title")
			if err != nil {
				return fmt.Errorf("failed to get title flag: %w", err)
			}
			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return fmt.Errorf("failed to get owner flag: %w", err)
			}
			outputType, err := cmd.Flags().GetString("type")
			if err != nil {
				return fmt.Errorf("failed to get type flag: %w", err)
			}

			if title != "" && owner != "" {
				// Create from flags
				job = &model.Job{
					ID:               store.GenerateJobID(),
					Title:            title,
					Owner:            owner,
					Created:          time.Now(),
					Due:              time.Now().Add(7 * 24 * time.Hour),
					Priority:         model.PriorityP2,
					Status:           model.StatusBacklog,
					OutputType:       model.OutputType(outputType),
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
			} else {
				// Interactive mode (original behavior)
				job = &model.Job{
					ID:               store.GenerateJobID(),
					Title:            "",
					Owner:            "",
					Created:          time.Now(),
					Due:              time.Now().Add(7 * 24 * time.Hour),
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
			}
		}

		// Validate required fields
		if job.Title == "" {
			return fmt.Errorf("title is required")
		}
		if job.Owner == "" {
			return fmt.Errorf("owner is required")
		}

		// Initialize first stage run if not present
		if len(job.StageRuns) == 0 {
			job.StageRuns = append(job.StageRuns, model.StageRun{
				Stage:     model.StatusBacklog,
				StartedAt: time.Now(),
			})
		}

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
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return fmt.Errorf("failed to get verbose flag: %w", err)
		}

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
				if verbose {
					fmt.Printf("      Freshness: %s, Reliability: %s\n", m.Freshness, m.Reliability)
				}
			}
		}

		if verbose {
			if len(job.SuccessCriteria) > 0 {
				fmt.Println("\nSuccess Criteria:")
				for i, sc := range job.SuccessCriteria {
					fmt.Printf("  %d. %s\n", i+1, sc)
				}
			}

			if len(job.Constraints) > 0 {
				fmt.Println("\nConstraints:")
				for i, c := range job.Constraints {
					fmt.Printf("  %d. %s\n", i+1, c)
				}
			}

			if len(job.ScopeIn) > 0 {
				fmt.Println("\nScope In:")
				for i, s := range job.ScopeIn {
					fmt.Printf("  %d. %s\n", i+1, s)
				}
			}

			if len(job.ScopeOut) > 0 {
				fmt.Println("\nScope Out:")
				for i, s := range job.ScopeOut {
					fmt.Printf("  %d. %s\n", i+1, s)
				}
			}

			if len(job.StageRuns) > 0 {
				fmt.Println("\nStage Runs:")
				for i, sr := range job.StageRuns {
					fmt.Printf("  %d. %s - Started: %s", i+1, sr.Stage, sr.StartedAt.Format("2006-01-02 15:04:05"))
					if sr.CompletedAt != nil {
						fmt.Printf(", Completed: %s", sr.CompletedAt.Format("2006-01-02 15:04:05"))
					}
					fmt.Println()
					if sr.Notes != "" {
						fmt.Printf("      Notes: %s\n", sr.Notes)
					}
				}
			}

			if len(job.Artifacts) > 0 {
				fmt.Println("\nArtifacts:")
				for i, art := range job.Artifacts {
					fmt.Printf("  %d. Type: %s, File: %s, Version: %s\n", i+1, art.Type, art.Content, art.Version)
					fmt.Printf("      Created: %s\n", art.CreatedAt.Format("2006-01-02 15:04:05"))
				}
			}
		}

		if len(job.QualityChecks) > 0 {
			fmt.Println("\nQuality Checks:")
			for _, qc := range job.QualityChecks {
				status := statusPass
				if !qc.Passed {
					status = statusFail
				}
				fmt.Printf("  %s: %s - Score: %d/%d - %s\n",
					qc.GateType, status, qc.Score, qc.MaxScore, qc.CheckedAt.Format("2006-01-02 15:04"))
				if len(qc.DefectCodes) > 0 {
					fmt.Printf("    Defects: %v\n", qc.DefectCodes)
				}
				if verbose && qc.Notes != "" {
					fmt.Printf("    Notes: %s\n", qc.Notes)
				}
			}
		}

		return nil
	},
}

var jobTimelineCmd = &cobra.Command{
	Use:   "timeline <job-id>",
	Short: "Show job timeline",
	Long:  `Show chronological timeline of job stages and quality checks.`,
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

		fmt.Printf("Timeline for job %s: %s\n\n", job.ID, job.Title)

		// Combine stage runs and quality checks into timeline
		type timelineEvent struct {
			Time    time.Time
			Type    string
			Stage   string
			Details string
		}

		var events []timelineEvent

		// Add stage runs
		for _, sr := range job.StageRuns {
			events = append(events, timelineEvent{
				Time:    sr.StartedAt,
				Type:    "STAGE_START",
				Stage:   string(sr.Stage),
				Details: fmt.Sprintf("Started %s", sr.Stage),
			})
			if sr.CompletedAt != nil {
				events = append(events, timelineEvent{
					Time:    *sr.CompletedAt,
					Type:    "STAGE_END",
					Stage:   string(sr.Stage),
					Details: fmt.Sprintf("Completed %s", sr.Stage),
				})
			}
		}

		// Add quality checks
		for _, qc := range job.QualityChecks {
			status := statusPass
			if !qc.Passed {
				status = statusFail
			}
			events = append(events, timelineEvent{
				Time:    qc.CheckedAt,
				Type:    "QUALITY_CHECK",
				Stage:   qc.GateType,
				Details: fmt.Sprintf("%s: %s (%d/%d)", qc.GateType, status, qc.Score, qc.MaxScore),
			})
		}

		// Sort by time
		for i := 0; i < len(events)-1; i++ {
			for j := i + 1; j < len(events); j++ {
				if events[i].Time.After(events[j].Time) {
					events[i], events[j] = events[j], events[i]
				}
			}
		}

		// Display timeline
		for _, event := range events {
			fmt.Printf("%s [%s] %s\n", event.Time.Format("2006-01-02 15:04:05"), event.Type, event.Details)
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
	jobCmd.AddCommand(jobTimelineCmd)
	jobCmd.AddCommand(jobAdvanceCmd)
	jobCmd.AddCommand(jobRejectCmd)

	jobCreateCmd.Flags().String("from-file", "", "Create job from JSON file")
	jobCreateCmd.Flags().String("title", "", "Job title")
	jobCreateCmd.Flags().String("owner", "", "Job owner")
	jobCreateCmd.Flags().String("type", "1pager", "Output type (1pager, comparison, prd, retrospective)")

	jobShowCmd.Flags().Bool("verbose", false, "Show detailed information")

	jobRejectCmd.Flags().String("reason", "", "Rejection reason")
	jobRejectCmd.Flags().StringSlice("defects", []string{}, "Defect codes (e.g., D01,D02)")
}
