package rclient

// Metrics
type DataInMetricCount struct {
	Count int `json:"count"`
}

type dataInMetricGauges struct {
	Value float64 `json:"value"`
}

type dataInMetricMeters struct {
	Count    int     `json:"count"`
	M15Rate  float64 `json:"m15_rate"`
	M1Rate   float64 `json:"m1_rate"`
	M5Rate   float64 `json:"m5_rate"`
	MeanRate float64 `json:"mean_rate"`
}

type dataInMetricTimers struct {
	Count    int     `json:"count"`
	Max      float64 `json:"max"`
	Mean     float64 `json:"mean"`
	P50      float64 `json:"p50"`
	P75      float64 `json:"p75"`
	P95      float64 `json:"p95"`
	P98      float64 `json:"p98"`
	P99      float64 `json:"p99"`
	P999     float64 `json:"p999"`
	Stddev   float64 `json:"stddev"`
	M15Rate  float64 `json:"m15_rate"`
	M1Rate   float64 `json:"m1_rate"`
	M5Rate   float64 `json:"m5_rate"`
	MeanRate float64 `json:"mean_rate"`
}
