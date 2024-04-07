package world

import (
	"time"
)

type Time struct {
	DeltaTime      time.Duration
	Current        time.Time
	UpdateInterval time.Duration
}

func newTime() *Time {
	return &Time{Current: time.Now().UTC(), UpdateInterval: time.Second / 60}
}

func (t *Time) update() {
	now := time.Now().UTC()
	t.DeltaTime = now.Sub(t.Current)
	t.Current = now
}
