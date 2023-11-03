package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Conf struct {
	Repos           map[string]Repo
	DefaultTemplate string
	OutputFilePath  string
}

var Config Conf

var templateSet = false
var reposSet = false
var filePathSet = false

var configTemplate = `{{ range $index, $value := .Repos }}repo: {{$value.Path}}::{{$value.Name}}::{{$value.Template}}
{{ end }}
template: {{.DefaultTemplate}}
filePath: {{.OutputFilePath}}`

func ReadConfig() error {
	file, err := os.Open("config.txt")
	templateSet = false
	reposSet = false
	filePathSet = false
	newConfig := Conf{make(map[string]Repo), "", ""}
	if err != nil {
		fmt.Println("Cannot open file. Starting a new config.")
		file, err = os.Create("config.txt")
		if err != nil {
			return err
		}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !reposSet && !strings.HasPrefix(line, "repo") && len(newConfig.Repos) > 0 {
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
			newConfig.Repos[repo.Path] = repo
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

func SaveConfig() error {
	tmpl := template.New("configTemplate")
	tmpl, err := tmpl.Parse(configTemplate)
	if err != nil {
		return err
	}
	file, err := os.Create("config.txt")
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, Config)
	if err != nil {
		return err
	}
	return nil
}

func AddRepos(repos []Repo) {
	for _, repo := range repos {
		Config.Repos[repo.Path] = repo
	}
}

func RemoveRepos(repos []Repo) {
	for _, repo := range repos {
		delete(Config.Repos, repo.Path)
	}
}

func GetRepos() []Repo {
	output := make([]Repo, 0, len(Config.Repos))
	for _, k := range Config.Repos {
		output = append(output, k)
	}
	return output
}
