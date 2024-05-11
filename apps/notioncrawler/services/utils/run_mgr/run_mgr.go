package run_mgr

import (
	"sync"
	"time"
)

const sleepResolution = 100 * time.Millisecond

type RunMgr struct {
	mu        sync.Mutex
	cancelRun bool
}

func New() *RunMgr {
	return &RunMgr{}
}

func (r *RunMgr) ShouldCancel() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cancelRun
}

func (r *RunMgr) CancelCurrentRun() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelRun = true
}

func (r *RunMgr) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelRun = false
}

func (r *RunMgr) SleepCancelAware(duration time.Duration) bool {
	repeats := int(duration / sleepResolution)
	for i := 0; i < repeats; i++ {
		time.Sleep(sleepResolution)
		if r.ShouldCancel() {
			return true
		}
	}
	return false
}
