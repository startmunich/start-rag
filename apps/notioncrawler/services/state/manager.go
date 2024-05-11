package state

import (
	"sync"
)

type Manager struct {
	state *State
	lock  sync.Mutex
}

func New() *Manager {
	return &Manager{
		state: &State{
			IsRunning:        false,
			InQueue:          0,
			Processed:        0,
			LastRunDuration:  0,
			LastRunStartedAt: 0,
			LastRunEndedAt:   0,
		},
		lock: sync.Mutex{},
	}
}

func (m *Manager) UpdateIsRunning(isRunning bool) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.IsRunning = isRunning
	return m
}

func (m *Manager) UpdateInQueue(inQueue uint64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.InQueue = inQueue
	return m
}

func (m *Manager) UpdateCacheMisses(cacheMisses uint64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.CacheMisses = cacheMisses
	return m
}

func (m *Manager) UpdateProcessed(processed uint64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.Processed = processed
	return m
}

func (m *Manager) UpdateLastRunDuration(lastRunDuration uint64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.LastRunDuration = lastRunDuration
	return m
}

func (m *Manager) UpdateLastRunEndedAt(lastRunEndedAt int64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.LastRunEndedAt = lastRunEndedAt
	return m
}

func (m *Manager) UpdateLastRunStartedAt(lastRunStartedAt int64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.LastRunStartedAt = lastRunStartedAt
	return m
}

func (m *Manager) UpdateNextRunAt(nextRunAt int64) *Manager {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.NextRunAt = nextRunAt
	return m
}

func (m *Manager) GetState() State {
	m.lock.Lock()
	defer m.lock.Unlock()
	return *m.state
}
