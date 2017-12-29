package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := makeGlobalApp()
	app.Commands = []cli.Command{}

	app.Run(os.Args)
}

func makeGlobalApp() *cli.App {
	app := cli.NewApp()
	app.Name = "diglett"
	// app.Before = config.Setup
	app.Usage = "JSON Compiler and Schema Validator"
	// app.Action = run
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debugging",
			EnvVar: "DEBUG",
		},
	}

	return app
}
