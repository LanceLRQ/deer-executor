package logger

import (
	"fmt"
	"time"
)

const (
	LogLevelDebug = iota + 1
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var LogLevelMapping = []string{
	"",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

var LogLevelStrMapping = map[string]int{
	"debug": LogLevelDebug,
	"info":  LogLevelInfo,
	"warn":  LogLevelWarn,
	"error": LogLevelError,
}

// 用于记录评测日志的工具
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

// 创建评测日志
func NewJudgeLogger() *JudgeLogger {
	logger := JudgeLogger{}
	logger.logs = make([]JudgeLogItem, 0, 5)
	return &logger
}

// 输出Log的基础函数
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

func (logger *JudgeLogger) Logf(level int, msg string, args ...interface{}) {
	logger.Log(level, fmt.Sprintf(msg, args...))
}

func (logger *JudgeLogger) Debug(msg string) {
	logger.Log(LogLevelDebug, msg)
}

func (logger *JudgeLogger) Info(msg string) {
	logger.Log(LogLevelInfo, msg)
}

func (logger *JudgeLogger) Warn(msg string) {
	logger.Log(LogLevelWarn, msg)
}

func (logger *JudgeLogger) Error(msg string) {
	logger.Log(LogLevelError, msg)
}

func (logger *JudgeLogger) Debugf(msg string, args ...interface{}) {
	logger.Logf(LogLevelDebug, msg, args...)
}

func (logger *JudgeLogger) Infof(msg string, args ...interface{}) {
	logger.Logf(LogLevelInfo, msg, args...)
}

func (logger *JudgeLogger) Warnf(msg string, args ...interface{}) {
	logger.Logf(LogLevelWarn, msg, args...)
}

func (logger *JudgeLogger) Errorf(msg string, args ...interface{}) {
	logger.Logf(LogLevelError, msg, args...)
}

func (logger *JudgeLogger) GetLogs() []JudgeLogItem {
	return logger.logs
}

func (logger *JudgeLogger) getDurationTimeStr(d time.Duration) string {
	m := d / time.Minute
	s := d / time.Second % 60
	ms := d / time.Millisecond % 1000
	return fmt.Sprintf("%02d:%02d.%03d", m, s, ms)
}

func (logger *JudgeLogger) SetStdoutPrint(swPrint bool) {
	logger.swPrint = swPrint
}

// 设置日志等级，会打印它和比它大级别的日志。比如设置为 WARN，则WARN和ERROR会被输出。
func (logger *JudgeLogger) SetLogLevel(level int) {
	logger.printLevel = level
}
