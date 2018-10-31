package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/gwicks/diglet/compiler"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func compileAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {
		compilerOpts := compiler.BuildOptions{
			SkipResolve:   c.Bool("skip-resolve"),
			SkipParenting: c.Bool("skip-parenting"),
			SkipValidate:  c.Bool("skip-validation"),
		}
		compileResult, err := compiler.CompileFile(cmdArgs[0], compilerOpts)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(cmdArgs) == 2 {
			ioutil.WriteFile(cmdArgs[1], []byte(compileResult), 0644)
		} else {
			fmt.Println(compileResult)
		}
	} else {
		fmt.Println("Must specify a JSON file to compile")
	}
	return nil
}

// CompileCommand Performs JSON compilation
func CompileCommand() cli.Command {
	return cli.Command{
		Name:    "compile",
		Aliases: []string{"c"},
		Usage:   "Compile the source file and resolve it's dependencies",
		Action:  compileAction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "skip-resolve",
			},
			cli.BoolFlag{
				Name: "skip-parenting",
			},
			cli.BoolFlag{
				Name: "skip-validation",
			},
		},
	}
}
