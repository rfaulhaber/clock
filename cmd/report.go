package cmd

import (
	"github.com/spf13/cobra"
	"time"
	"github.com/pkg/errors"
	"path/filepath"
	"github.com/rfaulhaber/clock/data"
	"io/ioutil"
	"bytes"
)

const parseTemplate = "01-02-2006"

var startDate, endDate string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunReport()

		if err != nil {
			stderr.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	weekAgo := time.Now().Add(-1 * 24 * 7 * time.Hour)

	reportCmd.Flags().StringVarP(&startDate, "start", "s", weekAgo.Format(parseTemplate), "start date for report")
	reportCmd.Flags().StringVarP(&endDate, "end", "e", time.Now().Format(parseTemplate), "end date for report")
}

func RunReport() error {
	parseStart, err := time.Parse(parseTemplate, startDate)

	if err != nil {
		return errors.Wrap(err, "could not parse start date")
	}

	parseEnd, err := time.Parse(parseTemplate, endDate)

	if err != nil {
		return errors.Wrap(err, "could not parse end date")
	}

	dates := fillRange(parseStart, parseEnd)

	dir := getDir()

	if dir == "" {
		return errors.New("no write dir specified")
	}

	var filenames []string

	for _, d := range dates {
		filenames = append(filenames, filepath.Join(dir, normalizeFile(getFileTimestamp(d, logTag), logTag)))
	}

	var records []*data.Record

	for _, fn := range filenames {
		b, err := ioutil.ReadFile(fn)

		if err != nil {
			return errors.Wrap(err, "could not read file")
		}

		table, err := data.Read(bytes.NewReader(b))

		if err != nil {
			return errors.Wrap(err, "could not parse table")
		}

		records = append(records, table.Records...)
	}

	duration := totalDuration(records)

	stdout.Println("entries: ", len(records))
	stdout.Println("duration: ", duration.String())

	return nil
}

type LogInfo struct {
	Entries uint
	TotalDuration time.Duration
}

func fillRange(start, stop time.Time) []time.Time {
	days := stop.Sub(start).Hours() / 24

	r := make([]time.Time, int(days))

	for i := 1; i < int(days) + 1; i++ {
		r[i - 1] = start.Add(24 * time.Duration(i) * time.Hour)
	}

	return r
}

func totalDuration(records []*data.Record) time.Duration {
	sum := time.Duration(0)

	for _, r := range records {
		sum += r.Duration()
	}

	return sum
}

