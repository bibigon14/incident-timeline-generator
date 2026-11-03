// Package render turns a merged timeline into a Markdown post-mortem
// skeleton. The template matches the shape used in
// github.com/bibigon14/inference-slo-stack/docs/incidents/ - header,
// Summary / Impact / Timeline / Root cause / Lessons / Follow-up.
// Only Timeline is populated; the rest is left blank for the human to
// fill in from the raw material.
package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

// Input is what Markdown needs to render a skeleton.
type Input struct {
	Title    string
	Severity string
	Author   string
	Since    time.Time
	Until    time.Time
	Events   []timeline.Event
}

// Markdown returns a post-mortem skeleton as a UTF-8 Markdown string.
func Markdown(in Input) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Incident: %s\n\n", in.Title)
	fmt.Fprintf(&b, "**Date:** %s → %s\n",
		in.Since.Format("2006-01-02 15:04 MST"),
		in.Until.Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "**Severity:** %s\n", in.Severity)
	fmt.Fprintf(&b, "**Status:** _TODO - resolved / ongoing / monitoring_\n")
	fmt.Fprintf(&b, "**Author:** %s\n\n", in.Author)

	b.WriteString("## Summary\n\n_TODO: one-paragraph plain-English description of what happened._\n\n")
	b.WriteString("## Impact\n\n_TODO: who was affected, for how long, what work was lost or delayed._\n\n")

	b.WriteString("## Timeline\n\n")
	b.WriteString("| Time | Source | Event |\n")
	b.WriteString("|------|--------|-------|\n")
	for _, e := range in.Events {
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
			e.Time.Format("2006-01-02 15:04:05"),
			e.Source,
			escapeTableCell(e.Summary)))
	}
	b.WriteString("\n")

	b.WriteString("## Root cause\n\n_TODO: what, mechanically, caused this. Trace back from the first observed symptom to the initial trigger._\n\n")
	b.WriteString("## Lessons\n\n_TODO: 3-5 numbered lessons. Prefer actionable statements over generic observations._\n\n")
	b.WriteString("## Follow-up actions\n\n")
	b.WriteString("- [ ] _TODO: specific checklist items with owners_\n\n")

	return b.String()
}

// escapeTableCell protects pipes in event messages from being interpreted
// as column boundaries.
func escapeTableCell(s string) string {
	return strings.ReplaceAll(s, "|", `\|`)
}
