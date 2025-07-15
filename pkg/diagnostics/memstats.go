package diagnostics

import (
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/logger"
)

func LogMemoryStats(stage string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logger.Info("Memory stats snapshot", log.Fields{
		"stage":        stage,
		"AllocKB":      m.Alloc / 1024,
		"TotalAllocKB": m.TotalAlloc / 1024,
		"SysKB":        m.Sys / 1024,
		"NumGC":        m.NumGC,
		"PauseTotalMS": m.PauseTotalNs / 1e6, // convert to milliseconds
		"HeapAllocKB":  m.HeapAlloc / 1024,
		"HeapIdleKB":   m.HeapIdle / 1024,
		"HeapInuseKB":  m.HeapInuse / 1024,
		"HeapReleased": m.HeapReleased / 1024,
		"StackInuseKB": m.StackInuse / 1024,
	})
}
