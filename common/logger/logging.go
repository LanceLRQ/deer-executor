package logger

import (
	"fmt"
	"time"
)

const (
	// LogLevelDebug log level: debug
	LogLevelDebug = iota + 1
	// LogLevelInfo log level: info
	LogLevelInfo
	// LogLevelWarn log level: warn
	LogLevelWarn
	// LogLevelError log level: error
	LogLevelError
)

// LogLevelMapping log level mapping
var LogLevelMapping = []string{
	"",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

// LogLevelStrMapping log level string mapping
var LogLevelStrMapping = map[string]int{
	"debug": LogLevelDebug,
	"info":  LogLevelInfo,
	"warn":  LogLevelWarn,
	"error": LogLevelError,
}

// JudgeLogger 评测日志
type JudgeLogger struct {
	// 日志数据
	logs []JudgeLogItem
	// T-0时间
	startTime time.Time
	// 是否打印日志
	swPrint bool
	// 设置日志打印的等级，默认是全部
	printLevel int
}

// NewJudgeLogger 创建评测日志
func NewJudgeLogger() *JudgeLogger {
	logger := JudgeLogger{}
	logger.logs = make([]JudgeLogItem, 0, 5)
	return &logger
}

// Log 输出Log的基础函数
func (logger *JudgeLogger) Log(level int, msg string) {
	nowTime := time.Now()
	if len(logger.logs) <= 0 {
		// 如果还没有写入过日志，则以这个时间作为起点
		logger.startTime = nowTime
	}
	timeDistance := nowTime.Sub(logger.startTime)
	log := JudgeLogItem{
		Message:   msg,
		Level:     level,
		Timestamp: float64(timeDistance.Nanoseconds()) / 1000000000.0,
	}
	logger.logs = append(logger.logs, log)
	if logger.swPrint && level >= logger.printLevel {
		fmt.Printf(
			"[%s] %s %s\n",
			logger.getDurationTimeStr(timeDistance),
			LogLevelMapping[level],
			msg,
		)
	}
}

// Logf 输出Log并格式化
func (logger *JudgeLogger) Logf(level int, msg string, args ...interface{}) {
	logger.Log(level, fmt.Sprintf(msg, args...))
}

// Debug 记录debug信息
func (logger *JudgeLogger) Debug(msg string) {
	logger.Log(LogLevelDebug, msg)
}

// Info 记录info信息
func (logger *JudgeLogger) Info(msg string) {
	logger.Log(LogLevelInfo, msg)
}

// Warn 记录warn信息
func (logger *JudgeLogger) Warn(msg string) {
	logger.Log(LogLevelWarn, msg)
}

// Error 记录error信息
func (logger *JudgeLogger) Error(msg string) {
	logger.Log(LogLevelError, msg)
}

// Debugf 格式化并记录debug信息
func (logger *JudgeLogger) Debugf(msg string, args ...interface{}) {
	logger.Logf(LogLevelDebug, msg, args...)
}

// Infof 格式化并记录info信息
func (logger *JudgeLogger) Infof(msg string, args ...interface{}) {
	logger.Logf(LogLevelInfo, msg, args...)
}

// Warnf 格式化并记录warn信息
func (logger *JudgeLogger) Warnf(msg string, args ...interface{}) {
	logger.Logf(LogLevelWarn, msg, args...)
}

// Errorf 格式化并记录error信息
func (logger *JudgeLogger) Errorf(msg string, args ...interface{}) {
	logger.Logf(LogLevelError, msg, args...)
}

// GetLogs 获取当前的日志列表
func (logger *JudgeLogger) GetLogs() []JudgeLogItem {
	return logger.logs
}

// 时间戳转文字格式
func (logger *JudgeLogger) getDurationTimeStr(d time.Duration) string {
	m := d / time.Minute
	s := d / time.Second % 60
	ms := d / time.Millisecond % 1000
	return fmt.Sprintf("%02d:%02d.%03d", m, s, ms)
}

// SetStdoutPrint 设置输出流
func (logger *JudgeLogger) SetStdoutPrint(swPrint bool) {
	logger.swPrint = swPrint
}

// SetLogLevel 设置日志等级，会打印它和比它大级别的日志。比如设置为 WARN，则WARN和ERROR会被输出。
func (logger *JudgeLogger) SetLogLevel(level int) {
	logger.printLevel = level
}
