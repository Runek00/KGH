package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Repo struct {
	Name     string
	Path     string
	Template string
}

func FindRepos(path string) []Repo {
	output := make([]Repo, 0)
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			direntries, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, entry := range direntries {
				if entry.IsDir() && entry.Name() == ".git" {
					output = append(output, Repo{info.Name(), path, ""})
					break
				}
			}
		}
		return nil
	})
	return output
}
