package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kpls",
	Short: "Knowledge Production Line System",
	Long:  `KPLS is a CLI tool for managing knowledge production workflows with quality gates.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(jobCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(artifactCmd)
	rootCmd.AddCommand(gateCmd)
}
