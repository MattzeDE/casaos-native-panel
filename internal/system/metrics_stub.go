package system

// non-Linux stub for GetMetrics so the project can build on Windows (or other OSes).
// The real implementation lives in metrics.go with a //go:build linux tag.

func GetMetrics() (Metrics, error) {
	return Metrics{}, nil
}
