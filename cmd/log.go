package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repoName, _ := cmd.Flags().GetString("repo")
		if repoName == "" {
			fmt.Println("Error: repository name is required. Use -r <repo>")
			return
		}

		repoPath := filepath.Join(".mygit", "repositories", repoName)
		commitsPath := filepath.Join(".mygit", "commits")

		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository does not exist:", repoName)
			return
		}

		// Read all commit files
		files, err := os.ReadDir(commitsPath)
		if err != nil {
			fmt.Println("Error reading commits folder:", err)
			return
		}

		type Commit struct {
			ID        int
			Message   string
			Timestamp string
			FileCount int
		}

		var commits []Commit

		for _, f := range files {
			name := f.Name()
			if !strings.HasPrefix(name, repoName+"_") || !strings.HasSuffix(name, ".json") {
				continue
			}

			// Extract commit ID
			idStr := strings.TrimSuffix(strings.TrimPrefix(name, repoName+"_"), ".json")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}

			// Load commit content
			data, err := os.ReadFile(filepath.Join(commitsPath, name))
			if err != nil {
				continue
			}

			var commitData map[string]interface{}
			json.Unmarshal(data, &commitData)

			filesMap, ok := commitData["files"].(map[string]interface{})
			fileCount := 0
			if ok {
				fileCount = len(filesMap)
			}

			message, _ := commitData["message"].(string)
			timestamp, _ := commitData["timestamp"].(string)

			commits = append(commits, Commit{
				ID:        id,
				Message:   message,
				Timestamp: timestamp,
				FileCount: fileCount,
			})
		}

		// Sort commits descending by ID
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].ID > commits[j].ID
		})

		fmt.Printf("Commit history for repository '%s':\n", repoName)
		for _, c := range commits {
			fmt.Printf("Commit %d | %s\n", c.ID, c.Timestamp)
			fmt.Printf("Message: %s\n", c.Message)
			fmt.Printf("Files: %d\n\n", c.FileCount)
		}

		if len(commits) == 0 {
			fmt.Println("No commits yet.")
		}
	},
}

func init() {
	logCmd.Flags().StringP("repo", "r", "", "Repository name")
}
