package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var repoName string

// addrepo command
var addRepoCmd = &cobra.Command{
	Use:   "addrepo",
	Short: "Add a new repository",
	Run: func(cmd *cobra.Command, args []string) {
		if repoName == "" {
			fmt.Println("Error: repository name is required. Use -n \"name\"")
			return
		}

		// Ensure .mygit exists
		mygitPath := filepath.Join(".", ".mygit")
		if _, err := os.Stat(mygitPath); os.IsNotExist(err) {
			fmt.Println("Error: Run 'mygit init' first to initialize the main repository.")
			return
		}

		// Create repositories folder if not exist
		reposPath := filepath.Join(mygitPath, "repositories")
		os.MkdirAll(reposPath, 0755)

		// Create the new repo folder
		projectPath := filepath.Join(reposPath, repoName)
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			err := os.Mkdir(projectPath, 0755)
			if err != nil {
				fmt.Println("Error creating repository folder:", err)
				return
			}
			fmt.Println("Repository folder created:", projectPath)
		} else {
			fmt.Println("Repository folder already exists:", projectPath)
		}

		// Update global.json with the new repo
		globalFile := filepath.Join(mygitPath, "global.json")

		var repos []string

		// Read existing global.json if it exists
		if _, err := os.Stat(globalFile); err == nil {
			bytes, err := os.ReadFile(globalFile)
			if err == nil {
				json.Unmarshal(bytes, &repos)
			}
		}

		// Check if repo already exists in the list
		for _, r := range repos {
			if r == repoName {
				fmt.Println("Repository already added to global.json")
				return
			}
		}

		// Append and save
		repos = append(repos, repoName)
		data, _ := json.MarshalIndent(repos, "", "  ")
		err := os.WriteFile(globalFile, data, 0644)
		if err != nil {
			fmt.Println("Error updating global.json:", err)
			return
		}

		fmt.Println("Repository added to global.json:", repoName)
	},
}


func init() {
	addRepoCmd.Flags().StringVarP(&repoName, "name", "n", "", "Name of the repository to add")
}