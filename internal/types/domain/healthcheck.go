package types

type Healthcheck struct {
	Test        []string `json:"test"`
	Interval    int32    `json:"interval"`
	Timeout     int32    `json:"timeout"`
	Retries     int32    `json:"retries"`
	StartPeriod int32    `json:"start_period"`
	Disable     bool     `json:"disable"`
}
