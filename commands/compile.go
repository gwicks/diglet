package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gwicks/lincoln/utils"
	"github.com/urfave/cli"
)

var rootJSON map[string]interface{}

func compileAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {
		rootJSON = make(map[string]interface{})

		fd, err := ioutil.ReadFile(cmdArgs[0])
		if err != nil {
			return err
		}

		json.Unmarshal(fd, &rootJSON)

		resultJSON, _ := utils.ParseFileRefs(cmdArgs[0], rootJSON)

		finalJSON, _ := utils.ParseFileParent(cmdArgs[0], resultJSON)

		resultStr, _ := json.MarshalIndent(finalJSON, "", "    ")

		fmt.Println(string(resultStr))
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
