package sources

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// Helper: build a journalctl-style JSON line with a given realtime
// timestamp (microseconds since epoch) plus arbitrary other fields.
func journalLine(usec int64, kv ...string) string {
	parts := []string{
		`"__REALTIME_TIMESTAMP":"` + strconv.FormatInt(usec, 10) + `"`,
	}
	for i := 0; i+1 < len(kv); i += 2 {
		parts = append(parts, `"`+kv[i]+`":"`+kv[i+1]+`"`)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func TestParseJournal_HappyPath(t *testing.T) {
	// 2026-10-29 03:15:00 UTC in microseconds since epoch
	ts1 := int64(1793325300) * 1_000_000
	ts2 := ts1 + 5_000_000
	out := strings.Join([]string{
		journalLine(ts1, "PRIORITY", "3", "_SYSTEMD_UNIT", "k3s.service", "MESSAGE", "failed to run kube-proxy"),
		journalLine(ts2, "PRIORITY", "4", "_SYSTEMD_UNIT", "containerd.service", "MESSAGE", "CNI config missing"),
	}, "\n") + "\n"

	since := time.Unix(0, (ts1-1_000_000)*1000)
	until := time.Unix(0, (ts2+1_000_000)*1000)
	events, err := ParseJournal([]byte(out), since, until)
	if err != nil {
		t.Fatalf("ParseJournal returned error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Source != "journal" {
		t.Errorf("expected source 'journal', got %q", events[0].Source)
	}
	wantPrefix := "[err] k3s.service:"
	if !strings.HasPrefix(events[0].Summary, wantPrefix) {
		t.Errorf("event[0] summary missing prefix %q: %q", wantPrefix, events[0].Summary)
	}
}

func TestParseJournal_MultilineMessageFlattened(t *testing.T) {
	ts := int64(1793325300) * 1_000_000
	// JSON-encode a newline inside the MESSAGE value
	line := `{"__REALTIME_TIMESTAMP":"` + strconv.FormatInt(ts, 10) + `","PRIORITY":"4","_SYSTEMD_UNIT":"containerd.service","MESSAGE":"line one\nline two"}`
	since := time.Unix(0, (ts-1_000_000)*1000)
	until := time.Unix(0, (ts+1_000_000)*1000)
	events, err := ParseJournal([]byte(line+"\n"), since, until)
	if err != nil {
		t.Fatalf("ParseJournal returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if strings.Contains(events[0].Summary, "\n") {
		t.Errorf("summary should be single-line, got %q", events[0].Summary)
	}
}

func TestParseJournal_WindowFilter(t *testing.T) {
	ts := int64(1793325300) * 1_000_000
	line := journalLine(ts, "PRIORITY", "6", "MESSAGE", "hello")
	// Window starts AFTER the event timestamp
	since := time.Unix(0, (ts+10_000_000)*1000)
	until := time.Unix(0, (ts+20_000_000)*1000)
	events, err := ParseJournal([]byte(line+"\n"), since, until)
	if err != nil {
		t.Fatalf("ParseJournal error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events (out of window), got %d", len(events))
	}
}

func TestParseJournal_SkipsUnparseableTimestamp(t *testing.T) {
	bad := `{"PRIORITY":"4","MESSAGE":"no timestamp"}`
	ts := int64(1793325300) * 1_000_000
	good := journalLine(ts, "PRIORITY", "4", "MESSAGE", "ok")
	since := time.Unix(0, (ts-1_000_000)*1000)
	until := time.Unix(0, (ts+1_000_000)*1000)
	events, err := ParseJournal([]byte(bad+"\n"+good+"\n"), since, until)
	if err != nil {
		t.Fatalf("ParseJournal error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event (good survives, bad silently skipped), got %d", len(events))
	}
}

func TestParseJournal_InvalidJSONReturnsError(t *testing.T) {
	_, err := ParseJournal([]byte("not json\n"), time.Now(), time.Now())
	if err == nil {
		t.Fatal("expected error on invalid JSON, got nil")
	}
}

func TestParseJournal_FallbackToSyslogIdentifier(t *testing.T) {
	ts := int64(1793325300) * 1_000_000
	// _SYSTEMD_UNIT absent, SYSLOG_IDENTIFIER set
	line := `{"__REALTIME_TIMESTAMP":"` + strconv.FormatInt(ts, 10) + `","PRIORITY":"4","SYSLOG_IDENTIFIER":"kernel","MESSAGE":"BPF program loaded"}`
	since := time.Unix(0, (ts-1_000_000)*1000)
	until := time.Unix(0, (ts+1_000_000)*1000)
	events, err := ParseJournal([]byte(line+"\n"), since, until)
	if err != nil {
		t.Fatalf("ParseJournal error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if !strings.Contains(events[0].Summary, "kernel:") {
		t.Errorf("expected fallback to SYSLOG_IDENTIFIER=kernel, got %q", events[0].Summary)
	}
}

func TestParseJournal_UnknownPriorityNoLabel(t *testing.T) {
	ts := int64(1793325300) * 1_000_000
	line := journalLine(ts, "PRIORITY", "99", "_SYSTEMD_UNIT", "foo.service", "MESSAGE", "weird")
	since := time.Unix(0, (ts-1_000_000)*1000)
	until := time.Unix(0, (ts+1_000_000)*1000)
	events, err := ParseJournal([]byte(line+"\n"), since, until)
	if err != nil {
		t.Fatalf("ParseJournal error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if strings.HasPrefix(events[0].Summary, "[") {
		t.Errorf("unknown priority should drop the [label] prefix, got %q", events[0].Summary)
	}
}
