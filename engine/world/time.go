package world

import (
	"time"
)

type Time struct {
	Current         time.Time
	DeltaTime       time.Duration
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
