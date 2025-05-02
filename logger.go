package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// 日志级别
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// 日志记录器
type Logger struct {
	mu     sync.Mutex
	level  LogLevel
	prefix string
	logger *log.Logger
}

// 单例日志实例
var (
	defaultLogger *Logger
	once          sync.Once
)

// 获取默认日志记录器
func GetLogger() *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			level:  InfoLevel,
			prefix: "[ImageMaster] ",
			logger: log.New(os.Stderr, "", log.LstdFlags),
		}
	})
	return defaultLogger
}

// 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// 设置日志前缀
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// 格式化并输出日志
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 添加日志级别前缀
	var levelStr string
	switch level {
	case DebugLevel:
		levelStr = "[DEBUG] "
	case InfoLevel:
		levelStr = "[INFO] "
	case WarnLevel:
		levelStr = "[WARN] "
	case ErrorLevel:
		levelStr = "[ERROR] "
	case FatalLevel:
		levelStr = "[FATAL] "
	}

	// 格式化消息
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fullMsg := fmt.Sprintf("%s %s%s%s", timeStr, l.prefix, levelStr, msg)

	// 输出日志
	l.logger.Println(fullMsg)

	// 如果是致命错误，退出程序
	if level == FatalLevel {
		os.Exit(1)
	}
}

// 不同级别的日志方法
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
}

// 简便的全局方法
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
} 