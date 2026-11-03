// Package timeline merges events from multiple sources into a single
// time-ordered stream suitable for rendering into a post-mortem skeleton.
package timeline

import (
	"sort"
	"time"
)

// Event is a single thing that happened: a Kubernetes event, a Prometheus
// alert transition, a git commit, or anything else a Source produces.
type Event struct {
	Time    time.Time
	Source  string // "kubectl", "prometheus", "git", etc.
	Summary string // one-line human-readable description
}

// Merge sorts events by timestamp, oldest first. Events with the same
// timestamp are stably ordered by Source then Summary so output is
// deterministic across runs.
func Merge(events []Event) []Event {
	out := make([]Event, len(events))
	copy(out, events)
	sort.SliceStable(out, func(i, j int) bool {
		if !out[i].Time.Equal(out[j].Time) {
			return out[i].Time.Before(out[j].Time)
		}
		if out[i].Source != out[j].Source {
			return out[i].Source < out[j].Source
		}
		return out[i].Summary < out[j].Summary
	})
	return out
}
