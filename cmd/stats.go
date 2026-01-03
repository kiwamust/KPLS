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

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics and metrics",
	Long:  `Display workflow statistics including WIP, defect rates, and stage durations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.NewFileStore()
		if err != nil {
			return err
		}

		jobs, err := s.ListJobs()
		if err != nil {
			return err
		}

		// WIP by stage
		fmt.Println("=== WIP by Stage ===")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Stage", "Count", "WIP Limit"})
		table.SetBorder(false)

		stages := workflow.GetAllStages()
		for _, stage := range stages {
			count := 0
			for _, job := range jobs {
				if job.Status == stage {
					count++
				}
			}

			limit := "-"
			if stage == model.StatusSkeleton {
				limit = "5"
			} else if stage == model.StatusDraft {
				limit = "3"
			}

			table.Append([]string{
				string(stage),
				fmt.Sprintf("%d", count),
				limit,
			})
		}
		table.Render()

		// Defect statistics
		fmt.Println("\n=== Defect Statistics ===")
		defectCounts := make(map[string]int)
		for _, job := range jobs {
			for _, qc := range job.QualityChecks {
				for _, code := range qc.DefectCodes {
					defectCounts[code]++
				}
			}
		}

		if len(defectCounts) > 0 {
			table2 := tablewriter.NewWriter(os.Stdout)
			table2.SetHeader([]string{"Defect Code", "Count"})
			table2.SetBorder(false)
			for code, count := range defectCounts {
				table2.Append([]string{code, fmt.Sprintf("%d", count)})
			}
			table2.Render()
		} else {
			fmt.Println("No defects recorded.")
		}

		// Stage duration (average)
		fmt.Println("\n=== Average Stage Duration ===")
		durations := make(map[model.JobStatus][]time.Duration)
		for _, job := range jobs {
			for i := 0; i < len(job.StageRuns)-1; i++ {
				run := job.StageRuns[i]
				if run.CompletedAt != nil {
					duration := run.CompletedAt.Sub(run.StartedAt)
					durations[run.Stage] = append(durations[run.Stage], duration)
				}
			}
		}

		table3 := tablewriter.NewWriter(os.Stdout)
		table3.SetHeader([]string{"Stage", "Avg Duration", "Samples"})
		table3.SetBorder(false)

		for _, stage := range stages {
			if durations[stage] != nil && len(durations[stage]) > 0 {
				var total time.Duration
				for _, d := range durations[stage] {
					total += d
				}
				avg := total / time.Duration(len(durations[stage]))
				table3.Append([]string{
					string(stage),
					avg.Round(time.Hour).String(),
					fmt.Sprintf("%d", len(durations[stage])),
				})
			}
		}
		table3.Render()

		return nil
	},
}
