// Package watch provides port monitoring primitives for portwatch.
//
// # Forecast
//
// Forecast tracks per-port presence history over a configurable rolling window
// and exposes a likelihood score (0.0–1.0) indicating how often a port has
// been observed open in recent scans.
//
// Use Predict to gate actions on ports that exceed a confidence threshold:
//
//	f, _ := watch.NewForecast(10)
//	f.Observe(443, true)
//	if f.Predict(443, 0.8) {
//	    // port 443 is very likely to be open
//	}
//
// Forecast is safe for concurrent use.
package watch
