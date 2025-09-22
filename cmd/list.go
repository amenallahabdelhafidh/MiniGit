package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all committed projects",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := filepath.Join(".", ".mygit")
		commitsPath := filepath.Join(repoPath, "commits")

		// Ensure repo exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository not initialized. Run 'mygit init' first.")
			return
		}

		files, err := os.ReadDir(commitsPath)
		if err != nil || len(files) == 0 {
			fmt.Println("No commits found.")
			return
		}

		projectSet := make(map[string]struct{})
		for _, f := range files {
			name := f.Name()
			// commit files are in format: projectFileName_ID.json
			parts := strings.SplitN(name, "_", 2)
			if len(parts) == 2 {
				projectSet[parts[0]] = struct{}{}
			}
		}

		fmt.Println("Committed projects:")
		for project := range projectSet {
			fmt.Println("  -", project)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
