# incident-timeline-generator

Generates a Markdown post-mortem skeleton from the raw signals you already
have during an incident: Kubernetes events, Prometheus alerts, and git
commits. Takes a time window, pulls from each source, merges by timestamp,
renders into the standard post-mortem shape you'd otherwise type by hand.

Built because I noticed I was reconstructing incident timelines manually
across two production post-mortems last week - the same mechanical
assembly of "kubectl events sorted by time" + "prometheus alerts fired"
+ "git commits deployed around the window". This tool does the first
draft so the human can focus on root cause and lessons.

## What it does

```
   kubectl events ─┐
                   │
   prometheus      ├─► merge by timestamp ─► markdown skeleton
   alerts          │                           - Summary (empty)
                   │                           - Timeline table (populated)
   git commits ────┘                           - Root cause (empty)
                                               - Lessons (empty)
                                               - Follow-up (empty)
```

The automation stops at the mechanical step. Root cause analysis and
lessons stay human - the tool just removes the typing.

## Install

```bash
go install github.com/bibigon14/incident-timeline-generator/cmd/itg@latest
```

Or build locally:

```bash
git clone https://github.com/bibigon14/incident-timeline-generator
cd incident-timeline-generator
make build
```

## Usage

```bash
itg \
  --since "2026-10-27 22:00" \
  --until "2026-10-28 00:30" \
  --kubectl-namespace vm \
  --prometheus-url http://localhost:8429 \
  --git-repo ~/git/inference-slo-stack \
  --title "vmalert CrashLoopBackOff" \
  --severity SEV-3 \
  --output docs/incidents/2026-10-27-vmalert-crashloop.md
```

What this does:

1. Runs `kubectl -n vm get events --sort-by=.lastTimestamp` and filters
   to events within the window
2. Queries `prometheus-url/api/v1/alerts` and
   `prometheus-url/api/v1/query_range` for ALERTS in the window
3. Runs `git -C <repo> log --since=... --until=... --format=...` and
   collects commits (including the one that likely triggered)
4. Sorts all three by timestamp into a unified table
5. Renders the markdown with a header pre-filled, the timeline populated,
   and headings waiting for your words in Root cause / Lessons / Follow-up

## Sources

Each source is an optional ingester. The tool runs whatever you point at
and skips the rest.

| Source | Flag | What it ingests |
|--------|------|------------------|
| Kubernetes events | `--kubectl-namespace` (or `--all-namespaces`) | `kubectl get events` for the window |
| Prometheus alerts | `--prometheus-url` | ALERTS metric for `firing` state transitions |
| git commits | `--git-repo` | commits in the window on the default branch |

Future ingesters planned for v0.2:
- Loki log excerpts around critical events
- Grafana dashboard snapshots (via API)
- `journalctl` for host-level events on single-node clusters

## Output template

The rendered markdown follows the shape of the two post-mortems in
[inference-slo-stack/docs/incidents/](https://github.com/bibigon14/inference-slo-stack/tree/main/docs/incidents),
because that's the shape I was tired of typing:

```
# Incident: <title>

**Date:** <from --since> → <to --until>
**Severity:** <--severity>
**Status:** <blank, for human>
**Author:** <--author or git config user.name>

## Summary

<blank, for human>

## Impact

<blank, for human>

## Timeline

| Time | Source | Event |
|------|--------|-------|
| ... (populated from merged ingesters) |

## Root cause

<blank, for human>

## Lessons

<blank, for human>

## Follow-up actions

<blank, for human>
```

See [examples/sample-output.md](examples/sample-output.md) for a real
rendered skeleton from the 2026-10-27 vmalert incident.

## Why not <existing tool>

- **blameless, Rootly, FireHydrant, PagerDuty** - all great SaaS, all
  presume you're already paying for them. This is for homelab and small
  SRE teams who want the output artifact without the SaaS.
- **k8s-event-logger / prom2log** - these transform events, don't
  structure them into a post-mortem shape.
- **Writing it by hand** - works, but the mechanical ingestion step is
  typing, not thinking.

## Status

v0.1 - working MVP. All three ingesters functional, markdown output
usable as-is. Known gaps:
- No log excerpt ingester yet (coming v0.2)
- Prometheus auth not supported (bearer token only in v0.2)
- Hardcoded to GET events; no watch mode

## License

MIT
