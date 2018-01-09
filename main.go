package main

import (
	"os"

	"github.com/gwicks/lincoln/commands"
	"github.com/urfave/cli"
)

func main() {
	app := makeGlobalApp()
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.CompileCommand())
	app.Commands = append(app.Commands, commands.BatchCommand())

	app.Run(os.Args)
}

func makeGlobalApp() *cli.App {
	app := cli.NewApp()
	app.Name = "lincoln"
	// app.Before = config.Setup
	app.Usage = "JSON Compiler and Schema Validator"
	// app.Action = run
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{}

	return app
}
