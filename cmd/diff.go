package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var diffProject string // the project file to diff

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show differences between current project file and last commit",
	Run: func(cmd *cobra.Command, args []string) {
		if diffProject == "" {
			fmt.Println("Error: specify the project file with -p")
			return
		}

		// Ensure the project file exists
		if _, err := os.Stat(diffProject); os.IsNotExist(err) {
			fmt.Println("Error: project file does not exist:", diffProject)
			return
		}

		repoPath := filepath.Join(".", ".mygit")
		commitsPath := filepath.Join(repoPath, "commits")

		// Ensure repo exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository not initialized. Run 'mygit init' first.")
			return
		}

		// Read current file
		currentContentBytes, err := os.ReadFile(diffProject)
		if err != nil {
			fmt.Println("Error reading project file:", err)
			return
		}
		currentContent := strings.Split(strings.ReplaceAll(string(currentContentBytes), "\x00", ""), "\n")

		// Find the latest commit for this project
		latestID := getLatestProjectCommitID(commitsPath, diffProject)
		if latestID == 0 {
			fmt.Println("No commits found for project:", diffProject)
			return
		}

		commitFile := filepath.Join(commitsPath, fmt.Sprintf("%s_%d.json", diffProject, latestID))
		data, err := os.ReadFile(commitFile)
		if err != nil {
			fmt.Println("Error reading commit file:", err)
			return
		}

		var commit map[string]interface{}
		_ = json.Unmarshal(data, &commit)

		files := commit["files"].(map[string]interface{})
		lastContent := strings.Split(files[diffProject].(string), "\n")

		fmt.Printf("Diff for project '%s' compared to last commit:\n", diffProject)
		printDiff(lastContent, currentContent)
	},
}

func init() {
	diffCmd.Flags().StringVarP(&diffProject, "project", "p", "", "Project (.txt) file to diff")
}



// helper: print simple line-by-line diff
func printDiff(oldLines, newLines []string) {
	maxLen := len(oldLines)
	if len(newLines) > maxLen {
		maxLen = len(newLines)
	}

	for i := 0; i < maxLen; i++ {
		oldLine, newLine := "", ""
		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}
		if oldLine != newLine {
			fmt.Printf("-%s\n+%s\n", oldLine, newLine)
		}
	}
}
