package cmd

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/rfaulhaber/clock/data"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a current clock and finalizes the record.",
	Long:  `Stops the current clock and finalizes the record.`,
	Run:   RunStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func RunStop(cmd *cobra.Command, args []string) {
	// find the current record

	dir := getDir()

	if dir == "" {
		stderr.Fatalln("dir not specified")
	}

	currentDir := filepath.Join(dir, ".current")

	currentFile, err := ioutil.ReadFile(filepath.Join(currentDir, normalizeCurrent(logTag)))

	if err != nil {
		stderr.Fatalln(err)
	}

	var table data.RecordTable

	err = toml.Unmarshal(currentFile, &table)

	if err != nil {
		stderr.Fatalln(err)
	}

	table.Records[0].Stop = time.Now()

	record := table.Records[0]

	// TODO allow config file to specify date format?

	filename := filepath.Join(dir, getFileTimestamp(record.Start, logTag))

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	defer f.Close()

	if err != nil {
		stderr.Fatalln(err)
	}

	err = table.Write(f)

	if err != nil {
		stderr.Fatalln(err)
	}

	stdout.Println("duration: ", record.Duration())

	// TODO implement cleanup

	// append to today's record, divided by tag, and delete current file
}

func getFileTimestamp(t time.Time, tag string) string {
	dateTemplate := "%02d-%02d-%02d"
	if tag != "" {
		return fmt.Sprintf("%s-"+dateTemplate, tag, t.Year(), t.Month(), t.Day())
	} else {
		return fmt.Sprintf(dateTemplate, t.Year(), t.Month(), t.Day())
	}
}
