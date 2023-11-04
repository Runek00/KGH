package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Repo struct {
	Path     string
	Name     string
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

func RepoTool() {
	printRepoList()
	repoNumPick()
}

func repoNumPick() {
	for {
		fmt.Print(`
	Which repo do you want to edit? 
	(number) repo number [number] from the list
	(s) save
	(x) exit 
	(l) show the list
	(p) pull all
	`)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		input = strings.ToLower(strings.TrimSpace(input))
		switch input {
		case "s":
			fallthrough
		case "save":
			SaveConfig()
		case "x":
			fallthrough
		case "exit":
			return
		case "l":
			fallthrough
		case "list":
			printRepoList()
			continue
		case "p":
			PullAll()
		default:
			idx, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("wrong command")
				continue
			}
			editMenu(idx)
		}
	}
}

func editMenu(idx int) {
	paths := make([]string, 0)
	for path := range Config.Repos {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	repo := Config.Repos[paths[idx]]
MenuLoop:
	for {
		fmt.Println(fmt.Sprint(idx) + ": path: " + repo.Path + ", name: " + repo.Name + ", template: " + repo.Template)
		fmt.Print(`
	What do you want to change?
	(p) change path
	(n) change repo's name
	(t) set/change custom template
	(d) set template to default
	(x) go back to repo list
	`)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch strings.TrimSpace(input) {
		case "p":
			repo = showFieldChange(repo, "new path", PathChange)
		case "n":
			repo = showFieldChange(repo, "new name", NameChange)
		case "t":
			repo = showFieldChange(repo, "new template", TemplateChange)
		case "d":
			repo = TemplateChange(repo, "")
		case "x":
			fmt.Print("\n")
			break MenuLoop
		default:
			fmt.Println("No such option")
		}
	}
}

func showFieldChange(repo Repo, text string, updater func(Repo, string) Repo) Repo {
	fmt.Println(text + " (leave empty to cancel):")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return repo
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return repo
	}
	return updater(repo, input)
}

func PathChange(repo Repo, path string) Repo {
	delete(Config.Repos, repo.Path)
	repo.Path = path
	Config.Repos[repo.Path] = repo
	return repo
}

func NameChange(repo Repo, name string) Repo {
	repo.Name = name
	Config.Repos[repo.Path] = repo
	return repo
}

func TemplateChange(repo Repo, template string) Repo {
	repo.Template = template
	Config.Repos[repo.Path] = repo
	return repo
}

func printRepoList() {
	fmt.Println("Repos with empty template use the default: " + Config.DefaultTemplate)
	paths := make([]string, 0)
	for path := range Config.Repos {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for idx, path := range paths {
		repo := Config.Repos[path]
		fmt.Println(fmt.Sprint(idx) + ": path: " + repo.Path + ", name: " + repo.Name + ", template: " + repo.Template)
	}
}

func PullAll() {
	panic("unimplemented")
}
