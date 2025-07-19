package monitoring

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// MemoryMonitor monitors memory usage and detects potential leaks
type MemoryMonitor struct {
	logger   *logrus.Logger
	interval time.Duration
	
	// Memory statistics
	mutex     sync.RWMutex
	snapshots []MemorySnapshot
	maxSnapshots int
	
	// Thresholds for alerts
	memoryThreshold   uint64 // Memory usage threshold in bytes
	goroutineThreshold int   // Goroutine count threshold
	
	// Control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// MemorySnapshot represents a point-in-time memory snapshot
type MemorySnapshot struct {
	Timestamp    time.Time `json:"timestamp"`
	HeapAlloc    uint64    `json:"heap_alloc"`     // Heap allocated bytes
	HeapSys      uint64    `json:"heap_sys"`       // Heap system bytes
	HeapInuse    uint64    `json:"heap_inuse"`     // Heap in-use bytes
	HeapReleased uint64    `json:"heap_released"`  // Heap released bytes
	StackInuse   uint64    `json:"stack_inuse"`    // Stack in-use bytes
	NumGC        uint32    `json:"num_gc"`         // Number of GC cycles
	NumGoroutine int       `json:"num_goroutine"`  // Number of goroutines
	
	// Calculated metrics
	HeapUtilization float64 `json:"heap_utilization"` // Heap utilization percentage
	GCPressure      float64 `json:"gc_pressure"`      // GC pressure indicator
}

// MemoryStats provides aggregated memory statistics
type MemoryStats struct {
	Current     MemorySnapshot `json:"current"`
	Peak        MemorySnapshot `json:"peak"`
	Average     MemorySnapshot `json:"average"`
	TrendSlope  float64        `json:"trend_slope"`  // Memory usage trend
	LeakSuspect bool           `json:"leak_suspect"` // Potential memory leak detected
}

// NewMemoryMonitor creates a new memory monitor
func NewMemoryMonitor(logger *logrus.Logger, interval time.Duration) *MemoryMonitor {
	if logger == nil {
		logger = logrus.New()
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &MemoryMonitor{
		logger:             logger,
		interval:           interval,
		maxSnapshots:       100, // Keep last 100 snapshots
		memoryThreshold:    100 * 1024 * 1024, // 100MB default threshold
		goroutineThreshold: 1000, // 1000 goroutines threshold
		ctx:                ctx,
		cancel:             cancel,
	}
}

// Start begins memory monitoring
func (m *MemoryMonitor) Start() {
	m.wg.Add(1)
	go m.monitorLoop()
}

// Stop stops memory monitoring
func (m *MemoryMonitor) Stop() {
	m.cancel()
	m.wg.Wait()
}

// SetThresholds sets memory and goroutine thresholds for alerts
func (m *MemoryMonitor) SetThresholds(memoryMB int, goroutines int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.memoryThreshold = uint64(memoryMB * 1024 * 1024)
	m.goroutineThreshold = goroutines
}

// GetCurrentSnapshot returns the current memory snapshot
func (m *MemoryMonitor) GetCurrentSnapshot() MemorySnapshot {
	return m.takeSnapshot()
}

// GetStats returns aggregated memory statistics
func (m *MemoryMonitor) GetStats() MemoryStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	if len(m.snapshots) == 0 {
		current := m.takeSnapshot()
		return MemoryStats{
			Current: current,
			Peak:    current,
			Average: current,
		}
	}
	
	current := m.snapshots[len(m.snapshots)-1]
	peak := m.findPeak()
	average := m.calculateAverage()
	trend := m.calculateTrend()
	leak := m.detectLeak()
	
	return MemoryStats{
		Current:     current,
		Peak:        peak,
		Average:     average,
		TrendSlope:  trend,
		LeakSuspect: leak,
	}
}

// GetSnapshots returns all memory snapshots
func (m *MemoryMonitor) GetSnapshots() []MemorySnapshot {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// Return a copy to avoid race conditions
	snapshots := make([]MemorySnapshot, len(m.snapshots))
	copy(snapshots, m.snapshots)
	return snapshots
}

// ForceGC forces garbage collection and returns before/after snapshots
func (m *MemoryMonitor) ForceGC() (before, after MemorySnapshot) {
	before = m.takeSnapshot()
	runtime.GC()
	runtime.GC() // Run twice to ensure cleanup
	after = m.takeSnapshot()
	
	m.logger.WithFields(logrus.Fields{
		"heap_before": before.HeapAlloc,
		"heap_after":  after.HeapAlloc,
		"freed":       before.HeapAlloc - after.HeapAlloc,
	}).Info("Forced garbage collection")
	
	return before, after
}

// monitorLoop is the main monitoring loop
func (m *MemoryMonitor) monitorLoop() {
	defer m.wg.Done()
	
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			snapshot := m.takeSnapshot()
			m.addSnapshot(snapshot)
			m.checkThresholds(snapshot)
			
		case <-m.ctx.Done():
			return
		}
	}
}

// takeSnapshot captures current memory statistics
func (m *MemoryMonitor) takeSnapshot() MemorySnapshot {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	numGoroutine := runtime.NumGoroutine()
	
	// Calculate derived metrics
	var heapUtilization float64
	if memStats.HeapSys > 0 {
		heapUtilization = float64(memStats.HeapInuse) / float64(memStats.HeapSys) * 100
	}
	
	// Simple GC pressure calculation (GCs per minute)
	var gcPressure float64
	if len(m.snapshots) > 0 {
		lastSnapshot := m.snapshots[len(m.snapshots)-1]
		timeDiff := time.Since(lastSnapshot.Timestamp).Minutes()
		if timeDiff > 0 {
			gcDiff := float64(memStats.NumGC - lastSnapshot.NumGC)
			gcPressure = gcDiff / timeDiff
		}
	}
	
	return MemorySnapshot{
		Timestamp:       time.Now(),
		HeapAlloc:       memStats.HeapAlloc,
		HeapSys:         memStats.HeapSys,
		HeapInuse:       memStats.HeapInuse,
		HeapReleased:    memStats.HeapReleased,
		StackInuse:      memStats.StackInuse,
		NumGC:           memStats.NumGC,
		NumGoroutine:    numGoroutine,
		HeapUtilization: heapUtilization,
		GCPressure:      gcPressure,
	}
}

// addSnapshot adds a snapshot to the history
func (m *MemoryMonitor) addSnapshot(snapshot MemorySnapshot) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.snapshots = append(m.snapshots, snapshot)
	
	// Keep only the last N snapshots
	if len(m.snapshots) > m.maxSnapshots {
		m.snapshots = m.snapshots[1:]
	}
}

// checkThresholds checks if memory usage exceeds thresholds
func (m *MemoryMonitor) checkThresholds(snapshot MemorySnapshot) {
	if snapshot.HeapAlloc > m.memoryThreshold {
		m.logger.WithFields(logrus.Fields{
			"heap_alloc":  snapshot.HeapAlloc,
			"threshold":   m.memoryThreshold,
			"goroutines":  snapshot.NumGoroutine,
		}).Warn("Memory usage exceeds threshold")
	}
	
	if snapshot.NumGoroutine > m.goroutineThreshold {
		m.logger.WithFields(logrus.Fields{
			"goroutines": snapshot.NumGoroutine,
			"threshold":  m.goroutineThreshold,
		}).Warn("Goroutine count exceeds threshold")
	}
}

// findPeak finds the snapshot with highest memory usage
func (m *MemoryMonitor) findPeak() MemorySnapshot {
	if len(m.snapshots) == 0 {
		return MemorySnapshot{}
	}
	
	peak := m.snapshots[0]
	for _, snapshot := range m.snapshots {
		if snapshot.HeapAlloc > peak.HeapAlloc {
			peak = snapshot
		}
	}
	return peak
}

// calculateAverage calculates average memory usage
func (m *MemoryMonitor) calculateAverage() MemorySnapshot {
	if len(m.snapshots) == 0 {
		return MemorySnapshot{}
	}
	
	var sum MemorySnapshot
	for _, snapshot := range m.snapshots {
		sum.HeapAlloc += snapshot.HeapAlloc
		sum.HeapSys += snapshot.HeapSys
		sum.HeapInuse += snapshot.HeapInuse
		sum.StackInuse += snapshot.StackInuse
		sum.NumGoroutine += snapshot.NumGoroutine
	}
	
	count := uint64(len(m.snapshots))
	return MemorySnapshot{
		Timestamp:    time.Now(),
		HeapAlloc:    sum.HeapAlloc / count,
		HeapSys:      sum.HeapSys / count,
		HeapInuse:    sum.HeapInuse / count,
		StackInuse:   sum.StackInuse / count,
		NumGoroutine: sum.NumGoroutine / len(m.snapshots),
	}
}

// calculateTrend calculates memory usage trend (slope)
func (m *MemoryMonitor) calculateTrend() float64 {
	if len(m.snapshots) < 2 {
		return 0
	}
	
	// Simple linear regression for trend calculation
	n := float64(len(m.snapshots))
	var sumX, sumY, sumXY, sumX2 float64
	
	for i, snapshot := range m.snapshots {
		x := float64(i)
		y := float64(snapshot.HeapAlloc)
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	// Calculate slope (trend)
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0
	}
	
	return (n*sumXY - sumX*sumY) / denominator
}

// detectLeak detects potential memory leaks
func (m *MemoryMonitor) detectLeak() bool {
	if len(m.snapshots) < 10 {
		return false // Need enough data points
	}
	
	// Check if memory usage is consistently increasing
	trend := m.calculateTrend()
	
	// Check if recent memory usage is significantly higher than average
	recent := m.snapshots[len(m.snapshots)-5:] // Last 5 snapshots
	var recentAvg uint64
	for _, snapshot := range recent {
		recentAvg += snapshot.HeapAlloc
	}
	recentAvg /= uint64(len(recent))
	
	overall := m.calculateAverage()
	
	// Leak indicators:
	// 1. Positive trend (memory increasing over time)
	// 2. Recent usage significantly higher than overall average
	// 3. High goroutine count
	trendThreshold := 1000.0 // bytes per snapshot
	recentThreshold := 1.5   // 50% higher than average
	
	return trend > trendThreshold && 
		   float64(recentAvg) > float64(overall.HeapAlloc)*recentThreshold
}

// PrintStats prints memory statistics to logger
func (m *MemoryMonitor) PrintStats() {
	stats := m.GetStats()
	
	m.logger.WithFields(logrus.Fields{
		"heap_alloc_mb":    stats.Current.HeapAlloc / 1024 / 1024,
		"heap_sys_mb":      stats.Current.HeapSys / 1024 / 1024,
		"goroutines":       stats.Current.NumGoroutine,
		"gc_cycles":        stats.Current.NumGC,
		"heap_utilization": stats.Current.HeapUtilization,
		"trend_slope":      stats.TrendSlope,
		"leak_suspect":     stats.LeakSuspect,
	}).Info("Memory statistics")
}
