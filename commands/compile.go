package commands

import (
	"fmt"

	"github.com/gwicks/lincoln/utils"
	"github.com/urfave/cli"
)

func compileAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {
		resultJSON, _ := utils.ParseFileParent(cmdArgs[0])
		// resultJSON, _ := utils.ParseFileRefs(cmdArgs[0])
		fmt.Println(resultJSON)
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
