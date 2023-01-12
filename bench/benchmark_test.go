package bench

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunHistorySort(t *testing.T) {
	hist := RunHistory{
		{Date: time.Now().Add(-5 * time.Minute).Unix()},
		{Date: time.Now().Add(5 * time.Minute).Unix()},
		{Date: time.Now().Unix()},
	}
	sort.Sort(hist)
	if hist[0].Date < hist[1].Date || hist[1].Date < hist[2].Date {
		t.Error("expected most recent (largest date) to come first")
		t.Fail()
	}
}

func TestRunHistoryLatest(t *testing.T) {
	ts := time.Now().Add(5 * time.Minute).Unix()
	hist := RunHistory{
		{Date: time.Now().Add(-5 * time.Minute).Unix()},
		{Date: time.Now().Unix()},
		{Date: ts},
	}
	latest := hist.Latest()
	assert.Equal(t, ts, latest.Date)
}
