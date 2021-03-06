package cmd

import (
	"bytes"
	"github.com/pelletier/go-toml"
	"github.com/rfaulhaber/clock/internal/record"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStart(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "clocktest")

	defer os.RemoveAll(tmpDir)

	if err != nil {
		panic(err)
	}

	var outBuf bytes.Buffer

	stdout = log.New(&outBuf, "", 0)

	saveDir = tmpDir

	err = RunStart()
	assert.NoError(t, err)

	dirStat, err := os.Stat(filepath.Join(tmpDir, ".current"))

	assert.NoError(t, err)
	assert.True(t, dirStat.IsDir())

	currentFile := filepath.Join(tmpDir, ".current", "current")

	fileStat, err := os.Stat(currentFile)

	assert.NoError(t, err)
	assert.True(t, fileStat.Size() > 0)

	var fileData record.RecordTable

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

func TestRunStop(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "clocktest")

	defer os.RemoveAll(tmpDir)

	if err != nil {
		panic(err)
	}

	var outBuf bytes.Buffer

	stdout = log.New(&outBuf, "", 0)

	saveDir = tmpDir

	err = RunStart()

	assert.NoError(t, err)

	err = RunStop()

	assert.NoError(t, err)

	infos, err := ioutil.ReadDir(tmpDir)

	assert.NoError(t, err)

	for _, info := range infos {
		if info.Name() == ".current" {
			continue
		}

		assert.False(t, info.IsDir())
	}

	var filename string

	for _, info := range infos {
		if info.Name() != ".current" {
			filename = info.Name()
		}
	}

	b, err := ioutil.ReadFile(filepath.Join(tmpDir, filename))

	assert.NoError(t, err)

	var table record.RecordTable

	err = toml.Unmarshal(b, &table)

	assert.NoError(t, err)

	assert.False(t, table.Records[0].Start.IsZero())
	assert.False(t, table.Records[0].Stop.IsZero())
}
