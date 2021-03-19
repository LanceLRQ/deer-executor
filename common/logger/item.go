package logger

// JudgeLogItem judge log item structs
type JudgeLogItem struct {
	// 时间戳
	Timestamp float64 `json:"timestamp" bson:"timestamp"`
	//
	Duration float64 `json:"duration" bson:"duration"`
	// 日志等级
	Level int `json:"level" bson:"level"`
	// 日志消息
	Message string `json:"msg" bson:"msg"`
}
