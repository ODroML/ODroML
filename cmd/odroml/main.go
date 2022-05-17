// Tealang runtime REPL tool.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/odroml/odroml/pkg/repl"
	"gopkg.in/urfave/cli.v1"
)

var (
	interactiveMode, graphvizMode bool
)

func runInteractiveShell(c *cli.Context) error {
	env := repl.New(repl.Config{OutputGraph: c.GlobalBool("graph")})
	reader := bufio.NewReader(os.Stdin)
	for env.Active {
		fmt.Fprint(os.Stdout, "> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			env.Stop()
		} else {
			output, err := env.Interpret(strings.TrimRight(input, "\n"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else if output != "" {
				fmt.Fprintln(os.Stdout, output)
			}
		}
	}
	fmt.Fprintln(os.Stdout)
	return nil
}

func executeProgramFile(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("required filename")
	}
	return repl.New(repl.Config{OutputGraph: c.GlobalBool("graph")}).Load(c.Args()[0])
}

func main() {
	app := cli.NewApp()
	app.Name = "odroml"
	app.Usage = "interactive runtime environment"
	app.Copyright = "2022 ODroML. All rights reserved."
	app.Author = "The ODroMLdas"
	app.Version = "v0.1.0-dev"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "graph",
			Usage:  "Show graphviz output instead of program result",
			Hidden: false,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "shell",
			Usage:  "Run an interactive interpreter",
			Action: runInteractiveShell,
		},
		{
			Name:   "run",
			Usage:  "Execute a program file",
			Action: executeProgramFile,
		},
	}
	app.Run(os.Args)
}
