// Package sources contains ingesters that produce timeline.Event values
// from external systems (kubectl, Prometheus, git).
package sources

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

// kubectlEvent mirrors the subset of core/v1.Event we need. We parse
// kubectl's own JSON output rather than taking a direct client-go
// dependency - this tool runs wherever kubectl runs, which is a wider
// surface than compiled client-go would support.
type kubectlEvent struct {
	LastTimestamp  string `json:"lastTimestamp"`
	EventTime      string `json:"eventTime"`
	FirstTimestamp string `json:"firstTimestamp"`
	Type           string `json:"type"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	InvolvedObject struct {
		Kind      string `json:"kind"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"involvedObject"`
}

type kubectlResponse struct {
	Items []kubectlEvent `json:"items"`
}

// FetchKubectlEvents runs `kubectl get events -o json` for the given
// namespace (empty = all namespaces) and returns events whose timestamp
// falls inside [since, until].
func FetchKubectlEvents(namespace string, since, until time.Time) ([]timeline.Event, error) {
	args := []string{"get", "events", "-o", "json"}
	if namespace != "" {
		args = append([]string{"-n", namespace}, args...)
	} else {
		args = append(args, "--all-namespaces")
	}

	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("kubectl exited %d: %s",
				ee.ExitCode(), strings.TrimSpace(string(ee.Stderr)))
		}
		return nil, fmt.Errorf("run kubectl: %w", err)
	}

	var resp kubectlResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("parse kubectl json: %w", err)
	}

	var events []timeline.Event
	for _, e := range resp.Items {
		ts := bestTimestamp(e)
		if ts.IsZero() {
			continue
		}
		if ts.Before(since) || ts.After(until) {
			continue
		}

		obj := fmt.Sprintf("%s/%s", strings.ToLower(e.InvolvedObject.Kind), e.InvolvedObject.Name)
		if e.InvolvedObject.Namespace != "" {
			obj = fmt.Sprintf("%s/%s", e.InvolvedObject.Namespace, obj)
		}

		summary := fmt.Sprintf("[%s] %s %s: %s",
			e.Type, obj, e.Reason, truncate(e.Message, 120))

		events = append(events, timeline.Event{
			Time:    ts.Local(),
			Source:  "kubectl",
			Summary: summary,
		})
	}
	return events, nil
}

// bestTimestamp picks the most recent populated timestamp from the Event.
// kubectl Events may populate lastTimestamp, eventTime, or firstTimestamp
// depending on the version and event series state.
func bestTimestamp(e kubectlEvent) time.Time {
	for _, s := range []string{e.LastTimestamp, e.EventTime, e.FirstTimestamp} {
		if s == "" {
			continue
		}
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
