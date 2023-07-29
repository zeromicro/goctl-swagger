package stat

import "time"

// A Task is a task that is reported to Metrics.
type Task struct {
	Drop        bool
	Duration    time.Duration
	Description string
}
