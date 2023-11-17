package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	ReadConfig()
	/*
		options:
			c	config						open config in vim / default editor
			a	add-repos <path>			add all repos in path to repo list
			r	remove-repos <path>			remove all repos in path from repo list
			l	repo-list					show repo list, interact with the repos
			t	template <tmpl>				set default template for the output. Can use variables {{.Hash}}, {{.Author}}, {{.Commiter}}, {{.Message}} and {{.RepoName}}
			p	pull-all					pull all repos from list
			f	find-commits				find all commits containing msg in the commit message in the repo list
				--no-clipboard				don't copy the output to the clipboard
				-f <path>					output to the file in path
				-p							print the result in command line
				<msg>						pattern to find
			g	gui							open in GUI mode
	*/
	if len(os.Args) < 2 {
		webGui()
	} else {
		switch os.Args[1] {
		case "c":
			fallthrough
		case "config":
			openConfigFile()
		case "a":
			fallthrough
		case "add-repos":
			addRepos(os.Args[2])
		case "r":
			fallthrough
		case "remove-repos":
			removeRepos(os.Args[2])
		case "l":
			fallthrough
		case "repo-list":
			RepoTool()
		case "t":
			fallthrough
		case "template":
			setDefaultTemplate(os.Args[2])
		case "p":
			fallthrough
		case "pull-all":
			pullAll()
		case "f":
			fallthrough
		case "find-commits":
			FindCommits()
		}
	}
}

func openConfigFile() {
	switch runtime.GOOS {
	case "windows":
		execute("cmd", "/c", ".\\config.txt")
	case "linux":
		execute("edit", "./config.txt")
	default:
		fmt.Println("I don't know this system")
	}
}

func execute(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func addRepos(s string) {
	AddRepos(FindRepos(s))
	err := SaveConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func removeRepos(s string) {
	RemoveRepos(FindRepos(s))
	err := SaveConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func setDefaultTemplate(template string) {
	Config.DefaultTemplate = template
	SaveConfig()
}

func pullAll() {
	outputChan := make(chan string)
	defer close(outputChan)
	go func() {
		for out := range outputChan {
			fmt.Println(out)
		}
	}()
	PullAll(outputChan)
}
