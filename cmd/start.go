package cmd

import (
	"github.com/rfaulhaber/clock/data"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "records start time",
	Long:  `Creates and records the start time of a log`,
	Run:   RunStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func RunStart(cmd *cobra.Command, args []string) {
	// TODO implement validation, cannot call start more than once, etc.

	dir := getDir()

	if dir == "" {
		stderr.Fatalln("write dir not specified")
	}

	currentDir := filepath.Join(dir, ".current")

	err := os.Mkdir(currentDir, 0700)

	if err != nil && !os.IsExist(err) {
		stderr.Fatalln(err)
	}

	f, err := os.OpenFile(filepath.Join(currentDir, normalizeCurrent(logTag)), os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		stderr.Fatalln(err)
	}

	startTime := time.Now()

	record := data.Record{Start: startTime}

	table := data.RecordTable{Tag: logTag, Records: []*data.Record{&record}}

	err = table.Write(f)

	f.Close()

	if err != nil {
		stderr.Fatalln(err)
	}

	stdout.Printf("started at: %d/%d/%d %02d:%02d:%02d\n", startTime.Month(), startTime.Day(), startTime.Year(), startTime.Hour(), startTime.Minute(), startTime.Second())
}
