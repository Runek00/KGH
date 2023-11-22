package base

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func HandleInput(args []string) {
	switch os.Args[1] {
	case "c":
		fallthrough
	case "config":
		openConfigFile()
	case "a":
		fallthrough
	case "add-repos":
		addRepos(args[2])
	case "r":
		fallthrough
	case "remove-repos":
		removeRepos(args[2])
	case "l":
		fallthrough
	case "repo-list":
		RepoTool()
	case "t":
		fallthrough
	case "template":
		setDefaultTemplate(args[2])
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
