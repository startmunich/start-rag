package state

type State struct {
	IsRunning   bool   `json:"isRunning"`
	InQueue     uint64 `json:"inQueue"`
	Processed   uint64 `json:"processed"`
	CacheMisses uint64 `json:"cacheMisses"`

	LastRunDuration  uint64 `json:"lastRunDuration"`
	LastRunStartedAt int64  `json:"lastRunStartedAt"`
	LastRunEndedAt   int64  `json:"lastRunEndedAt"`
	NextRunAt        int64  `json:"nextRunAt"`
}
