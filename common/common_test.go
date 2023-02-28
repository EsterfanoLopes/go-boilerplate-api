package common_test

import (
	"go-boilerplate/common"
	"testing"
	"time"
)

func TestToTime(t *testing.T) {
	value := "2020-02-15T17:01:00-03:00"
	parsed, err := common.ToTime(value)
	if err != nil {
		t.Errorf("error parsing value to time %s", value)
		return
	}

	result := parsed.Format(time.RFC3339)
	if result != value {
		t.Errorf("wrong parsed time %s", result)
		return
	}
}

func TestQuotedStringBytes(t *testing.T) {
	result := common.QuotedStringBytes("VALUE")
	if string(result) != `"VALUE"` {
		t.Errorf("unexpected quoted string value %s", result)
	}
}
