package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"kpls/internal/model"
	"kpls/internal/store"
)

var gateCmd = &cobra.Command{
	Use:   "gate",
	Short: "Manage quality gates",
	Long:  `Record and view quality gate check results.`,
}

var gateCheckCmd = &cobra.Command{
	Use:   "check <job-id>",
	Short: "Record a quality gate check",
	Long:  `Record the result of a quality gate check for a job.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID := args[0]

		gateType, err := cmd.Flags().GetString("gate")
		if err != nil {
			return fmt.Errorf("failed to get gate flag: %w", err)
		}
		if gateType == "" {
			return fmt.Errorf("--gate is required (IQC, IPQC-1, IPQC-2, FQC)")
		}

		score, err := cmd.Flags().GetInt("score")
		if err != nil {
			return fmt.Errorf("failed to get score flag: %w", err)
		}

		maxScore, err := cmd.Flags().GetInt("max-score")
		if err != nil {
			return fmt.Errorf("failed to get max-score flag: %w", err)
		}
		if maxScore == 0 {
			maxScore = 10 // default
		}

		notes, _ := cmd.Flags().GetString("notes")
		defects, _ := cmd.Flags().GetStringSlice("defects")
		checker, _ := cmd.Flags().GetString("checker")
		if checker == "" {
			checker = os.Getenv("USER")
			if checker == "" {
				checker = "unknown"
			}
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

		// Create quality check
		passed := score >= (maxScore * 9 / 10) // 90% threshold
		if gateType == "FQC" {
			passed = score >= (maxScore * 9 / 10) // 90% for FQC
		} else if gateType == "IQC" {
			passed = score >= (maxScore * 8 / 10) // 80% for IQC
		} else {
			passed = score >= (maxScore * 8 / 10) // 80% for IPQC
		}

		qc := model.QualityCheck{
			ID:          store.GenerateDefectID(), // Reuse ID generator
			JobID:       jobID,
			GateType:    gateType,
			CheckedAt:   time.Now(),
			Checker:     checker,
			Score:       score,
			MaxScore:    maxScore,
			Passed:      passed,
			DefectCodes: defects,
			Notes:       notes,
		}

		job.QualityChecks = append(job.QualityChecks, qc)

		// Save job
		if err := s.SaveJob(job); err != nil {
			return fmt.Errorf("failed to save job: %w", err)
		}

		status := "PASS"
		if !passed {
			status = "FAIL"
		}
		fmt.Printf("Quality check recorded: %s %s (%d/%d) - %s\n", gateType, status, score, maxScore, jobID)
		return nil
	},
}

var gateHistoryCmd = &cobra.Command{
	Use:   "history <job-id>",
	Short: "Show quality gate history",
	Long:  `Display all quality gate check results for a job.`,
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

		if len(job.QualityChecks) == 0 {
			fmt.Printf("No quality checks found for job %s\n", jobID)
			return nil
		}

		fmt.Printf("Quality gate history for job %s:\n\n", jobID)
		for i, qc := range job.QualityChecks {
			status := "PASS"
			if !qc.Passed {
				status = "FAIL"
			}
			fmt.Printf("%d. %s - %s (%d/%d)\n", i+1, qc.GateType, status, qc.Score, qc.MaxScore)
			fmt.Printf("   Checked by: %s at %s\n", qc.Checker, qc.CheckedAt.Format("2006-01-02 15:04:05"))
			if len(qc.DefectCodes) > 0 {
				fmt.Printf("   Defects: %v\n", qc.DefectCodes)
			}
			if qc.Notes != "" {
				fmt.Printf("   Notes: %s\n", qc.Notes)
			}
			if qc.RejectReason != "" {
				fmt.Printf("   Reject reason: %s\n", qc.RejectReason)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	gateCmd.AddCommand(gateCheckCmd)
	gateCmd.AddCommand(gateHistoryCmd)

	gateCheckCmd.Flags().String("gate", "", "Gate type (IQC, IPQC-1, IPQC-2, FQC)")
	gateCheckCmd.Flags().Int("score", 0, "Score")
	gateCheckCmd.Flags().Int("max-score", 10, "Maximum score")
	gateCheckCmd.Flags().String("notes", "", "Notes")
	gateCheckCmd.Flags().StringSlice("defects", []string{}, "Defect codes (e.g., D01,D02)")
	gateCheckCmd.Flags().String("checker", "", "Checker name (default: $USER)")
}
