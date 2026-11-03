package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bibigon14/incident-timeline-generator/internal/timeline"
)

// FetchPrometheusAlerts queries the Prometheus (or VictoriaMetrics) ALERTS
// synthetic series for the window and returns one event per alert
// state=firing transition into the window.
//
// We use ALERTS over /api/v1/alerts because /api/v1/alerts only shows
// currently-firing state; ALERTS as a metric gives us historical state
// transitions, which is what we need to reconstruct an incident timeline.
func FetchPrometheusAlerts(baseURL string, since, until time.Time) ([]timeline.Event, error) {
	// Query ALERTS metric with step=30s over the window. Each sample where
	// alertstate="firing" represents the alert being active at that
	// evaluation. We detect transitions (gaps > step) to emit one event
	// per firing episode.
	q := `ALERTS{alertstate="firing"}`
	step := 30 * time.Second

	u, err := url.Parse(strings.TrimRight(baseURL, "/") + "/api/v1/query_range")
	if err != nil {
		return nil, fmt.Errorf("parse prometheus url: %w", err)
	}
	qs := u.Query()
	qs.Set("query", q)
	qs.Set("start", strconv.FormatInt(since.Unix(), 10))
	qs.Set("end", strconv.FormatInt(until.Unix(), 10))
	qs.Set("step", step.String())
	u.RawQuery = qs.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus %d: %s", resp.StatusCode, truncate(string(body), 200))
	}

	var pr promRangeResponse
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if pr.Status != "success" {
		return nil, fmt.Errorf("prometheus status=%s", pr.Status)
	}

	var events []timeline.Event
	for _, series := range pr.Data.Result {
		alertName := series.Metric["alertname"]
		if alertName == "" {
			alertName = "unknown"
		}
		severity := series.Metric["severity"]

		// Detect firing episodes: consecutive samples <= step*2 apart
		// are the same episode; a bigger gap starts a new one.
		var episodeStart time.Time
		var lastSample time.Time
		for _, pair := range series.Values {
			if len(pair) < 2 {
				continue
			}
			tsFloat, ok := pair[0].(float64)
			if !ok {
				continue
			}
			ts := time.Unix(int64(tsFloat), 0).Local()

			if lastSample.IsZero() || ts.Sub(lastSample) > step*2 {
				// Emit previous episode if any.
				if !episodeStart.IsZero() {
					events = append(events, firingEvent(alertName, severity, series.Metric, episodeStart))
				}
				episodeStart = ts
			}
			lastSample = ts
		}
		if !episodeStart.IsZero() {
			events = append(events, firingEvent(alertName, severity, series.Metric, episodeStart))
		}
	}
	return events, nil
}

func firingEvent(alertName, severity string, labels map[string]string, firedAt time.Time) timeline.Event {
	// Build a compact labelset for context, dropping the obvious ones.
	var ctx []string
	for k, v := range labels {
		switch k {
		case "alertname", "alertstate", "__name__", "severity":
			continue
		}
		ctx = append(ctx, fmt.Sprintf("%s=%s", k, v))
	}

	summary := fmt.Sprintf("alert %s fired", alertName)
	if severity != "" {
		summary = fmt.Sprintf("alert %s fired (severity=%s)", alertName, severity)
	}
	if len(ctx) > 0 {
		summary += " [" + strings.Join(ctx, ", ") + "]"
	}

	return timeline.Event{
		Time:    firedAt,
		Source:  "prometheus",
		Summary: summary,
	}
}

type promRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}
