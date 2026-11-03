package sources

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

// FetchGitCommits returns commits in the repo whose committer timestamp
// falls inside [since, until], across all branches reachable from HEAD.
// Deploys and configuration changes land in git before they land in
// production; having them inline with the alert timeline often makes the
// trigger obvious.
func FetchGitCommits(repoPath string, since, until time.Time) ([]timeline.Event, error) {
	// %cI = committer ISO 8601 strict, %h = short SHA, %s = subject,
	// %an = author name. Separator chosen as unlikely to appear in real
	// commit messages.
	const sep = "||"
	const format = "%cI" + sep + "%h" + sep + "%an" + sep + "%s"

	cmd := exec.Command("git",
		"-C", repoPath,
		"log",
		"--since="+since.Format(time.RFC3339),
		"--until="+until.Format(time.RFC3339),
		"--format="+format,
		"--all",
	)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git exited %d: %s",
				ee.ExitCode(), strings.TrimSpace(string(ee.Stderr)))
		}
		return nil, fmt.Errorf("run git: %w", err)
	}

	var events []timeline.Event
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, sep, 4)
		if len(parts) != 4 {
			continue
		}
		ts, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			continue
		}
		sha, author, subject := parts[1], parts[2], parts[3]

		events = append(events, timeline.Event{
			Time:    ts.Local(),
			Source:  "git",
			Summary: fmt.Sprintf("commit %s by %s: %s", sha, author, truncate(subject, 100)),
		})
	}
	return events, nil
}
