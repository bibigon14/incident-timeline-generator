// Command itg generates a Markdown post-mortem skeleton by merging
// Kubernetes events, Prometheus alerts, and git commits for a given
// time window. The automation stops at the mechanical assembly step -
// root cause analysis and lessons remain human.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/render"
	"github.com/bibigon14/incident-timeline-generator/internal/sources"
	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

const timeLayout = "2006-01-02 15:04"

type config struct {
	since     string
	until     string
	title     string
	severity  string
	author    string
	namespace string
	allNS     bool
	promURL   string
	gitRepo   string
	output    string
}

func main() {
	cfg := parseFlags()

	sinceT, err := time.ParseInLocation(timeLayout, cfg.since, time.Local)
	if err != nil {
		die("parse --since: %v (expected format: %q)", err, timeLayout)
	}
	untilT, err := time.ParseInLocation(timeLayout, cfg.until, time.Local)
	if err != nil {
		die("parse --until: %v (expected format: %q)", err, timeLayout)
	}
	if !untilT.After(sinceT) {
		die("--until (%s) must be after --since (%s)", cfg.until, cfg.since)
	}

	author := cfg.author
	if author == "" {
		author = gitUserName()
	}

	var events []timeline.Event

	// Ingest Kubernetes events if asked.
	if cfg.namespace != "" || cfg.allNS {
		ns := cfg.namespace
		if cfg.allNS {
			ns = ""
		}
		k8sEvents, err := sources.FetchKubectlEvents(ns, sinceT, untilT)
		if err != nil {
			warn("kubectl events: %v (continuing without them)", err)
		} else {
			events = append(events, k8sEvents...)
		}
	}

	// Ingest Prometheus alerts if a URL is given.
	if cfg.promURL != "" {
		promAlerts, err := sources.FetchPrometheusAlerts(cfg.promURL, sinceT, untilT)
		if err != nil {
			warn("prometheus alerts: %v (continuing without them)", err)
		} else {
			events = append(events, promAlerts...)
		}
	}

	// Ingest git commits if a repo path is given.
	if cfg.gitRepo != "" {
		commits, err := sources.FetchGitCommits(cfg.gitRepo, sinceT, untilT)
		if err != nil {
			warn("git commits: %v (continuing without them)", err)
		} else {
			events = append(events, commits...)
		}
	}

	if len(events) == 0 {
		die("no events found from any source in [%s, %s]", cfg.since, cfg.until)
	}

	merged := timeline.Merge(events)

	md := render.Markdown(render.Input{
		Title:    cfg.title,
		Severity: cfg.severity,
		Author:   author,
		Since:    sinceT,
		Until:    untilT,
		Events:   merged,
	})

	if cfg.output == "" || cfg.output == "-" {
		fmt.Print(md)
		return
	}

	if err := os.WriteFile(cfg.output, []byte(md), 0o644); err != nil {
		die("write %s: %v", cfg.output, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d events across %d sources)\n",
		cfg.output, len(merged), countSources(merged))
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.since, "since", "", "window start, format: "+timeLayout+" (local time)")
	flag.StringVar(&cfg.until, "until", "", "window end, format: "+timeLayout+" (local time)")
	flag.StringVar(&cfg.title, "title", "Untitled incident", "incident title for the post-mortem header")
	flag.StringVar(&cfg.severity, "severity", "SEV-?", "severity label (SEV-1..SEV-4)")
	flag.StringVar(&cfg.author, "author", "", "author name (defaults to git config user.name)")
	flag.StringVar(&cfg.namespace, "kubectl-namespace", "", "kubernetes namespace to pull events from")
	flag.BoolVar(&cfg.allNS, "all-namespaces", false, "pull kubectl events across all namespaces")
	flag.StringVar(&cfg.promURL, "prometheus-url", "", "prometheus base URL (http://host:port)")
	flag.StringVar(&cfg.gitRepo, "git-repo", "", "path to a git repo to pull commits from")
	flag.StringVar(&cfg.output, "output", "-", "output file path, or '-' for stdout")
	flag.Parse()

	if cfg.since == "" || cfg.until == "" {
		fmt.Fprintln(os.Stderr, "both --since and --until are required")
		flag.Usage()
		os.Exit(2)
	}
	if cfg.namespace == "" && !cfg.allNS && cfg.promURL == "" && cfg.gitRepo == "" {
		fmt.Fprintln(os.Stderr, "give at least one source: --kubectl-namespace / --all-namespaces / --prometheus-url / --git-repo")
		os.Exit(2)
	}
	return cfg
}

func gitUserName() string {
	out, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func countSources(events []timeline.Event) int {
	seen := map[string]struct{}{}
	for _, e := range events {
		seen[e.Source] = struct{}{}
	}
	return len(seen)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

func warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "warn: "+format+"\n", args...)
}
