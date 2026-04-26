package envfile

import "sort"

// DiffStatus represents the state of a key when comparing two env files.
type DiffStatus string

const (
	// StatusAdded indicates the key exists only in the right/new env.
	StatusAdded DiffStatus = "added"
	// StatusRemoved indicates the key exists only in the left/old env.
	StatusRemoved DiffStatus = "removed"
	// StatusChanged indicates the key exists in both envs but with different values.
	StatusChanged DiffStatus = "changed"
	// StatusUnchanged indicates the key exists in both envs with the same value.
	StatusUnchanged DiffStatus = "unchanged"
)

// DiffEntry represents a single key comparison result between two env maps.
type DiffEntry struct {
	Key      string
	Status   DiffStatus
	OldValue string
	NewValue string
}

// Diff compares two parsed env maps and returns a sorted slice of DiffEntry.
// The left map is treated as the "old" environment and right as the "new" one.
func Diff(left, right map[string]string) []DiffEntry {
	results := make([]DiffEntry, 0)

	// Track all keys seen across both maps.
	allKeys := make(map[string]struct{})
	for k := range left {
		allKeys[k] = struct{}{}
	}
	for k := range right {
		allKeys[k] = struct{}{}
	}

	for key := range allKeys {
		leftVal, inLeft := left[key]
		rightVal, inRight := right[key]

		var entry DiffEntry
		entry.Key = key

		switch {
		case inLeft && inRight:
			if leftVal == rightVal {
				entry.Status = StatusUnchanged
			} else {
				entry.Status = StatusChanged
			}
			entry.OldValue = leftVal
			entry.NewValue = rightVal
		case inLeft && !inRight:
			entry.Status = StatusRemoved
			entry.OldValue = leftVal
		case !inLeft && inRight:
			entry.Status = StatusAdded
			entry.NewValue = rightVal
		}

		results = append(results, entry)
	}

	// Sort by key name for deterministic output.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

// FilterByStatus returns only the diff entries matching one of the provided statuses.
func FilterByStatus(entries []DiffEntry, statuses ...DiffStatus) []DiffEntry {
	statusSet := make(map[DiffStatus]struct{}, len(statuses))
	for _, s := range statuses {
		statusSet[s] = struct{}{}
	}

	filtered := make([]DiffEntry, 0, len(entries))
	for _, e := range entries {
		if _, ok := statusSet[e.Status]; ok {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// HasChanges reports whether any entry in the diff represents a meaningful change
// (i.e., added, removed, or changed — not unchanged).
func HasChanges(entries []DiffEntry) bool {
	for _, e := range entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}
