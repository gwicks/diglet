package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/gwicks/diglet/compiler"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var rootJSON interface{}

func compileAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {
		compileResult, err := compiler.CompileFile(cmdArgs[0])
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
	}
}
