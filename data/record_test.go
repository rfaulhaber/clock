package data

import (
	"testing"
	"time"
		"github.com/stretchr/testify/assert"
			"bytes"
	"github.com/BurntSushi/toml"
	)

func TestNewRecord(t *testing.T) {
	testStart := time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC)
	testEnd := testStart.Add(8 * time.Hour)

	result := NewRecord(testStart, testEnd)

	assert.Equal(t, result.Start, testStart)
	assert.Equal(t, result.Stop, testEnd)
}

func TestRecord_Duration(t *testing.T) {
	testCase := []struct{
		Input *Record
		Expected time.Duration
	}{
		{
			NewRecord(time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC), time.Date(1992, time.December, 27, 18, 15, 15, 0, time.UTC)),
			3 * time.Hour,
		},
		{
			NewRecord(time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC), time.Time{}),
			0,
		},
	}

	for _, tc := range testCase {
		result := tc.Input.Duration()
		assert.Equal(t, result, tc.Expected)
	}
}

func TestRead(t *testing.T) {
	testStart := time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC)
	testEnd := testStart.Add(8 * time.Hour)

	var buf bytes.Buffer

	toml.NewEncoder(&buf).Encode(RecordTable{
		Records: []*Record{NewRecord(testStart, testEnd)},
	})

	result, err := Read(&buf)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Records))
	assert.Equal(t, result.Records[0].Start, testStart)
	assert.Equal(t, result.Records[0].Stop, testEnd)

}

func TestRecordTable_Write(t *testing.T) {
	testStart := time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC)
	testEnd := testStart.Add(8 * time.Hour)

	testTable := RecordTable{"test tag", []*Record{NewRecord(testStart, testEnd)}}

	var buf bytes.Buffer

	err := testTable.Write(&buf)

	assert.NoError(t, err)

	expectedStr := `Tag = "test tag"

[[Records]]
  start = 1992-12-27T15:15:15Z
  stop = 1992-12-27T23:15:15Z
`

	assert.Equal(t, expectedStr, buf.String())
}

func TestRecordTable_Add(t *testing.T) {
	testStart := time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC)
	testEnd := testStart.Add(8 * time.Hour)
	testTable := RecordTable{"test tag", []*Record{NewRecord(testStart, testEnd)}}

	testTable.Add(NewRecord(testStart, testEnd))

	assert.Equal(t, 2, len(testTable.Records))
}

func TestRecordTable_Update(t *testing.T) {
	testStart := time.Date(1992, time.December, 27, 15, 15, 15, 0, time.UTC)
	testEnd := testStart.Add(8 * time.Hour)
	testTable := RecordTable{"test tag", []*Record{NewRecord(testStart, time.Time{})}}

	testTable.Update(testEnd)

	assert.NotEmpty(t, testTable.Records[0].Stop)
	assert.Equal(t, testEnd.String(), testTable.Records[0].Stop.String())
}