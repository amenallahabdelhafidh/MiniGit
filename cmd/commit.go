package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var commitMessage string
var projectFile string // the .txt file you want to commit

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new commit for a single project (.txt file)",
	Run: func(cmd *cobra.Command, args []string) {
		if commitMessage == "" {
			fmt.Println("Error: commit message is required. Use -m \"message\"")
			return
		}

		if projectFile == "" {
			fmt.Println("Error: specify the project file with -p")
			return
		}

		// Ensure the project file exists
		if _, err := os.Stat(projectFile); os.IsNotExist(err) {
			fmt.Println("Error: project file does not exist:", projectFile)
			return
		}

		repoPath := filepath.Join(".", ".mygit")
		commitsPath := filepath.Join(repoPath, "commits")

		// Ensure repo exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository not initialized. Run 'mygit init' first.")
			return
		}

		// Read the project file
		contentBytes, err := os.ReadFile(projectFile)
		if err != nil {
			fmt.Println("Error reading project file:", err)
			return
		}

		contentStr := strings.ReplaceAll(string(contentBytes), "\x00", "")

		// Determine next commit ID for this project only
		projectID := getNextProjectCommitID(commitsPath, projectFile)

		// Prepare commit object
		commit := make(map[string]interface{})
		commit["id"] = projectID
		commit["project"] = projectFile
		commit["message"] = commitMessage
		commit["timestamp"] = time.Now().Format(time.RFC3339)
		commit["files"] = map[string]string{
			projectFile: contentStr,
		}

		// Save commit as JSON file
		commitFile := filepath.Join(commitsPath, fmt.Sprintf("%s_%d.json", projectFile, projectID))
		data, _ := json.MarshalIndent(commit, "", "  ")
		os.WriteFile(commitFile, data, 0644)

		fmt.Printf("Committed project '%s' as commit %d\n", projectFile, projectID)
	},
}

func init() {
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
	commitCmd.Flags().StringVarP(&projectFile, "project", "p", "", "Project (.txt) file to commit")
}

// helper function to determine next commit ID for a specific project
func getNextProjectCommitID(commitsPath, project string) int {
	files, err := os.ReadDir(commitsPath)
	if err != nil {
		return 1
	}

	maxID := 0
	for _, f := range files {
		// Commit filename format: projectFileName_ID.json
		var id int
		if n, _ := fmt.Sscanf(f.Name(), project+"_%d.json", &id); n == 1 {
			if id > maxID {
				maxID = id
			}
		}
	}
	return maxID + 1
}
