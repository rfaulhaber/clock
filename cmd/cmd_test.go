package cmd

import (
	"bytes"
	"github.com/pelletier/go-toml"
	"github.com/rfaulhaber/clock/data"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStart(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "clocktest")

	defer os.RemoveAll(tmpDir)

	if err != nil {
		panic(err)
	}

	var outBuf bytes.Buffer

	stdout = log.New(&outBuf, "", 0)

	saveDir = tmpDir

	runStartNoArg()

	dirStat, err := os.Stat(filepath.Join(tmpDir, ".current"))

	assert.NoError(t, err)
	assert.True(t, dirStat.IsDir())

	currentFile := filepath.Join(tmpDir, ".current", "current")

	fileStat, err := os.Stat(currentFile)

	assert.NoError(t, err)
	assert.True(t, fileStat.Size() > 0)

	var fileData data.RecordTable

	b, err := ioutil.ReadFile(currentFile)

	assert.NoError(t, err)

	err = toml.Unmarshal(b, &fileData)

	assert.NoError(t, err)
	assert.False(t, fileData.Records[0].Start.IsZero())
	assert.True(t, fileData.Records[0].Stop.IsZero())
	assert.Equal(t, 1, len(fileData.Records))
	assert.True(t, strings.Index(outBuf.String(), "started at:") > -1)
}

func TestRunStopWithoutDir(t *testing.T) {
	err := RunStop()

	assert.Error(t, err)
	assert.Equal(t, "Cannot find any records. You must run `start` before calling stop.", err.Error())
}

func TestRunStopBeforeStart(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "clocktest")

	defer os.RemoveAll(tmpDir)

	if err != nil {
		panic(err)
	}

	saveDir = tmpDir

	err = RunStop()

	assert.Error(t, err)
	assert.Equal(t, "Cannot find any records. You must run `start` before calling stop.", err.Error())
}

func runStartNoArg() {
	RunStart(nil, []string{})
}
