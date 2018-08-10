package cmd

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	stdout = log.New(os.Stdout, "clock: ", 0)
	stderr = log.New(os.Stderr, "clock: ", 0)
)

func getDir() string {
	if saveDir == "" {
		return viper.GetString("saveDir")
	}

	return saveDir
}

func normalizeCurrent(tag string) string {
	if tag == "" {
		return "current"
	} else {
		return logTag + "-current"
	}
}
