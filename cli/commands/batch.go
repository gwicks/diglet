package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gwicks/diglet/compiler"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func batchAction(c *cli.Context) error {
	cmdArgs := c.Args()
	if len(cmdArgs) > 0 {
		inFile, err := os.Open(cmdArgs[0])
		if err != nil {
			log.Error(err)
			return err
		}
		defer inFile.Close()

		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			compilerOpts := compiler.BuildOptions{}
			batchline := scanner.Text()

			taskData := strings.Split(batchline, " ")
			fmt.Println(fmt.Sprintf("Processing %s into %s", taskData[0], taskData[1]))

			taskItemRes, err := compiler.CompileFile(taskData[0], compilerOpts)
			if err != nil {
				log.Error(err)
				return err
			}
			ioutil.WriteFile(taskData[1], []byte(taskItemRes), 0644)
		}
	} else {
		fmt.Println("Must specify a batch file to compile")
	}
	return nil
}

// BatchCommand Performs JSON compilation
func BatchCommand() cli.Command {
	return cli.Command{
		Name:    "batchfile",
		Aliases: []string{"b"},
		Usage:   "Compile each file within the specified batchfile",
		Action:  batchAction,
	}
}
