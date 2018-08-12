package cmd

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
	"fmt"
)

var (
	stdout = log.New(os.Stdout, "clock: ", 0)
	stderr = log.New(os.Stderr, "clock: ", 0)
)

func getDir() string {
	if saveDir == "" {
		return viper.GetString("saveDir")
	} else {
		return saveDir
	}
}

func normalizeFile(filename string, tag string) string {
	if tag == "" {
		return filename
	} else {
		return tag + "-filename"
	}
}

func getFileTimestamp(t time.Time, tag string) string {
	dateTemplate := "%02d-%02d-%02d"
	if tag != "" {
		return fmt.Sprintf("%s-"+dateTemplate, tag, t.Year(), t.Month(), t.Day())
	} else {
		return fmt.Sprintf(dateTemplate, t.Year(), t.Month(), t.Day())
	}
}

