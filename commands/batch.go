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
			fmt.Println("ERROR ON BATCH")
			log.Error(err)
			return err
		}
		defer inFile.Close()

		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			batchline := scanner.Text()

			taskData := strings.Split(batchline, " ")

			if len(taskData) == 2 {
				fmt.Print("Processing ")
				fmt.Print(taskData[0])
				fmt.Print(" into ")
				fmt.Println(taskData[1])
				cerr := doCompile(taskData[0], taskData[1])
				if cerr != nil {
					return cerr
				}
			}
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
