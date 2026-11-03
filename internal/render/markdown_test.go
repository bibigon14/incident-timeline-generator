package render

import (
	"strings"
	"testing"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

func TestMarkdownRendersHeader(t *testing.T) {
	since := time.Date(2026, 10, 27, 22, 0, 0, 0, time.UTC)
	until := time.Date(2026, 10, 28, 0, 30, 0, 0, time.UTC)

	md := Markdown(Input{
		Title:    "vmalert CrashLoopBackOff",
		Severity: "SEV-3",
		Author:   "Dmitry Stepanov",
		Since:    since,
		Until:    until,
		Events:   nil,
	})

	for _, want := range []string{
		"# Incident: vmalert CrashLoopBackOff",
		"**Severity:** SEV-3",
		"**Author:** Dmitry Stepanov",
		"## Summary",
		"## Impact",
		"## Timeline",
		"## Root cause",
		"## Lessons",
		"## Follow-up actions",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestMarkdownRendersTimelineTable(t *testing.T) {
	ts := time.Date(2026, 10, 27, 22, 38, 15, 0, time.UTC)
	md := Markdown(Input{
		Title: "test",
		Since: ts,
		Until: ts.Add(time.Hour),
		Events: []timeline.Event{
			{Time: ts, Source: "kubectl", Summary: "pod went Running"},
		},
	})

	if !strings.Contains(md, "| 2026-10-27 22:38:15 | kubectl | pod went Running |") {
		t.Errorf("expected row not found in output:\n%s", md)
	}
}

func TestMarkdownEscapesPipesInSummary(t *testing.T) {
	ts := time.Date(2026, 10, 27, 22, 38, 15, 0, time.UTC)
	md := Markdown(Input{
		Title: "test",
		Since: ts,
		Until: ts.Add(time.Hour),
		Events: []timeline.Event{
			{Time: ts, Source: "kubectl", Summary: "foo | bar | baz"},
		},
	})

	if !strings.Contains(md, `foo \| bar \| baz`) {
		t.Errorf("pipes not escaped in output:\n%s", md)
	}
}
