package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
			batchline := scanner.Text()

			taskData := strings.Split(batchline, " ")
			fmt.Println(fmt.Sprintf("Processing %s into %s", taskData[0], taskData[1]))
			doCompile(taskData[0], taskData[1])
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
