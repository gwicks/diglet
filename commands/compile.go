package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gwicks/lincoln/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var rootJSON interface{}

func compileAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {

		fd, err := ioutil.ReadFile(cmdArgs[0])
		if err != nil {
			log.Error(err)
			return err
		}

		json.Unmarshal(fd, &rootJSON)

		resultJSON, _ := utils.ParseFileRefs(cmdArgs[0], rootJSON)

		if resultJSONObj, ok := resultJSON.(map[string]interface{}); ok {
			finalJSON, _ := utils.ParseFileParent(cmdArgs[0], resultJSONObj)

			validatedJSON, verr := utils.ParseFileSchema(cmdArgs[0], finalJSON)

			if verr != nil {
				log.Error(verr)
			} else {
				resultStr, _ := json.MarshalIndent(validatedJSON, "", "    ")

				fmt.Println(string(resultStr))
			}
		} else {
			resultStr, _ := json.MarshalIndent(resultJSON, "", "    ")

			fmt.Println(string(resultStr))
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
