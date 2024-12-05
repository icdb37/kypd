// Package logx 日志
package logx

import (
	"fmt"
	"strings"
)

// Level 日志等级
type Level = uint8

const (
	// LevelDebug 调试
	LevelDebug Level = 0
	// LevelInfo 信息
	LevelInfo Level = 1
	// LevelWarn 警告
	LevelWarn Level = 2
	// LevelError 错误
	LevelError Level = 3
)

// 日志等级
var (
	level Level = LevelWarn
)

// SetLevel 设置日志等级
func SetLevel(l string) {
	l = strings.ToLower(l)
	switch l {
	case "debug":
		level = LevelDebug
	case "info":
		level = LevelInfo
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	}
}

// Logger 日志等级
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Debug 调试日志
func Debug(args ...interface{}) {
	if level > LevelDebug {
		return
	}
	fmt.Println(args...)
}

// Debugf 调试日志
func Debugf(format string, args ...interface{}) {
	if level > LevelDebug {
		return
	}
	fmt.Printf(format+"\n", args...)
}

// Infof 信息日志
func Infof(format string, args ...interface{}) {
	if level > LevelInfo {
		return
	}
	fmt.Printf(format+"\n", args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	if level > LevelInfo {
		return
	}
	fmt.Println(args...)
}

// Warnf 警告日志
func Warnf(format string, args ...interface{}) {
	if level > LevelWarn {
		return
	}
	fmt.Printf(format+"\n", args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	if level > LevelWarn {
		return
	}
	fmt.Println(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	if level < LevelError {
		return
	}
	fmt.Println(args...)
}

// Errorf 错误日志
func Errorf(format string, args ...interface{}) {
	if level > LevelError {
		return
	}
	fmt.Printf(format+"\n", args...)
}
