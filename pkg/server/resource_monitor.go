package server

import "runtime"

type ResourceMonitor interface {
	GetUsedMemory() int
}

type ResourceMonitorImpl struct {
}

// GetUsedMemory implements ResourceMonitor.
func (r *ResourceMonitorImpl) GetUsedMemory() int {
	// Get server status
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	return int(memStat.Sys)
}

var _ ResourceMonitor = (*ResourceMonitorImpl)(nil)

type ResourceMonitorMock struct {
	UsedMemory int
}

// GetUsedMemory implements ResourceMonitor.
func (r *ResourceMonitorMock) GetUsedMemory() int {
	return r.UsedMemory
}

var _ ResourceMonitor = (*ResourceMonitorMock)(nil)
