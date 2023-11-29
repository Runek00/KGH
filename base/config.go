package base

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Conf struct {
	Repos           map[string]Repo
	DefaultTemplate string
	OutputFilePath  string
}

var Config Conf

func ReadConfig() error {
	file, err := os.Open("config.json")
	newConfig := Conf{make(map[string]Repo), "", ""}
	if err != nil {
		fmt.Println("Cannot open file. Starting a new config.")
		file, err = os.Create("config.txt")
		if err != nil {
			return err
		}
	}
	defer file.Close()
	byteFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Cannot read file")
		return err
	}
	err = json.Unmarshal(byteFile, &newConfig)
	if err != nil {
		fmt.Println("Cannot parse file")
		return err
	}
	Config = newConfig
	return nil
}

func SaveConfig() error {
	file, err := os.Create("config.json")
	writer := bufio.NewWriter(file)
	if err != nil {
		fmt.Println("Cannot even create a writer")
		return err
	}
	byteFile, err := json.Marshal(Config)
	if err != nil {
		fmt.Println("Config marshalling failed")
		return err
	}
	_, err = writer.WriteString(string(byteFile))
	if err != nil {
		fmt.Println("Cannot write file")
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
