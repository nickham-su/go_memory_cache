package memory_manager

import (
	"regexp"
	"strconv"
	"strings"
)

// NewMemorySize 解析内存大小
// minSize 最小内存使用量，防止用户设置的太小，或者没设置
func NewMemorySize(minSize uint64) MemorySize {
	return MemorySize{
		minSize: minSize,
	}
}

type MemorySize struct {
	minSize   uint64 // 最小内存使用量，防止用户设置的太小，或者没设置
	limitSize uint64 // 内存限制
}

// Parse 解析内存大小，支持KB、MB、GB
func (m *MemorySize) Parse(size string) bool {
	re, _ := regexp.Compile(`^(\d+)([KMG]B)$`)
	result := re.FindStringSubmatch(strings.ToUpper(size))
	if len(result) != 3 {
		return false
	}
	num, err := strconv.ParseInt(result[1], 10, 64)
	if err != nil {
		return false
	}
	if result[2] == "KB" {
		num *= 1024
	} else if result[2] == "MB" {
		num *= 1024 * 1024
	} else if result[2] == "GB" {
		num *= 1024 * 1024 * 1024
	}
	m.limitSize = uint64(num)
	return true
}

// GetLimitSize 获取内存限制大小
func (m *MemorySize) GetLimitSize() uint64 {
	if m.limitSize < m.minSize {
		return m.minSize
	}
	return m.limitSize
}
