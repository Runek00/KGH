package main

import (
	"KGH/base"
	"KGH/web"
	"fmt"
	"os"
)

func main() {
	base.ReadConfig()
	halp := `	
		If you don't use any options, the program will be accessible at http://localhost:8079
		options:
			c	config						open config in vim / default editor
			a	add-repos <path>				add all repos in path to repo list
			r	remove-repos <path>				remove all repos in path from repo list
			l	repo-list					show repo list, interact with the repos
			t	template <tmpl>					set default template for the output. Can use variables {{.Hash}}, {{.Author}}, {{.Commiter}}, {{.Message}} and {{.RepoName}}
			p	pull-all					pull all repos from list
			f	find-commits					find all commits containing msg in the commit message in the repo list
				--no-clipboard					don't copy the output to the clipboard
				-f <path>					output to the file in path
				-p						print the result in command line
				<msg>						pattern to find`
	// g	gui							open in GUI mode

	if len(os.Args) < 2 {
		web.WebGui()
	} else if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Print(halp)
	} else {
		base.HandleInput(os.Args)
	}
}
