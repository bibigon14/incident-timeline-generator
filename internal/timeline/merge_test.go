package timeline

import (
	"testing"
	"time"
)

func TestMergeSortsByTimeStably(t *testing.T) {
	t1 := time.Date(2026, 10, 27, 22, 38, 0, 0, time.UTC)
	t2 := time.Date(2026, 10, 27, 22, 40, 0, 0, time.UTC)
	t3 := time.Date(2026, 10, 27, 22, 42, 0, 0, time.UTC)

	in := []Event{
		{Time: t3, Source: "git", Summary: "commit c"},
		{Time: t1, Source: "kubectl", Summary: "event a"},
		{Time: t2, Source: "prometheus", Summary: "alert b"},
	}
	got := Merge(in)

	if len(got) != 3 {
		t.Fatalf("want 3 events, got %d", len(got))
	}
	if got[0].Summary != "event a" {
		t.Errorf("want oldest first (event a), got %q", got[0].Summary)
	}
	if got[2].Summary != "commit c" {
		t.Errorf("want newest last (commit c), got %q", got[2].Summary)
	}
}

func TestMergeStableOrderingForSameTimestamp(t *testing.T) {
	ts := time.Date(2026, 10, 27, 22, 38, 0, 0, time.UTC)
	in := []Event{
		{Time: ts, Source: "prometheus", Summary: "alert"},
		{Time: ts, Source: "git", Summary: "commit"},
		{Time: ts, Source: "kubectl", Summary: "event"},
	}
	got := Merge(in)

	// Expect alphabetical source order within the same timestamp.
	if got[0].Source != "git" || got[1].Source != "kubectl" || got[2].Source != "prometheus" {
		t.Errorf("want git/kubectl/prometheus order within same ts, got %q/%q/%q",
			got[0].Source, got[1].Source, got[2].Source)
	}
}

func TestMergeEmpty(t *testing.T) {
	got := Merge(nil)
	if len(got) != 0 {
		t.Errorf("want empty, got %v", got)
	}
}

func TestMergeDoesNotMutateInput(t *testing.T) {
	t1 := time.Date(2026, 10, 27, 22, 38, 0, 0, time.UTC)
	t2 := time.Date(2026, 10, 27, 22, 40, 0, 0, time.UTC)
	in := []Event{
		{Time: t2, Source: "b", Summary: "second"},
		{Time: t1, Source: "a", Summary: "first"},
	}
	original := make([]Event, len(in))
	copy(original, in)

	_ = Merge(in)

	for i, e := range in {
		if e != original[i] {
			t.Errorf("input mutated at index %d: want %+v, got %+v", i, original[i], e)
		}
	}
}
