package cmd

import (
	"github.com/rfaulhaber/clock/data"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
	"github.com/pkg/errors"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "records start time",
	Long:  `Creates and records the start time of a log`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunStart()

		if err != nil {
			stderr.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func RunStart() error {
	// TODO implement validation, cannot call start more than once, etc.

	dir := getDir()

	if dir == "" {
		return errors.New("write dir not specified")
	}

	_, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)

		if err != nil {
			return errors.Wrap(err, "mkdir if not exist")
		}
	}

	currentDir := filepath.Join(dir, ".current")

	err = os.Mkdir(currentDir, 0700)

	if err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "could not mkdir")
	}

	logFile := filepath.Join(currentDir, normalizeFile("current", logTag))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		return errors.Wrap(err, "could not open file: " + logFile)
	}

	startTime := time.Now()

	record := data.Record{Start: startTime}

	table := data.RecordTable{Records: []*data.Record{&record}}

	err = table.Write(f)

	f.Close()

	if err != nil {
		return errors.Wrap(err, "could not write table")
	}

	stdout.Printf("started at: %d/%d/%d %02d:%02d:%02d\n", startTime.Month(), startTime.Day(), startTime.Year(), startTime.Hour(), startTime.Minute(), startTime.Second())

	return nil
}
