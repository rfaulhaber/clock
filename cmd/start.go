package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"log"
	"github.com/rfaulhaber/clock/data"
	"time"
	"github.com/spf13/viper"
	"os/user"
	"fmt"
)

var (
	saveDir string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "records start time",
	Long: `Creates and records the start time of a log`,
	Run: RunStart,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVarP(&saveDir, "save-dir", "s", "", "directory to save this log to")
}

func RunStart(cmd *cobra.Command, args []string) {
	dir := getDir()

	currentDir := filepath.Join(dir, ".current")

	err := os.Mkdir(currentDir, 0700)

	u, _ := user.Current()

	fmt.Println(u.Username)

	if err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(filepath.Join(currentDir, logTag + "current"), os.O_CREATE | os.O_WRONLY, 0700)

	if err != nil {
		log.Fatalln(err)
	}

	record := data.Record{time.Now(), time.Time{}}

	table := data.RecordTable{Tag: logTag, Records: []*data.Record{&record}}

	err = table.Write(f)

	if err != nil {
		log.Fatalln(err)
	}
}

func getDir() string {
	if saveDir == "" {
		return viper.GetString("saveDir")
	}

	return saveDir
}
