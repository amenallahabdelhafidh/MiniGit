package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := filepath.Join(".", ".mygit")
		commitsPath := filepath.Join(repoPath, "commits")

		// Ensure repo exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository not initialized. Run 'mygit init' first.")
			return
		}

		// Read all commit files
		files, err := os.ReadDir(commitsPath)
		if err != nil || len(files) == 0 {
			fmt.Println("No commits found.")
			return
		}

		// Sort files by name (assuming numeric IDs)
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() > files[j].Name() // latest first
		})

		for _, f := range files {
			commitFile := filepath.Join(commitsPath, f.Name())
			data, err := os.ReadFile(commitFile)
			if err != nil {
				continue
			}

			var commit map[string]interface{}
			if err := json.Unmarshal(data, &commit); err != nil {
				continue
			}

			id := commit["id"]
			message := commit["message"]
			timestamp := commit["timestamp"]

			fmt.Printf("Commit %v: %v\n", id, message)
			fmt.Printf("Timestamp: %v\n", timestamp)

			if filesMap, ok := commit["files"].(map[string]interface{}); ok {
				fmt.Println("Files:")
				for fname := range filesMap {
					fmt.Printf("  - %s\n", fname)
				}
			}
			fmt.Println("--------------------------------------------------")
		}
	},
}

