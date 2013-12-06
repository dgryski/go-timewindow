package timewindow

import (
	"testing"
)

func TestWindow(t *testing.T) {

	var tests = []struct {
		addEpoch int64
		delta    int
		total    int
		epoch    int64
	}{
		{100, 1, 1, 100},
		{100, 2, 3, 100},
		{101, 2, 5, 101},
		{105, 1, 6, 105},
		{90, 1, 6, 105},
		{95, 1, 6, 105},
		{96, 1, 7, 105},
		{120, 1, 1, 120},
		{130, 0, 0, 130},
	}

	w := New(100, 10)
	for i, tt := range tests {
		w.Add(tt.addEpoch, tt.delta)

		if total := w.Total(); total != tt.total {
			t.Errorf("failed test %d: total=%d wanted %d\n", i, total, tt.total)
		}

		if epoch := w.Epoch(); epoch != tt.epoch {
			t.Errorf("failed test %d: epoch=%d wanted %d\n", i, epoch, tt.epoch)
		}
	}
}
