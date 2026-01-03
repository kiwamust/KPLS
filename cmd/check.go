package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"kpls/internal/checker"
	"kpls/internal/model"
	"kpls/internal/store"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Perform quality gate checks",
	Long:  `Run quality gate checks (IQC, IPQC-1, IPQC-2, FQC) on jobs.`,
}

var checkIQCCmd = &cobra.Command{
	Use:   "iqc [job-id]",
	Short: "Run IQC (Incoming Quality Control) check",
	Long:  `Perform incoming quality control check on a job.`,
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

		if job.Status != model.StatusIQC {
			return fmt.Errorf("job is not in IQC stage (current: %s)", job.Status)
		}

		iqcChecker := checker.NewIQCChecker(s)
		check, err := iqcChecker.Check(job)
		if err != nil {
			return err
		}

		job.QualityChecks = append(job.QualityChecks, *check)

		fmt.Printf("IQC Check Results:\n")
		fmt.Printf("  Score: %d/%d\n", check.Score, check.MaxScore)
		fmt.Printf("  Passed: %v\n", check.Passed)
		if len(check.DefectCodes) > 0 {
			fmt.Printf("  Defect Codes: %v\n", check.DefectCodes)
		}

		if check.Passed {
			fmt.Println("\nJob passed IQC. You can advance to next stage.")
		} else {
			fmt.Println("\nJob failed IQC. Please fix issues before advancing.")
		}

		return s.SaveJob(job)
	},
}

var checkIPQC1Cmd = &cobra.Command{
	Use:   "ipqc1 [job-id]",
	Short: "Run IPQC-1 (Skeleton Gate) check",
	Long:  `Perform IPQC-1 skeleton gate check on a job's artifact.`,
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

		if job.Status != model.StatusIPQC1 {
			return fmt.Errorf("job is not in IPQC-1 stage (current: %s)", job.Status)
		}

		// Get artifact content (simplified - in real implementation, load from file)
		artifactContent := ""
		if len(job.Artifacts) > 0 {
			artifactContent = job.Artifacts[len(job.Artifacts)-1].Content
		}

		ipqcChecker := checker.NewIPQCChecker(s, "IPQC-1")
		check, err := ipqcChecker.Check(job, artifactContent)
		if err != nil {
			return err
		}

		job.QualityChecks = append(job.QualityChecks, *check)

		fmt.Printf("IPQC-1 Check Results:\n")
		fmt.Printf("  Score: %d/%d\n", check.Score, check.MaxScore)
		fmt.Printf("  Passed: %v\n", check.Passed)
		if len(check.DefectCodes) > 0 {
			fmt.Printf("  Defect Codes: %v\n", check.DefectCodes)
		}

		if check.Passed {
			fmt.Println("\nJob passed IPQC-1. You can advance to next stage.")
		} else {
			fmt.Println("\nJob failed IPQC-1. Please fix issues before advancing.")
		}

		return s.SaveJob(job)
	},
}

var checkIPQC2Cmd = &cobra.Command{
	Use:   "ipqc2 [job-id]",
	Short: "Run IPQC-2 (Evidence Gate) check",
	Long:  `Perform IPQC-2 evidence gate check on a job's artifact.`,
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

		if job.Status != model.StatusIPQC2 {
			return fmt.Errorf("job is not in IPQC-2 stage (current: %s)", job.Status)
		}

		// Get artifact content
		artifactContent := ""
		if len(job.Artifacts) > 0 {
			artifactContent = job.Artifacts[len(job.Artifacts)-1].Content
		}

		ipqcChecker := checker.NewIPQCChecker(s, "IPQC-2")
		check, err := ipqcChecker.Check(job, artifactContent)
		if err != nil {
			return err
		}

		job.QualityChecks = append(job.QualityChecks, *check)

		fmt.Printf("IPQC-2 Check Results:\n")
		fmt.Printf("  Score: %d/%d\n", check.Score, check.MaxScore)
		fmt.Printf("  Passed: %v\n", check.Passed)
		if len(check.DefectCodes) > 0 {
			fmt.Printf("  Defect Codes: %v\n", check.DefectCodes)
		}

		if check.Passed {
			fmt.Println("\nJob passed IPQC-2. You can advance to next stage.")
		} else {
			fmt.Println("\nJob failed IPQC-2. Please fix issues before advancing.")
		}

		return s.SaveJob(job)
	},
}

var checkFQCCmd = &cobra.Command{
	Use:   "fqc [job-id]",
	Short: "Run FQC (Final Quality Control) check",
	Long:  `Perform final quality control check on a job.`,
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

		if job.Status != model.StatusFQC {
			return fmt.Errorf("job is not in FQC stage (current: %s)", job.Status)
		}

		// Get artifact content
		artifactContent := ""
		if len(job.Artifacts) > 0 {
			artifactContent = job.Artifacts[len(job.Artifacts)-1].Content
		}

		fqcChecker := checker.NewFQCChecker(s)
		check, err := fqcChecker.Check(job, artifactContent)
		if err != nil {
			return err
		}

		job.QualityChecks = append(job.QualityChecks, *check)

		fmt.Printf("FQC Check Results:\n")
		fmt.Printf("  Score: %d/%d\n", check.Score, check.MaxScore)
		fmt.Printf("  Passed: %v\n", check.Passed)
		if len(check.DefectCodes) > 0 {
			fmt.Printf("  Defect Codes: %v\n", check.DefectCodes)
		}

		if check.Passed {
			fmt.Println("\nJob passed FQC. You can mark as Done.")
		} else {
			fmt.Println("\nJob failed FQC. Please fix issues before completing.")
		}

		return s.SaveJob(job)
	},
}

func init() {
	checkCmd.AddCommand(checkIQCCmd)
	checkCmd.AddCommand(checkIPQC1Cmd)
	checkCmd.AddCommand(checkIPQC2Cmd)
	checkCmd.AddCommand(checkFQCCmd)
}
