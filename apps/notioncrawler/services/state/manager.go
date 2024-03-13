package state

import "google.golang.org/genproto/googleapis/type/date"

type Manager struct {
	state *State
}

func New() *Manager {
	return &Manager{
		state: &State{
			IsRunning:       false,
			InQueue:         0,
			Processed:       0,
			LastRunDuration: 0,
			LastRunEndedAt:  date.Date{},
		},
	}
}
