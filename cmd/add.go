package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files/folders to repository staging",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if repoName == "" {
			fmt.Println("Error: specify repository with -r")
			return
		}

		repoPath := filepath.Join(".mygit", "repositories", repoName)
		stagedPath := filepath.Join(repoPath, "staged")
		os.MkdirAll(stagedPath, 0755)

		for _, path := range args {
			info, err := os.Stat(path)
			if err != nil {
				fmt.Println("Cannot access:", path)
				continue
			}
			if info.IsDir() {
				copyDir(path, filepath.Join(stagedPath, filepath.Base(path)))
			} else {
				copyFile(path, filepath.Join(stagedPath, filepath.Base(path)))
			}
			fmt.Println("Added", path, "to repository", repoName)
		}
	},
}

// Copy functions
func copyFile(src, dst string) {
	input, _ := os.Open(src)
	defer input.Close()
	os.MkdirAll(filepath.Dir(dst), 0755)
	output, _ := os.Create(dst)
	defer output.Close()
	io.Copy(output, input)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, info.Mode())
	})
}

func init() {
	addCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repository name")
}
