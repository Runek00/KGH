package main

import (
	"fmt"
	"os"
)

func main() {
	/*
		flags:
			c	--config				open config in vim / default editor
			a	--add-repos <path>		add all repos in dir to repo list
			r	--remove-repos <path>	remove all repos in dir from repo list
			l	--repo-list				show repo list
			t	--template <tmpl>		set template for the output. Can use variables "repoName" and "hash"
			p	--pull-all				pull all repos from list
			f	--find-commits <msg>	find all commits containing msg in the commit message in the repo list
					--no-clipboard			don't copy the output to the clipboard
				-F	--file, --File [path]	output to the file (sets path if present)
				-r	--print					print the result in command line
			-g	-gui					open in GUI mode
	*/
	fmt.Println(os.Args[1:])
}
