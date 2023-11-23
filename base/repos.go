package base

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func getParamsAndFindCommits() string {
	fCmd := flag.NewFlagSet("f", flag.ExitOnError)
	fNoClip := fCmd.Bool("no-clipboard", false, "no-clipboard")
	fFile := fCmd.String("f", "", "f")
	fPrint := fCmd.Bool("p", false, "p")
	fCmd.Parse(os.Args[2:])
	searchPhrase := fCmd.Args()[0]

	commitSet := FindCommits(searchPhrase)
	found := ""
	foundCnt := 0
	for str := range commitSet {
		found += str + "\n"
	}
	foundCnt = len(commitSet)
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

func FindCommits(searchPhrase string) map[string]bool {
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
			if strings.Contains(c.Message, searchPhrase) {
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
	commitSet := make(map[string]bool, 0)
	go func() {
		for str := range commitsChan {
			commitSet[str] = true
		}
	}()
	wg.Wait()
	close(commitsChan)
	return commitSet
}
