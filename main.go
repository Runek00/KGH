package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	/*
		flags:
			c	config						open config in vim / default editor
			a	add-repos <path>			add all repos in path to repo list
			r	remove-repos <path>			remove all repos in path from repo list
			l	repo-list					show repo list
			t	template <tmpl>				set template for the output. Can use variables "repoName" and "hash"
			p	pull-all					pull all repos from list
			f	find-commits <msg>			find all commits containing msg in the commit message in the repo list
					--no-clipboard			don't copy the output to the clipboard
				-F	--file, --File [path]	output to the file (sets path if present)
				-r	--print					print the result in command line
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
			printRepoList()
		case "t":
			fallthrough
		case "template":
			setCommitTemplate(os.Args[2])
		case "p":
			fallthrough
		case "pull-all":
			pullAll()
		case "f":
			fallthrough
		case "find-commits":
			findCommits(os.Args[2])
		}
	}
	fmt.Println(os.Args[1:])
}

func openConfigFile() {
	switch runtime.GOOS {
	case "windows":
		exec.Command(".\\config.txt")
	case "linux":
		exec.Command("edit", "./config.txt")
	default:
		fmt.Println("I don't know this system")
	}
}

func addRepos(s string) {
	panic("unimplemented")
}

func removeRepos(s string) {
	panic("unimplemented")
}

func printRepoList() {
	panic("unimplemented")
}

func setCommitTemplate(s string) {
	panic("unimplemented")
}

func pullAll() {
	panic("unimplemented")
}

func findCommits(s string) {
	panic("unimplemented")
}
