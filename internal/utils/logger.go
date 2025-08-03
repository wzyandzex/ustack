package utils

import (
	"log"
	"os"
)

// Logger 日志记录器
type Logger struct {
	*log.Logger
	level int
}

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var (
	DefaultLogger = NewLogger(INFO)
)

// NewLogger 创建新的日志记录器
func NewLogger(level int) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

// Debug 调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.Printf("[DEBUG] "+format, v...)
	}
}

// Info 信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.Printf("[INFO] "+format, v...)
	}
}

// Warn 警告日志
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.Printf("[WARN] "+format, v...)
	}
}

// Error 错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.Printf("[ERROR] "+format, v...)
	}
}

// LogPacket 记录数据包信息
func (l *Logger) LogPacket(direction, protocol string, src, dst string, length int) {
	l.Info("%s %s packet: %s -> %s (%d bytes)", direction, protocol, src, dst, length)
}

// LogConnection 记录连接信息
func (l *Logger) LogConnection(event, src, dst string) {
	l.Info("%s connection: %s -> %s", event, src, dst)
}
