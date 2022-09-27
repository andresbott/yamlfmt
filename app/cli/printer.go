package cli

import "fmt"

type printer struct {
	quiet   bool
	verbose bool
}

func (p printer) print(msg string) {
	if p.quiet {
		return
	}
	fmt.Println(msg)
}

func (p printer) printVerbose(msg string) {
	if p.quiet {
		return
	}
	if p.verbose {
		fmt.Println(msg)
	}
}
