package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates",
	Long:  `List, show, add, and validate templates.`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available output templates.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		entries, err := os.ReadDir(templatesDir)
		if err != nil {
			return fmt.Errorf("failed to read templates directory: %w", err)
		}

		var templates []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(entry.Name(), ".md") {
				name := strings.TrimSuffix(entry.Name(), ".md")
				templates = append(templates, name)
			}
		}

		if len(templates) == 0 {
			fmt.Println("No templates found.")
			return nil
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
		path := fmt.Sprintf("%s/%s.md", templatesDir, templateName)

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("template not found: %s", templateName)
		}

		fmt.Print(string(content))
		return nil
	},
}

var templateAddCmd = &cobra.Command{
	Use:   "add <file.md>",
	Short: "Add a template",
	Long:  `Add a template file to the templates directory.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		sourceFile := args[0]
		templatesDir := "kpls/templates"

		// Read source file
		content, err := os.ReadFile(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to read source file: %w", err)
		}

		// Validate YAML frontmatter
		if err := validateTemplateFrontmatter(content); err != nil {
			return fmt.Errorf("template validation failed: %w", err)
		}

		// Extract template name from filename
		baseName := filepath.Base(sourceFile)
		if !strings.HasSuffix(baseName, ".md") {
			return fmt.Errorf("template file must have .md extension")
		}

		// Ensure templates directory exists
		if err := os.MkdirAll(templatesDir, 0755); err != nil {
			return fmt.Errorf("failed to create templates directory: %w", err)
		}

		// Copy file to templates directory
		destPath := filepath.Join(templatesDir, baseName)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write template file: %w", err)
		}

		fmt.Printf("Template added: %s\n", baseName)
		return nil
	},
}

var templateValidateCmd = &cobra.Command{
	Use:   "validate [template-name]",
	Short: "Validate a template",
	Long:  `Validate template structure and required fields.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		var templateName string
		if len(args) > 0 {
			templateName = args[0]
		} else {
			// Validate all templates
			templatesDir := "kpls/templates"
			entries, err := os.ReadDir(templatesDir)
			if err != nil {
				return fmt.Errorf("failed to read templates directory: %w", err)
			}

			allValid := true
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
					continue
				}
				name := strings.TrimSuffix(entry.Name(), ".md")
				path := filepath.Join(templatesDir, entry.Name())
				content, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("❌ %s: failed to read file\n", name)
					allValid = false
					continue
				}
				if err := validateTemplateFrontmatter(content); err != nil {
					fmt.Printf("❌ %s: %v\n", name, err)
					allValid = false
				} else {
					fmt.Printf("✅ %s: valid\n", name)
				}
			}

			if !allValid {
				return fmt.Errorf("some templates failed validation")
			}
			return nil
		}

		// Validate single template
		path := fmt.Sprintf("%s/%s.md", templatesDir, templateName)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("template not found: %s", templateName)
		}

		if err := validateTemplateFrontmatter(content); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		fmt.Printf("✅ %s: valid\n", templateName)
		return nil
	},
}

// validateTemplateFrontmatter validates YAML frontmatter in template
func validateTemplateFrontmatter(content []byte) error {
	contentStr := string(content)

	// Check for frontmatter delimiters
	if !strings.HasPrefix(contentStr, "---\n") {
		return fmt.Errorf("missing YAML frontmatter (must start with '---')")
	}

	// Find end of frontmatter
	endIdx := strings.Index(contentStr[4:], "\n---\n")
	if endIdx == -1 {
		return fmt.Errorf("missing YAML frontmatter end delimiter")
	}

	frontmatter := contentStr[4 : endIdx+4]

	// Basic validation: check for 'type' field
	if !strings.Contains(frontmatter, "type:") {
		return fmt.Errorf("missing required field 'type' in frontmatter")
	}

	return nil
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateAddCmd)
	templateCmd.AddCommand(templateValidateCmd)
}
