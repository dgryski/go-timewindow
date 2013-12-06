// Package timewindow provides counts for events in a sliding window of epochs
package timewindow

// Window is a sliding window of event counts.
type Window struct {
	counts  []int
	epoch   int64 // to match time.Now().Unix()
	headIdx int
	tailIdx int
	total   int
}

// New returns a sliding window starting at epoch0 with size seconds of history
func New(epoch0 int64, size int) *Window {

	w := &Window{
		counts:  make([]int, size),
		epoch:   epoch0,
		headIdx: 0,
		tailIdx: 1,
	}

	return w
}

// Add delta to the counter for epoch and adjust the window if necessary.
func (w *Window) Add(epoch int64, delta int) {

	// usual case -- update the present
	if epoch == w.epoch {
		w.total += delta
		w.counts[w.headIdx] += delta
		return
	}

	// common case -- advance our ring buffer
	if epoch > w.epoch {

		// FIXME(dgryski): we do too much work if zeroOut > len(count)
		zeroOut := int(epoch - w.epoch)
		for i := 0; i < zeroOut; i++ {
			w.total -= w.counts[w.tailIdx]
			w.counts[w.tailIdx] = 0
			w.tailIdx++
			if w.tailIdx >= len(w.counts) {
				w.tailIdx = 0
			}
		}

		w.headIdx += zeroOut
		for w.headIdx >= len(w.counts) {
			w.headIdx -= len(w.counts)

		}

		w.epoch = epoch
		w.total += delta
		w.counts[w.headIdx] += delta
		return
	}

	// less common -- update the past
	back := int(w.epoch - epoch)

	if back >= len(w.counts) {
		// too far in the past, ignore
		return
	}

	idx := w.headIdx - back

	if idx < 0 {
		// need to wrap around
		idx += len(w.counts)
	}

	w.total += delta
	w.counts[idx] += delta
}

// Total returns the sum of all counters in the window
func (w *Window) Total() int {
	return w.total
}

// Epoch returns most recent second for which data has been inserted
func (w *Window) Epoch() int64 {
	return w.epoch
}
