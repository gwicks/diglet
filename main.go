package main

import (
	"os"

	"github.com/gwicks/diglet/commands"
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
	app.Name = "diglet"
	app.Usage = "JSON Compiler and Schema Validator"
	app.Version = "0.2.5"
	app.Flags = []cli.Flag{}

	return app
}
