package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new mini git repository",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := filepath.Join(".", ".mygit")
		if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
			fmt.Println("Repository already initialized.")
			return
		}

		// Create .mygit/ structure
		os.MkdirAll(filepath.Join(repoPath, "commits"), 0755)
		os.MkdirAll(filepath.Join(repoPath, "repositories"), 0755) // new for multiple repos

		// Create empty index.json (global metadata)
		indexFile := filepath.Join(repoPath, "index.json")
		index := make(map[string]interface{})
		data, _ := json.MarshalIndent(index, "", "  ")
		os.WriteFile(indexFile, data, 0644)

		// Optional: create global.json for repo list
		globalFile := filepath.Join(repoPath, "global.json")
		repos := []string{} // empty list of repositories
		data, _ = json.MarshalIndent(repos, "", "  ")
		os.WriteFile(globalFile, data, 0644)

		fmt.Println("Initialized empty mini git repository in", repoPath)
		fmt.Println("You can now add repositories using `addrepo` command.")
	},
}
