package bench

import (
	"sort"
	"testing"
	"time"
)

func TestRunHistory_Sort(t *testing.T) {
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
