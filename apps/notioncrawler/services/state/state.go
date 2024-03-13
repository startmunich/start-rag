package state

import "google.golang.org/genproto/googleapis/type/date"

type State struct {
	IsRunning bool
	InQueue   uint64
	Processed uint64

	LastRunDuration uint64
	LastRunEndedAt  date.Date
}
