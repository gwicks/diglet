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

func doCompile(filePath string, targetPath string) error {
	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return err
	}

	json.Unmarshal(fd, &rootJSON)

	resultJSON, _ := utils.ParseFileRefs(filePath, rootJSON)

	if resultJSONObj, ok := resultJSON.(map[string]interface{}); ok {
		finalJSON, _ := utils.ParseFileParent(filePath, resultJSONObj)

		validatedJSON, verr := utils.ParseFileSchema(filePath, finalJSON)
		if verr != nil {
			log.Error(verr)
		} else {
			resultStr, _ := json.MarshalIndent(validatedJSON, "", "    ")

			if len(targetPath) > 0 {
				ioutil.WriteFile(targetPath, resultStr, 0644)
			} else {
				fmt.Println(string(resultStr))
			}
		}
	} else {
		resultStr, _ := json.MarshalIndent(resultJSON, "", "    ")

		fmt.Println(string(resultStr))
	}
	return nil
}

func compileAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {
		if len(cmdArgs) == 2 {
			doCompile(cmdArgs[0], cmdArgs[1])
		} else {
			doCompile(cmdArgs[0], "")
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
