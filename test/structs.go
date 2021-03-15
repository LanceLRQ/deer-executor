package test

type JudgementBenchmark struct {
    TimeUsed         float64        `json:"time_used"`
    CurrectCounter   map[string]int `json:"currect_counter"`
    IncurrectCounter map[string]int `json:"incurrect_counter"`
    Message          string         `json:"message"`
}
