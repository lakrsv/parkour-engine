package world

import (
	"time"
)

type Time struct {
	DeltaTime       time.Duration
	Current         time.Time
	Timestep        time.Duration
	PhysicsTimestep time.Duration
}

func newTime() *Time {
	//return &Time{Current: time.Now().UTC(), Timestep: time.Second / 60}
	return &Time{Current: time.Now().UTC(), Timestep: 0, PhysicsTimestep: time.Second / 60}
}

func (t *Time) update() {
	now := time.Now().UTC()
	t.DeltaTime = now.Sub(t.Current)
	t.Current = now
}
