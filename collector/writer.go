package collector

import (
	"fmt"
	"io"
	"encoding/json"
	"log"
)

const collectdIntGaugeTemplate = "PUTVAL %s/docker_stats-%s.%s/gauge-%s %d:%d\n"

// CollectdWriter is responsible for writing data
// to wrapped writer in collectd exec plugin format
type CollectdWriter struct {
	host     string
	writer   io.Writer
	prevCPUTotal map[string]uint64
	prevSystemTotal map[string]uint64
}

// NewCollectdWriter creates new CollectdWriter
// with specified hostname and writer
func NewCollectdWriter(host string, writer io.Writer) CollectdWriter {

	return CollectdWriter{
		host:   host,
		writer: writer,
		prevCPUTotal: make(map[string]uint64),
		prevSystemTotal: make(map[string]uint64),
	}
}

func (w *CollectdWriter) Write(s Stats) error {
	return w.writeInts(s)
}

func (w *CollectdWriter) writeInts(s Stats) error {
	t := uint64(s.Stats.Read.Unix())
	var ct float32
	if cpuTotal, ok := w.prevCPUTotal[s.App]; ok {
		ct = (float32)(s.Stats.CPUStats.CPUUsage.TotalUsage - cpuTotal) /
		      (float32)(s.Stats.CPUStats.SystemCPUUsage - w.prevSystemTotal[s.App])
		fmt.Printf("CPU %.4f\n", (ct * 100))
	} else {
		ct = 1000000 * 100
	}
	w.prevCPUTotal[s.App] = s.Stats.CPUStats.CPUUsage.TotalUsage
	w.prevSystemTotal[s.App] = s.Stats.CPUStats.SystemCPUUsage

	metrics := map[string]uint64{
		"ts": t,

		"cpu.perc": uint64(ct * 1000000),
		// "cpu.user":   s.Stats.CPUStats.CPUUsage.UsageInUsermode,
		// "cpu.system": s.Stats.CPUStats.CPUUsage.UsageInKernelmode,
		"cpu.total":  s.Stats.CPUStats.CPUUsage.TotalUsage,
		"cpu.system": s.Stats.CPUStats.SystemCPUUsage,

		"memory.limit": s.Stats.MemoryStats.Limit,
		"memory.max":   s.Stats.MemoryStats.MaxUsage,
		"memory.usage": s.Stats.MemoryStats.Usage,

		// "memory.active_anon":   s.Stats.MemoryStats.Stats.TotalActiveAnon,
		// "memory.active_file":   s.Stats.MemoryStats.Stats.TotalActiveFile,
		// "memory.cache":         s.Stats.MemoryStats.Stats.TotalCache,
		// "memory.inactive_anon": s.Stats.MemoryStats.Stats.TotalInactiveAnon,
		// "memory.inactive_file": s.Stats.MemoryStats.Stats.TotalInactiveFile,
		// "memory.mapped_file":   s.Stats.MemoryStats.Stats.TotalMappedFile,
		// "memory.pg_fault":      s.Stats.MemoryStats.Stats.TotalPgfault,
		// "memory.pg_in":         s.Stats.MemoryStats.Stats.TotalPgpgin,
		// "memory.pg_out":        s.Stats.MemoryStats.Stats.TotalPgpgout,
		// "memory.rss":           s.Stats.MemoryStats.Stats.TotalRss,
		// "memory.rss_huge":      s.Stats.MemoryStats.Stats.TotalRssHuge,
		// "memory.unevictable":   s.Stats.MemoryStats.Stats.TotalUnevictable,
		// "memory.writeback":     s.Stats.MemoryStats.Stats.TotalWriteback,
	}

	data := map[string]interface{}{
		"app": s.App,
		"metrics": metrics,
	}

	// for _, network := range s.Stats.Networks {
	// 	metrics["net.rx_bytes"] += network.RxBytes
	// 	metrics["net.rx_dropped"] += network.RxDropped
	// 	metrics["net.rx_errors"] += network.RxErrors
	// 	metrics["net.rx_packets"] += network.RxPackets
	//
	// 	metrics["net.tx_bytes"] += network.TxBytes
	// 	metrics["net.tx_dropped"] += network.TxDropped
	// 	metrics["net.tx_errors"] += network.TxErrors
	// 	metrics["net.tx_packets"] += network.TxPackets
	// }

	b, _ := json.Marshal(data)
	log.Println(string(b))
	// for k, v := range metrics {
	// 	err := w.writeInt(s, k, t, v)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// func (w CollectdWriter) writeInt(s Stats, k string, t int64, v uint64) error {
// 	msg := fmt.Sprintf(collectdIntGaugeTemplate, w.host, s.App, s.Task, k, t, v)
// 	_, err := w.writer.Write([]byte(msg))
// 	return err
// }
