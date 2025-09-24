package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var commitMessage string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new commit for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repoName, _ := cmd.Flags().GetString("repo")
		if repoName == "" {
			fmt.Println("Error: repository name is required. Use -r <repo>")
			return
		}
		if commitMessage == "" {
			fmt.Println("Error: commit message is required. Use -m \"message\"")
			return
		}

		repoPath := filepath.Join(".mygit", "repositories", repoName)
		commitsPath := filepath.Join(".mygit", "commits")

		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository does not exist:", repoName)
			return
		}

		// Load latest commit for this repo
		latestFiles := make(map[string]string)
		files, _ := os.ReadDir(commitsPath)
		maxID := 0
		for _, f := range files {
			name := f.Name()
			if !strings.HasPrefix(name, repoName+"_") || !strings.HasSuffix(name, ".json") {
				continue
			}
			idStr := strings.TrimSuffix(strings.TrimPrefix(name, repoName+"_"), ".json")
			id := 0
			fmt.Sscanf(idStr, "%d", &id)
			if id > maxID {
				maxID = id
			}
		}
		if maxID > 0 {
			data, _ := os.ReadFile(filepath.Join(commitsPath, fmt.Sprintf("%s_%d.json", repoName, maxID)))
			var commitData map[string]interface{}
			json.Unmarshal(data, &commitData)
			if filesMap, ok := commitData["files"].(map[string]interface{}); ok {
				for k, v := range filesMap {
					if s, ok := v.(string); ok {
						latestFiles[k] = s
					}
				}
			}
		}

		// Track changed files
		changedFiles := make(map[string]string)

		filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			relPath, _ := filepath.Rel(repoPath, path)
			contentBytes, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			contentStr := strings.ReplaceAll(string(contentBytes), "\x00", "")
			if prev, ok := latestFiles[relPath]; !ok || prev != contentStr {
				changedFiles[relPath] = contentStr
			}
			return nil
		})

		if len(changedFiles) == 0 {
			fmt.Println("No changes to commit.")
			return
		}

		// Prepare new commit
		newID := maxID + 1
		commit := map[string]interface{}{
			"id":        newID,
			"repo":      repoName,
			"message":   commitMessage,
			"timestamp": time.Now().Format(time.RFC3339),
			"files":     changedFiles,
		}

		commitFile := filepath.Join(commitsPath, fmt.Sprintf("%s_%d.json", repoName, newID))
		data, _ := json.MarshalIndent(commit, "", "  ")
		os.WriteFile(commitFile, data, 0644)

		fmt.Printf("Committed %d changed files as commit %d\n", len(changedFiles), newID)
	},
}

func init() {
	commitCmd.Flags().StringP("repo", "r", "", "Repository name")
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
}
