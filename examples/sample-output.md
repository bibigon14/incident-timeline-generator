# Incident: vmalert CrashLoopBackOff thundering herd

**Date:** 2026-10-27 22:00 PDT → 2026-10-28 00:30 PDT
**Severity:** SEV-3
**Status:** _TODO - resolved / ongoing / monitoring_
**Author:** Dmitry Stepanov

## Summary

_TODO: one-paragraph plain-English description of what happened._

## Impact

_TODO: who was affected, for how long, what work was lost or delayed._

## Timeline

| Time | Source | Event |
|------|--------|-------|
| 2026-10-27 15:52:14 | git | commit 8c9cd2f by Dmitry Stepanov: initial scaffold: SLO + burn-rate alerts for LLM inference |
| 2026-10-27 15:53:02 | kubectl | [Normal] vm/pod/vmalert-main-68b46f8b4c-86mmh Scheduled: Successfully assigned vm/vmalert-main-68b46f8b4c-86mmh to ebpf-observer |
| 2026-10-27 15:53:48 | kubectl | [Warning] vm/pod/vmalert-main-68b46f8b4c-86mmh Unhealthy: Readiness probe failed: HTTP probe failed with statuscode: 503 |
| 2026-10-27 22:28:01 | prometheus | alert ScrapeTargetMissing fired (severity=warning) [cluster=oci-ampere-a1, instance=10.42.0.27:8080, job=vmalert-main] |
| 2026-10-27 22:28:01 | prometheus | alert ScrapeTargetMissing fired (severity=warning) [cluster=oci-ampere-a1, instance=10.42.0.27:8435, job=vmalert-main-reloader-http] |
| 2026-10-27 22:52:46 | kubectl | [Warning] vm/pod/vmalert-main-68b46f8b4c-86mmh BackOff: Back-off restarting failed container vmalert in pod |
| 2026-10-27 22:58:09 | kubectl | [Normal] vm/pod/vmsingle-main-79c8586f99-t5jb8 Created: Container created |
| 2026-10-27 22:58:10 | kubectl | [Normal] vm/pod/vmsingle-main-79c8586f99-t5jb8 Started: Container started |
| 2026-10-27 23:04:00 | kubectl | [Normal] vm/pod/vmalert-main-68b46f8b4c-rrgnr Killing: Stopping container vmalert |
| 2026-10-28 00:14:15 | git | commit b6aea33 by Dmitry Stepanov: docs: add post-mortem for 2026-10-27 vmalert CrashLoopBackOff incident |

## Root cause

_TODO: what, mechanically, caused this. Trace back from the first observed symptom to the initial trigger._

## Lessons

_TODO: 3-5 numbered lessons. Prefer actionable statements over generic observations._

## Follow-up actions

- [ ] _TODO: specific checklist items with owners_
