package memory_manager

import (
	"runtime"
	"time"
)

func NewMemoryManager() *MemoryManager {
	managerByProcess := &MemoryManager{
		MemorySize:      NewMemorySize(1024 * 1024),
		memoryUsage:     0,
		initMemoryUsage: 0,
	}

	// 更新内存使用量定时任务
	go func() {
		for range time.Tick(time.Second * 10) {
			managerByProcess.readMemStats()
		}
	}()

	return managerByProcess
}

// MemoryManager 内存管理器
// 异步更新进程内存使用量
// 由于计算的是整个进程的内存，不能单独计算缓存数据的内存，所以只适合独立部署的缓存服务
type MemoryManager struct {
	MemorySize
	memoryUsage     uint64 // 当前内存使用量
	initMemoryUsage uint64 // 初始内存使用量
}

// SetMaxMemory 设置最大内存
func (m *MemoryManager) SetMaxMemory(size string) bool {
	return m.Parse(size)
}

// IsFull 是否达到内存上限
func (m *MemoryManager) IsFull() bool {
	return (m.memoryUsage - m.initMemoryUsage) > m.GetLimitSize()
}

// readMemStats 读取内存使用量
func (m *MemoryManager) readMemStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	if m.initMemoryUsage == 0 {
		m.initMemoryUsage = memStats.Alloc
	}
	m.memoryUsage = memStats.Alloc
}
