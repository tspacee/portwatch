// Package watch provides building blocks for port monitoring and analysis.
//
// # HeatMap
//
// HeatMap tracks how frequently each port has been observed across scan
// cycles. Each call to Record increments the internal counter for that port.
//
// Typical usage:
//
//	hm := watch.NewHeatMap()
//	for _, port := range openPorts {
//		_ = hm.Record(port)
//	}
//	fmt.Println("hottest port:", hm.Hottest())
//
// HeatMap is safe for concurrent use.
package watch
