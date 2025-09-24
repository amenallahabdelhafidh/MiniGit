package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var statusRepo string
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of files in a repository",
	Run: func(cmd *cobra.Command, args []string) {
		if statusRepo == "" {
			fmt.Println("Error: specify the repository with -r")
			return
		}

		repoPath := filepath.Join(".", ".mygit", "repositories", statusRepo)
		commitsPath := filepath.Join(".", ".mygit", "commits")

		// Ensure repository exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Println("Repository does not exist:", statusRepo)
			return
		}

		// Gather all commits for this repository
		commitFiles, err := os.ReadDir(commitsPath)
		if err != nil {
			fmt.Println("Error reading commits folder:", err)
			return
		}

		// Map: file path -> last commit info
		lastCommit := make(map[string]string)

		// Sort commit files to process in chronological order
		sort.Slice(commitFiles, func(i, j int) bool {
			return commitFiles[i].Name() < commitFiles[j].Name()
		})

		for _, f := range commitFiles {
			// Only commits for this repository
			if !filepath.HasPrefix(f.Name(), statusRepo+"_") {
				continue
			}

			commitPath := filepath.Join(commitsPath, f.Name())
			data, err := os.ReadFile(commitPath)
			if err != nil {
				continue
			}

			var commit map[string]interface{}
			if err := json.Unmarshal(data, &commit); err != nil {
				continue
			}

			msg, _ := commit["message"].(string)
			files, _ := commit["files"].(map[string]interface{})

			for file := range files {
				lastCommit[file] = msg
			}
		}

		// Walk repository files
		fmt.Printf("Project Status for '%s':\n", statusRepo)
		filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			relPath, _ := filepath.Rel(repoPath, path)
			msg, ok := lastCommit[relPath]
			if !ok {
				msg = "Untracked"
			}
			fmt.Printf("  %s - Last commit: %s\n", relPath, msg)
			return nil
		})
	},
}

func init() {
	statusCmd.Flags().StringVarP(&statusRepo, "repo", "r", "", "Repository name")
}
