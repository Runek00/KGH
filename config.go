package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Conf struct {
	repos           map[string]Repo
	DefaultTemplate string
	OutputFilePath  string
}

var Config Conf

var templateSet = false
var reposSet = false
var filePathSet = false

func ReadConfig() error {
	file, err := os.Open("config.txt")
	templateSet = false
	reposSet = false
	filePathSet = false
	newConfig := Conf{}
	if err != nil {
		fmt.Println("Cannot open file. Starting a new config.")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !reposSet && !strings.HasPrefix(line, "repo") && len(newConfig.repos) > 0 {
			reposSet = true
		}
		if strings.HasPrefix(line, "repo") {
			if reposSet {
				return errors.New("repos are not in a single block in config")
			}
			repoLine := strings.Split(strings.Split(line, "repo: ")[1], "::")
			if len(repoLine) < 2 || len(repoLine) > 3 {
				return errors.New("Error reading repo " + line)
			}
			var repo Repo
			if len(repoLine) == 2 {
				repo = Repo{repoLine[0], repoLine[1], ""}
			} else {
				repo = Repo{repoLine[0], repoLine[1], repoLine[2]}
			}
			newConfig.repos[repo.path] = repo
		}
		if strings.HasPrefix(line, "template") {
			if templateSet {
				return errors.New("more then one default template in config")
			}
			newConfig.DefaultTemplate = strings.Split(line, "template: ")[1]
			templateSet = true
		}
		if strings.HasPrefix(line, "filePath") {
			if filePathSet {
				return errors.New("more than one output file path in config")
			}
			newConfig.OutputFilePath = strings.Split(line, "filePath: ")[1]
			filePathSet = true
		}
	}
	Config = newConfig
	return nil
}

func SaveConfig() {
	// TODO config template -> config to tamplate -> to file
}

func AddRepos(repos []Repo) {
	for _, repo := range repos {
		Config.repos[repo.path] = repo
	}
}

func RemoveRepos(repos []Repo) {
	for _, repo := range repos {
		delete(Config.repos, repo.path)
	}
}

func GetRepos() []Repo {
	output := make([]Repo, 0, len(Config.repos))
	for _, k := range Config.repos {
		output = append(output, k)
	}
	return output
}
