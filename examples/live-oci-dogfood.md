# Live run - OCI Ampere A1 (2026-11-05)

This file is a verbatim dump of what `itg v0.2.0` produces on a real
cluster, not a hand-crafted template. Everything below the horizontal
rule was emitted by the tool itself with zero hand-editing. I ran it
on the OCI Ampere A1 VM that hosts `inference-slo-stack`,
`kafka-slo-observability`, `ebpf-tcp-observer` and the rest of the
`sre.dstepanov.dev`-adjacent portfolio stack.

**Command:**

```bash
./itg \
  --since "2026-11-05 20:00" \
  --until "2026-11-06 02:00" \
  --all-namespaces \
  --prometheus-url http://10.43.169.103:9090 \
  --git-repo ~/inference-slo-stack \
  --journalctl --journalctl-priority 4 \
  --title "OCI Ampere A1 - 6h window dogfood run of itg v0.2.0" \
  --severity SEV-? \
  --output itg-live-oci-run.md
# wrote itg-live-oci-run.md (145 events across 3 sources)
```

**Breakdown:**

| Source | Events | Note |
|--------|--------|------|
| kubectl | 138 | `k get events --all-namespaces` over 6h - mostly k6 loadgen CronJob churn plus a real `kafka/demo-brokers` CrashLoop that itg surfaced alongside the noise |
| prometheus | 5 | `ALERTS{alertstate="firing"}` over the window - k3s control-plane-absent cohort (KubeAPIDown, KubeControllerManagerDown, KubeProxyDown, KubeSchedulerDown) plus Watchdog |
| git | 0 | quiet repo for the window |
| journal | 2 | `-o json --priority=4` - one `ssh.service` error, one other `err`-level entry |

**Note on `Author: unknown`** - deliberate signal, not a bug. `itg`
defaults `--author` to `git config user.name`; the OCI VM didn't have
one set, and the tool fell back to a placeholder rather than failing.
Shows the graceful-degradation stance the code takes everywhere: miss
a source? warn and continue. Miss a config? fall back and continue.

---

# Incident: OCI Ampere A1 - 6h window dogfood run of itg v0.2.0

**Date:** 2026-11-05 20:00 UTC → 2026-11-06 02:00 UTC
**Severity:** SEV-?
**Status:** _TODO - resolved / ongoing / monitoring_
**Author:** unknown

## Summary

_TODO: one-paragraph plain-English description of what happened._

## Impact

_TODO: who was affected, for how long, what work was lost or delayed._

## Timeline

| Time | Source | Event |
|------|--------|-------|
| 2026-11-05 20:00:00 | prometheus | alert KubeAPIDown fired (severity=critical) [job=apiserver] |
| 2026-11-05 20:00:00 | prometheus | alert KubeControllerManagerDown fired (severity=critical) [job=kube-controller-manager] |
| 2026-11-05 20:00:00 | prometheus | alert KubeProxyDown fired (severity=critical) [job=kube-proxy] |
| 2026-11-05 20:00:00 | prometheus | alert KubeSchedulerDown fired (severity=critical) [job=kube-scheduler] |
| 2026-11-05 20:00:00 | prometheus | alert Watchdog fired (severity=none) |
| 2026-11-05 23:10:47 | journal | [err] ssh.service: error: kex_exchange_identification: read: Connection reset by peer |
| 2026-11-06 01:00:00 | kubectl | [Normal] inference/job/k6-loadgen-29898780 SuccessfulCreate: Created pod: k6-loadgen-29898780-fvtdk |
| 2026-11-06 01:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898780-fvtdk Created: Container created |
| 2026-11-06 01:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898780-fvtdk Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898780-fvtdk Started: Container started |
| 2026-11-06 01:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898780-fvtdk Scheduled: Successfully assigned inference/k6-loadgen-29898780-fvtdk to ebpf-observer |
| 2026-11-06 01:04:03 | kubectl | [Warning] inference/job/k6-loadgen-29898780 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:05:00 | kubectl | [Normal] inference/job/k6-loadgen-29898785 SuccessfulCreate: Created pod: k6-loadgen-29898785-56lvr |
| 2026-11-06 01:05:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898785-56lvr Created: Container created |
| 2026-11-06 01:05:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898785-56lvr Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:05:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898785-56lvr Scheduled: Successfully assigned inference/k6-loadgen-29898785-56lvr to ebpf-observer |
| 2026-11-06 01:05:01 | kubectl | [Normal] inference/pod/k6-loadgen-29898785-56lvr Started: Container started |
| 2026-11-06 01:09:05 | kubectl | [Warning] inference/job/k6-loadgen-29898785 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:09:09 | kubectl | [Normal] monitoring/pod/monitoring-kube-state-metrics-764584797b-q72sg Created: Container created |
| 2026-11-06 01:10:00 | kubectl | [Normal] inference/job/k6-loadgen-29898790 SuccessfulCreate: Created pod: k6-loadgen-29898790-26f7p |
| 2026-11-06 01:10:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898790-26f7p Created: Container created |
| 2026-11-06 01:10:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898790-26f7p Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:10:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898790-26f7p Started: Container started |
| 2026-11-06 01:10:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898790-26f7p Scheduled: Successfully assigned inference/k6-loadgen-29898790-26f7p to ebpf-observer |
| 2026-11-06 01:14:04 | kubectl | [Warning] inference/job/k6-loadgen-29898790 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:15:00 | kubectl | [Normal] inference/job/k6-loadgen-29898795 SuccessfulCreate: Created pod: k6-loadgen-29898795-2sgx4 |
| 2026-11-06 01:15:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898795-2sgx4 Created: Container created |
| 2026-11-06 01:15:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898795-2sgx4 Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:15:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898795-2sgx4 Started: Container started |
| 2026-11-06 01:15:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898795-2sgx4 Scheduled: Successfully assigned inference/k6-loadgen-29898795-2sgx4 to ebpf-observer |
| 2026-11-06 01:18:53 | kubectl | [Normal] kafka/pod/kafka-exporter-5bc797fd55-db99x Created: Container created |
| 2026-11-06 01:19:04 | kubectl | [Warning] inference/job/k6-loadgen-29898795 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:19:27 | kubectl | [Normal] monitoring/pod/monitoring-grafana-5cd4c45755-4mzpk Started: Container started |
| 2026-11-06 01:20:00 | kubectl | [Normal] inference/job/k6-loadgen-29898800 SuccessfulCreate: Created pod: k6-loadgen-29898800-4cp8v |
| 2026-11-06 01:20:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898800-4cp8v Created: Container created |
| 2026-11-06 01:20:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898800-4cp8v Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:20:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898800-4cp8v Started: Container started |
| 2026-11-06 01:20:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898800-4cp8v Scheduled: Successfully assigned inference/k6-loadgen-29898800-4cp8v to ebpf-observer |
| 2026-11-06 01:23:08 | kubectl | [Normal] vm/pod/vmoperator-victoria-metrics-operator-7c67bd477c-d28k2 Created: Container created |
| 2026-11-06 01:23:29 | kubectl | [Warning] vm/pod/vmoperator-victoria-metrics-operator-7c67bd477c-d28k2 Unhealthy: Readiness probe failed: Get "http://10.42.0.64:8081/health": context deadline exceeded (Client.Timeout exceeded while aw... |
| 2026-11-06 01:23:59 | kubectl | [Normal] monitoring/pod/monitoring-kube-prometheus-operator-849fc5d448-72tbf Killing: Container kube-prometheus-stack failed liveness probe, will be restarted |
| 2026-11-06 01:24:04 | kubectl | [Warning] inference/job/k6-loadgen-29898800 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:25:00 | kubectl | [Normal] inference/job/k6-loadgen-29898805 SuccessfulCreate: Created pod: k6-loadgen-29898805-rqzmx |
| 2026-11-06 01:25:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898805-rqzmx Created: Container created |
| 2026-11-06 01:25:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898805-rqzmx Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:25:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898805-rqzmx Started: Container started |
| 2026-11-06 01:25:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898805-rqzmx Scheduled: Successfully assigned inference/k6-loadgen-29898805-rqzmx to ebpf-observer |
| 2026-11-06 01:29:05 | kubectl | [Warning] inference/job/k6-loadgen-29898805 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:29:20 | kubectl | [Normal] kube-system/pod/metrics-server-6dc596dfb8-k88jx Created: Container created |
| 2026-11-06 01:29:26 | kubectl | [Warning] kube-system/pod/metrics-server-6dc596dfb8-k88jx Unhealthy: Readiness probe failed: Get "https://10.42.0.7:10250/readyz": net/http: request canceled while waiting for connection (C... |
| 2026-11-06 01:30:00 | kubectl | [Normal] inference/job/k6-loadgen-29898810 SuccessfulCreate: Created pod: k6-loadgen-29898810-vvpgt |
| 2026-11-06 01:30:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898810-vvpgt Created: Container created |
| 2026-11-06 01:30:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898810-vvpgt Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:30:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898810-vvpgt Started: Container started |
| 2026-11-06 01:30:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898810-vvpgt Scheduled: Successfully assigned inference/k6-loadgen-29898810-vvpgt to ebpf-observer |
| 2026-11-06 01:31:48 | kubectl | [Normal] kafka/pod/demo-entity-operator-689fdbbb75-tnxzv Pulled: Container image "quay.io/strimzi/operator:0.47.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:33:47 | kubectl | [Normal] monitoring/pod/monitoring-kube-state-metrics-764584797b-q72sg Killing: Container kube-state-metrics failed liveness probe, will be restarted |
| 2026-11-06 01:34:05 | kubectl | [Warning] inference/job/k6-loadgen-29898810 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:35:00 | kubectl | [Normal] inference/job/k6-loadgen-29898815 SuccessfulCreate: Created pod: k6-loadgen-29898815-lhfmz |
| 2026-11-06 01:35:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898815-lhfmz Created: Container created |
| 2026-11-06 01:35:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898815-lhfmz Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:35:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898815-lhfmz Started: Container started |
| 2026-11-06 01:35:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898815-lhfmz Scheduled: Successfully assigned inference/k6-loadgen-29898815-lhfmz to ebpf-observer |
| 2026-11-06 01:36:49 | kubectl | [Normal] vm/pod/vmagent-main-5b68496ccf-t7xkp Created: Container created |
| 2026-11-06 01:39:03 | kubectl | [Warning] inference/job/k6-loadgen-29898815 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:39:24 | kubectl | [Warning] monitoring/pod/monitoring-kube-state-metrics-764584797b-q72sg Unhealthy: Readiness probe failed: Get "http://10.42.0.186:8081/readyz": dial tcp 10.42.0.186:8081: connect: connection refused |
| 2026-11-06 01:40:00 | kubectl | [Normal] inference/job/k6-loadgen-29898820 SuccessfulCreate: Created pod: k6-loadgen-29898820-gpgpf |
| 2026-11-06 01:40:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898820-gpgpf Created: Container created |
| 2026-11-06 01:40:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898820-gpgpf Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:40:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898820-gpgpf Started: Container started |
| 2026-11-06 01:40:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898820-gpgpf Scheduled: Successfully assigned inference/k6-loadgen-29898820-gpgpf to ebpf-observer |
| 2026-11-06 01:40:24 | kubectl | [Warning] kube-system/pod/metrics-server-6dc596dfb8-k88jx Unhealthy: Readiness probe failed: Get "https://10.42.0.7:10250/readyz": dial tcp 10.42.0.7:10250: connect: connection refused |
| 2026-11-06 01:44:04 | kubectl | [Warning] inference/job/k6-loadgen-29898820 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:44:57 | kubectl | [Warning] monitoring/pod/monitoring-kube-state-metrics-764584797b-q72sg Unhealthy: Liveness probe failed: Get "http://10.42.0.186:8080/livez": dial tcp 10.42.0.186:8080: connect: connection refused |
| 2026-11-06 01:45:00 | kubectl | [Normal] inference/job/k6-loadgen-29898825 SuccessfulCreate: Created pod: k6-loadgen-29898825-mwlh5 |
| 2026-11-06 01:45:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898825-mwlh5 Created: Container created |
| 2026-11-06 01:45:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898825-mwlh5 Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:45:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898825-mwlh5 Started: Container started |
| 2026-11-06 01:45:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898825-mwlh5 Scheduled: Successfully assigned inference/k6-loadgen-29898825-mwlh5 to ebpf-observer |
| 2026-11-06 01:47:05 | kubectl | [Normal] vm/pod/vmalertmanager-main-0 Created: Container created |
| 2026-11-06 01:47:34 | kubectl | [Normal] monitoring/pod/monitoring-grafana-5cd4c45755-4mzpk Created: Container created |
| 2026-11-06 01:47:34 | kubectl | [Normal] monitoring/pod/monitoring-grafana-5cd4c45755-4mzpk Pulled: Container image "quay.io/kiwigrid/k8s-sidecar:2.14.1" already present on machine and can be accessed by the pod |
| 2026-11-06 01:48:18 | journal | [err] ssh.service: error: kex_exchange_identification: read: Connection reset by peer |
| 2026-11-06 01:49:03 | kubectl | [Warning] inference/job/k6-loadgen-29898825 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:49:47 | kubectl | [Warning] monitoring/pod/monitoring-kube-prometheus-operator-849fc5d448-72tbf BackOff: Back-off restarting failed container kube-prometheus-stack in pod monitoring-kube-prometheus-operator-849fc5d448-72tbf_m... |
| 2026-11-06 01:49:48 | kubectl | [Warning] kafka/pod/demo-brokers-0 Unhealthy: Readiness probe failed:   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current                 ... |
| 2026-11-06 01:50:00 | kubectl | [Normal] inference/job/k6-loadgen-29898830 SuccessfulCreate: Created pod: k6-loadgen-29898830-mct5d |
| 2026-11-06 01:50:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898830-mct5d Created: Container created |
| 2026-11-06 01:50:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898830-mct5d Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:50:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898830-mct5d Started: Container started |
| 2026-11-06 01:50:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898830-mct5d Scheduled: Successfully assigned inference/k6-loadgen-29898830-mct5d to ebpf-observer |
| 2026-11-06 01:50:30 | kubectl | [Normal] monitoring/pod/monitoring-grafana-5cd4c45755-4mzpk Pulled: Container image "quay.io/kiwigrid/k8s-sidecar:2.14.1" already present on machine and can be accessed by the pod |
| 2026-11-06 01:53:46 | kubectl | [Normal] kafka/pod/demo-entity-operator-689fdbbb75-tnxzv Killing: Container user-operator failed liveness probe, will be restarted |
| 2026-11-06 01:54:03 | kubectl | [Warning] inference/job/k6-loadgen-29898830 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:54:04 | kubectl | [Normal] monitoring/pod/monitoring-kube-prometheus-operator-849fc5d448-72tbf Created: Container created |
| 2026-11-06 01:54:39 | kubectl | [Warning] monitoring/pod/monitoring-kube-prometheus-operator-849fc5d448-72tbf Unhealthy: Liveness probe failed: Get "https://10.42.0.90:10250/healthz": dial tcp 10.42.0.90:10250: connect: connection refused |
| 2026-11-06 01:54:47 | kubectl | [Warning] kube-system/pod/hubble-relay-674cf45d78-wtngj BackOff: Back-off restarting failed container hubble-relay in pod hubble-relay-674cf45d78-wtngj_kube-system(6b0b7682-0d8e-4d63-97... |
| 2026-11-06 01:54:55 | kubectl | [Normal] kafka/pod/demo-brokers-0 Created: Container created |
| 2026-11-06 01:54:55 | kubectl | [Normal] kafka/pod/demo-brokers-0 Pulled: Container image "quay.io/strimzi/kafka:0.47.0-kafka-3.9.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:54:57 | kubectl | [Warning] kafka/pod/demo-entity-operator-689fdbbb75-tnxzv BackOff: Back-off restarting failed container topic-operator in pod demo-entity-operator-689fdbbb75-tnxzv_kafka(de6df798-cfa9-439... |
| 2026-11-06 01:55:00 | kubectl | [Normal] inference/job/k6-loadgen-29898835 SuccessfulCreate: Created pod: k6-loadgen-29898835-zdkg7 |
| 2026-11-06 01:55:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898835-zdkg7 Created: Container created |
| 2026-11-06 01:55:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898835-zdkg7 Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:55:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898835-zdkg7 Started: Container started |
| 2026-11-06 01:55:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898835-zdkg7 Scheduled: Successfully assigned inference/k6-loadgen-29898835-zdkg7 to ebpf-observer |
| 2026-11-06 01:55:16 | kubectl | [Warning] kafka/pod/demo-brokers-0 Unhealthy: Readiness probe failed:   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current                 ... |
| 2026-11-06 01:55:21 | kubectl | [Warning] kafka/pod/demo-brokers-2 BackOff: Back-off restarting failed container kafka in pod demo-brokers-2_kafka(3468a7a1-c229-456a-bc28-3bf963ba2d10) |
| 2026-11-06 01:55:24 | kubectl | [Warning] kube-system/pod/local-path-provisioner-77b9867795-pxwck BackOff: Back-off restarting failed container local-path-provisioner in pod local-path-provisioner-77b9867795-pxwck_kube-system(3... |
| 2026-11-06 01:55:39 | kubectl | [Warning] kafka/pod/demo-brokers-1 BackOff: Back-off restarting failed container kafka in pod demo-brokers-1_kafka(d5de7b67-b98e-4268-91d5-9edbbff1a166) |
| 2026-11-06 01:55:48 | kubectl | [Warning] kafka/pod/demo-producer-7d47bfcf97-bjr8p BackOff: Back-off restarting failed container producer in pod demo-producer-7d47bfcf97-bjr8p_kafka(59c4999d-a23f-4341-9f6d-221a96... |
| 2026-11-06 01:56:29 | kubectl | [Normal] monitoring/pod/monitoring-kube-state-metrics-764584797b-q72sg Pulled: Container image "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.20.0" already present on machine and can be ac... |
| 2026-11-06 01:57:07 | kubectl | [Normal] vm/pod/vmoperator-victoria-metrics-operator-7c67bd477c-d28k2 Pulled: Container image "victoriametrics/operator:v0.75.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:57:26 | kubectl | [Normal] kube-system/pod/metrics-server-6dc596dfb8-k88jx Pulled: Container image "rancher/mirrored-metrics-server:v0.9.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:57:40 | kubectl | [Normal] kube-system/pod/hubble-relay-674cf45d78-wtngj Pulled: Container image "quay.io/cilium/hubble-relay:v1.17.7@sha256:9394312ce65c3c253a8c26a6c292f58736e75c78d1446ecfcd244f1418be... |
| 2026-11-06 01:57:48 | kubectl | [Normal] observability/pod/log-emitter-7d6d7c5b74-v4brk Pulled: Container image "ghcr.io/open-telemetry/opentelemetry-collector-contrib/telemetrygen:v0.113.0" already present on machin... |
| 2026-11-06 01:57:52 | kubectl | [Warning] vm/pod/vmalertmanager-main-0 BackOff: Back-off restarting failed container config-init in pod vmalertmanager-main-0_vm(9a3ad3d1-78d0-4a92-bb71-c3df824bff9e) |
| 2026-11-06 01:58:13 | kubectl | [Normal] vm/pod/vmalertmanager-main-0 Pulled: Container image "victoriametrics/operator:config-reloader-v0.75.0" already present on machine and can be accessed by the... |
| 2026-11-06 01:58:20 | kubectl | [Warning] vm/pod/vmagent-main-5b68496ccf-t7xkp BackOff: Back-off restarting failed container config-init in pod vmagent-main-5b68496ccf-t7xkp_vm(4e357fcc-3c51-4040-9231-6738046... |
| 2026-11-06 01:58:32 | kubectl | [Normal] kafka/pod/demo-producer-7d47bfcf97-bjr8p Pulled: Container image "quay.io/strimzi/kafka:0.47.0-kafka-3.9.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:58:43 | kubectl | [Normal] kafka/pod/demo-brokers-2 Created: Container created |
| 2026-11-06 01:58:43 | kubectl | [Normal] kafka/pod/demo-brokers-2 Pulled: Container image "quay.io/strimzi/kafka:0.47.0-kafka-3.9.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:58:48 | kubectl | [Normal] kube-system/pod/hubble-relay-674cf45d78-wtngj Killing: Container hubble-relay failed startup probe, will be restarted |
| 2026-11-06 01:58:58 | kubectl | [Normal] kafka/pod/demo-entity-operator-689fdbbb75-tnxzv Pulled: Container image "quay.io/strimzi/operator:0.47.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:59:04 | kubectl | [Warning] inference/job/k6-loadgen-29898835 BackoffLimitExceeded: Job has reached the specified backoff limit |
| 2026-11-06 01:59:06 | kubectl | [Normal] kafka/pod/demo-consumer-b86478f55-zc8bs Created: Container created |
| 2026-11-06 01:59:06 | kubectl | [Normal] kafka/pod/demo-consumer-b86478f55-zc8bs Pulled: Container image "quay.io/strimzi/kafka:0.47.0-kafka-3.9.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:59:11 | kubectl | [Warning] vm/pod/vmoperator-victoria-metrics-operator-7c67bd477c-d28k2 BackOff: Back-off restarting failed container operator in pod vmoperator-victoria-metrics-operator-7c67bd477c-d28k2_vm(37ef7c62-0... |
| 2026-11-06 01:59:12 | kubectl | [Normal] kube-system/pod/local-path-provisioner-77b9867795-pxwck Created: Container created |
| 2026-11-06 01:59:12 | kubectl | [Normal] kube-system/pod/local-path-provisioner-77b9867795-pxwck Pulled: Container image "rancher/local-path-provisioner:v0.0.37" already present on machine and can be accessed by the pod |
| 2026-11-06 01:59:16 | kubectl | [Warning] kafka/pod/demo-entity-operator-689fdbbb75-tnxzv Unhealthy: Liveness probe failed: Get "http://10.42.0.198:8081/healthy": dial tcp 10.42.0.198:8081: connect: connection refused |
| 2026-11-06 01:59:30 | kubectl | [Warning] kube-system/pod/hubble-relay-674cf45d78-wtngj Unhealthy: Startup probe failed: service unhealthy (responded with "NOT_SERVING") |
| 2026-11-06 01:59:35 | kubectl | [Warning] kafka/pod/demo-brokers-2 Unhealthy: Readiness probe failed:   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current                 ... |
| 2026-11-06 01:59:36 | kubectl | [Normal] kafka/pod/demo-brokers-1 Created: Container created |
| 2026-11-06 01:59:36 | kubectl | [Normal] kafka/pod/demo-brokers-1 Pulled: Container image "quay.io/strimzi/kafka:0.47.0-kafka-3.9.0" already present on machine and can be accessed by the pod |
| 2026-11-06 01:59:38 | kubectl | [Normal] vm/pod/vmagent-main-5b68496ccf-t7xkp Pulled: Container image "victoriametrics/operator:config-reloader-v0.75.0" already present on machine and can be accessed by the... |
| 2026-11-06 01:59:50 | kubectl | [Warning] kafka/pod/kafka-exporter-5bc797fd55-db99x BackOff: Back-off restarting failed container kafka-exporter in pod kafka-exporter-5bc797fd55-db99x_kafka(0a7e5b29-e199-48c8-afca... |
| 2026-11-06 01:59:52 | kubectl | [Warning] monitoring/pod/monitoring-kube-state-metrics-764584797b-q72sg BackOff: Back-off restarting failed container kube-state-metrics in pod monitoring-kube-state-metrics-764584797b-q72sg_monitoring... |
| 2026-11-06 01:59:55 | kubectl | [Warning] kube-system/pod/coredns-7cfb7bc9c7-w8bch Unhealthy: Readiness probe failed: HTTP probe failed with statuscode: 503 |
| 2026-11-06 01:59:57 | kubectl | [Warning] kafka/pod/demo-brokers-1 Unhealthy: Readiness probe failed:   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current                 ... |
| 2026-11-06 01:59:59 | kubectl | [Warning] kafka/pod/demo-brokers-0 BackOff: Back-off restarting failed container kafka in pod demo-brokers-0_kafka(393a4b34-c49b-457e-b9b7-1b8ad0b65f1c) |
| 2026-11-06 02:00:00 | kubectl | [Normal] inference/cronjob/k6-loadgen SuccessfulCreate: (combined from similar events): Created job k6-loadgen-29898840 |
| 2026-11-06 02:00:00 | kubectl | [Normal] inference/job/k6-loadgen-29898840 SuccessfulCreate: Created pod: k6-loadgen-29898840-lc4lv |
| 2026-11-06 02:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898840-lc4lv Created: Container created |
| 2026-11-06 02:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898840-lc4lv Pulled: Container image "grafana/k6:0.54.0" already present on machine and can be accessed by the pod |
| 2026-11-06 02:00:00 | kubectl | [Normal] inference/pod/k6-loadgen-29898840-lc4lv Started: Container started |

## Root cause

_TODO: what, mechanically, caused this. Trace back from the first observed symptom to the initial trigger._

## Lessons

_TODO: 3-5 numbered lessons. Prefer actionable statements over generic observations._

## Follow-up actions

- [ ] _TODO: specific checklist items with owners_

