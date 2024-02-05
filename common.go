package utils

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	logger "github.com/skyfox2000/nect-utils/logger"
)

var Logger *logger.LoggerEntry

func ShowSysStat(prefix, reqId string, maxMem int, logger *logger.LoggerEntry) {
	// runtime.GC()
	numGoroutines := runtime.NumGoroutine()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	stackInUse := memStats.StackInuse / 1024 // 将字节数转换为KB

	logger.Warn(prefix, " ReqId: ", reqId)

	logger.Warn("内存使用情况")
	logger.Warn("  当前的goroutine数量: ", numGoroutines)
	logger.Warn("  内存使用量（字节）:")
	logger.Warn("    总分配内存: ", formatBytes(memStats.TotalAlloc))
	logger.Warn("    堆上对象数: ", memStats.HeapObjects)
	logger.Warn("    堆内存使用量/可用量: ", formatBytes(memStats.HeapAlloc), "/", formatBytes(uint64(maxMem)))
	logger.Warn("    堆内存回收量: ", formatBytes(memStats.HeapReleased))
	logger.Warn("  栈的总数量: ", stackInUse)

	if memStats.HeapAlloc > uint64(maxMem) {
		logger.Fatal("系统运行异常，内存溢出")
		os.Exit(2)
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func Monitor(logger *logger.LoggerEntry) {
	// 定时获取 Goroutine 的信息
	// GoRouteMap := make(map[string]*GoroutineInfo)

	for {
		// 计算第一次触发的等待时间
		nextRunTime := time.Now().Add(time.Second * 10)
		// 定时等待直到next时间
		<-time.After(time.Until(nextRunTime))

		ShowSysStat("Begin:", "", int(10240000000), logger)
		// 获取 Goroutine 的信息
		pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	}
}
