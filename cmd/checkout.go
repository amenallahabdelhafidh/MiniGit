package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	

	"github.com/spf13/cobra"
)

var checkoutProject string
var checkoutID int

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Restore a project file to a previous commit",
	Run: func(cmd *cobra.Command, args []string) {
		if checkoutProject == "" {
			fmt.Println("Error: specify the project file with -p")
			return
		}

		repoPath := filepath.Join(".", ".mygit")
		commitsPath := filepath.Join(repoPath, "commits")

		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository not initialized. Run 'mygit init' first.")
			return
		}

		// Determine which commit to restore
		var commitFile string
		if checkoutID == 0 {
			// Get latest commit for this project
			checkoutID = getLatestProjectCommitID(commitsPath, checkoutProject)
			if checkoutID == 0 {
				fmt.Println("No commits found for project:", checkoutProject)
				return
			}
		}

		commitFile = filepath.Join(commitsPath, fmt.Sprintf("%s_%d.json", checkoutProject, checkoutID))

		data, err := os.ReadFile(commitFile)
		if err != nil {
			fmt.Println("Error reading commit file:", err)
			return
		}

		commit := make(map[string]interface{})
		if err := json.Unmarshal(data, &commit); err != nil {
			fmt.Println("Error parsing commit JSON:", err)
			return
		}

		filesMap := commit["files"].(map[string]interface{})
		content := filesMap[checkoutProject].(string)

		// Restore project file
		if err := os.WriteFile(checkoutProject, []byte(content), 0644); err != nil {
			fmt.Println("Error restoring project file:", err)
			return
		}

		fmt.Printf("Project '%s' restored to commit %d (%s)\n", checkoutProject, checkoutID, commit["timestamp"])
	},
}

func init() {
	checkoutCmd.Flags().StringVarP(&checkoutProject, "project", "p", "", "Project (.txt) file to restore")
	checkoutCmd.Flags().IntVarP(&checkoutID, "commit", "c", 0, "Commit ID to restore (default: latest)")
}

// Get latest commit ID for a specific project
func getLatestProjectCommitID(commitsPath, project string) int {
	files, err := os.ReadDir(commitsPath)
	if err != nil {
		return 0
	}

	maxID := 0
	for _, f := range files {
		var id int
		if n, _ := fmt.Sscanf(f.Name(), project+"_%d.json", &id); n == 1 {
			if id > maxID {
				maxID = id
			}
		}
	}
	return maxID
}
