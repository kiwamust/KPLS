package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates",
	Long:  `List and show templates.`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available output templates.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		templates := []string{
			"1pager_decision_memo",
			"comparison_evaluation",
			"prd_system",
		}

		fmt.Println("Available templates:")
		for _, t := range templates {
			fmt.Printf("  - %s\n", t)
		}

		return nil
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show [template-name]",
	Short: "Show template content",
	Long:  `Display the content of a template.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		templateName := args[0]
		path := fmt.Sprintf("kpls/templates/%s.md", templateName)

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("template not found: %s", templateName)
		}

		fmt.Print(string(content))
		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
}
