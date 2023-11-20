package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/atotto/clipboard"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo struct {
	Path     string
	Name     string
	Template string
}

type CommitInfo struct {
	Hash     string
	Author   string
	Commiter string
	Message  string
	RepoName string
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
					output = append(output, Repo{path, info.Name(), ""})
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
			PullAll(make(chan string))
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

func PullAll(statusChan chan string) {
	wg := sync.WaitGroup{}

	pull := func(repo Repo, errorsChan chan string) {
		r, err := git.PlainOpen(repo.Path)
		defer wg.Done()
		if err != nil {
			statusChan <- "Error: " + err.Error()
			return
		}
		w, err := r.Worktree()
		if err != nil {
			statusChan <- "Error: " + err.Error()
			return
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err.Error() != "already up-to-date" {
			statusChan <- "Error: " + repo.Name + ": " + err.Error()
			return
		}
	}

	statusChan <- "Beginning pulling " + fmt.Sprint(len(Config.Repos)) + " repos"
	for _, repo := range Config.Repos {
		wg.Add(1)
		go pull(repo, statusChan)
		time.Sleep(time.Second)
	}
	statusChan <- "All started"
	time.Sleep(time.Second * 3)
	wg.Wait()
	statusChan <- "Pull All finished"
}

func FindCommits() string {
	fCmd := flag.NewFlagSet("f", flag.ExitOnError)
	fNoClip := fCmd.Bool("no-clipboard", false, "no-clipboard")
	fFile := fCmd.String("f", "", "f")
	fPrint := fCmd.Bool("p", false, "p")
	fCmd.Parse(os.Args[2:])
	taskId := fCmd.Args()[0]
	commitsChan := make(chan string)

	wg := sync.WaitGroup{}

	findCommit := func(repo Repo, commitsC chan string) {
		defer wg.Done()
		tmpl := template.New("repoTemplate")
		repoTemplate := repo.Template
		if repoTemplate == "" {
			repoTemplate = Config.DefaultTemplate
		}
		tmpl, err := tmpl.Parse(repoTemplate)
		if err != nil {
			fmt.Println(err)
			return
		}
		r, err := git.PlainOpen(repo.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		ci, err := r.CommitObjects()
		if err != nil {
			fmt.Println(err)
			return
		}
		ci.ForEach(func(c *object.Commit) error {
			buf := &bytes.Buffer{}
			if strings.Contains(c.Message, taskId) {
				info := CommitInfo{c.Hash.String(), c.Author.Name, c.Committer.Name, c.Message, repo.Name}
				tmpl.Execute(buf, info)
				commitsC <- buf.String()
			}
			return nil
		})
	}
	for _, repo := range Config.Repos {
		wg.Add(1)
		fmt.Println(repo.Name)
		go findCommit(repo, commitsChan)
	}
	found := ""
	foundCnt := 0
	go func() {
		commitSet := make(map[string]bool, 0)
		for str := range commitsChan {
			commitSet[str] = true
		}
		for str := range commitSet {
			found += str + "\n"
		}
		foundCnt = len(commitSet)
	}()
	wg.Wait()
	close(commitsChan)
	if !*fNoClip {
		clipboard.WriteAll(found)
	}
	if *fFile != "" {
		file, err := os.Create(*fFile)
		if err != nil {
			fmt.Println("Couldn't create output file")
		}
		file.WriteString(found)
	}
	if *fPrint {
		fmt.Println(found)
	}
	fmt.Println("Done (" + fmt.Sprint(foundCnt) + " results)")
	return found
}
