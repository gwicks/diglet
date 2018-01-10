package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dc0d/workerpool"
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

		jobs := make(chan func(), 10)

		myAppCtx, myAppCancel := context.WithCancel(context.Background())
		_ = myAppCancel

		pool, _ := workerpool.WithContext(myAppCtx, -1, jobs)

		mutex := &sync.Mutex{}
		for scanner.Scan() {
			batchline := scanner.Text()

			taskData := strings.Split(batchline, " ")

			if len(taskData) == 2 {
				fmt.Println(fmt.Sprintf("Processing %s into %s", taskData[0], taskData[1]))
				jobs <- func() {
					mutex.Lock()
					cerr := doCompile(taskData[0], taskData[1])
					mutex.Unlock()
					if cerr != nil {
						log.Error(cerr)
					}
				}

			}
		}
		pool.StopWait()
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
