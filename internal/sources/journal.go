// Package sources contains ingesters that produce timeline.Event values
// from external systems (kubectl, Prometheus, git, journalctl).
package sources

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

// journalEntry mirrors the fields we care about from `journalctl -o json`.
// journalctl emits a stream of JSON objects (one per line, not an array).
// Timestamps come as microseconds-since-epoch encoded as strings.
type journalEntry struct {
	RealtimeTimestamp string `json:"__REALTIME_TIMESTAMP"`
	Priority          string `json:"PRIORITY"`
	SystemdUnit       string `json:"_SYSTEMD_UNIT"`
	SyslogIdentifier  string `json:"SYSLOG_IDENTIFIER"`
	Message           string `json:"MESSAGE"`
}

// JournalOptions controls optional filters passed to journalctl. All
// fields are optional - a zero value means "no filter on this axis".
type JournalOptions struct {
	// Units limits the query to one or more systemd units. Each entry is
	// passed as a separate --unit=X flag to journalctl.
	Units []string
	// MaxPriority filters to entries at or below this numeric priority
	// (0=emerg..7=debug, per RFC 5424). Negative value = no filter.
	MaxPriority int
}

// FetchJournalEntries shells out to `journalctl -o json` for the given
// time window and returns events. Pass a zero JournalOptions for no
// filters. Entries without a parseable __REALTIME_TIMESTAMP are
// silently skipped so one broken line does not drop the whole window.
func FetchJournalEntries(since, until time.Time, opts JournalOptions) ([]timeline.Event, error) {
	args := []string{
		"-o", "json", "--no-pager",
		"--since", since.Format("2006-01-02 15:04:05"),
		"--until", until.Format("2006-01-02 15:04:05"),
	}
	for _, u := range opts.Units {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		args = append(args, "--unit="+u)
	}
	if opts.MaxPriority >= 0 && opts.MaxPriority <= 7 {
		args = append(args, "--priority="+strconv.Itoa(opts.MaxPriority))
	}

	cmd := exec.Command("journalctl", args...)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("journalctl %v: %w: %s",
				args, err, strings.TrimSpace(string(ee.Stderr)))
		}
		return nil, fmt.Errorf("journalctl %v: %w", args, err)
	}
	return ParseJournal(out, since, until)
}

// ParseJournal is exported for unit testing without a real systemd.
// It takes the raw output of `journalctl -o json` (one JSON object per
// line) and returns timeline events within [since, until].
func ParseJournal(out []byte, since, until time.Time) ([]timeline.Event, error) {
	var events []timeline.Event
	for i, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var e journalEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("journal line %d: %w", i+1, err)
		}
		ts, err := parseJournalTime(e.RealtimeTimestamp)
		if err != nil {
			// Skip lines without a parseable timestamp rather than aborting;
			// a single broken entry should not drop the whole window.
			continue
		}
		if ts.Before(since) || ts.After(until) {
			continue
		}
		events = append(events, timeline.Event{
			Time:    ts,
			Source:  "journal",
			Summary: journalSummary(e),
		})
	}
	return events, nil
}

// parseJournalTime converts journalctl's "microseconds since epoch as
// string" field into a time.Time in local timezone.
func parseJournalTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	usec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, usec*1000).In(time.Local), nil
}

// journalSummary renders a one-line human summary for the timeline table.
// Falls back to SYSLOG_IDENTIFIER when _SYSTEMD_UNIT is absent (kernel
// messages, early-boot, etc.).
func journalSummary(e journalEntry) string {
	unit := e.SystemdUnit
	if unit == "" {
		unit = e.SyslogIdentifier
	}
	if unit == "" {
		unit = "-"
	}
	msg := strings.ReplaceAll(e.Message, "\n", " ")
	if lvl := priorityLabel(e.Priority); lvl != "" {
		return fmt.Sprintf("[%s] %s: %s", lvl, unit, msg)
	}
	return fmt.Sprintf("%s: %s", unit, msg)
}

// priorityLabel maps the numeric syslog priority to a short label.
// Returns "" for unrecognized values so the summary stays readable.
func priorityLabel(p string) string {
	switch p {
	case "0":
		return "emerg"
	case "1":
		return "alert"
	case "2":
		return "crit"
	case "3":
		return "err"
	case "4":
		return "warn"
	case "5":
		return "notice"
	case "6":
		return "info"
	case "7":
		return "debug"
	default:
		return ""
	}
}
