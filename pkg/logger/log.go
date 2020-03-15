package logger

import (
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	DefaultLogFileMaxSize = 500

	DefaultLogFileMaxAge = 30

	UnavailableLogLevel = "unavailable log level"

	UnavailableLogFile = "unavailable log file path"
)

var Log *zap.Logger

func InitLog() error {

	logFilePathInput := config.ProjectConfig.LogPath
	logLevelInput := config.ProjectConfig.LogLevel
	logFileMaxSize := config.ProjectConfig.LogFileMaxSize
	logFileMaxAge := config.ProjectConfig.LogFileMaxAge

	logFilepath, err := getLogFilepath(logFilePathInput)
	if err != nil {
		return err
	}

	level, err := getLogLevel(logLevelInput)
	if err != nil {
		return err
	}

	// If the log file size and log file retention time are not configured,
	// the default configuration will be used.

	// the default maximum size of a single log file is 500M,
	// and the log file retention time is 30 days.
	hook := lumberjack.Logger{
		Filename: logFilepath,
		MaxSize:  500,
		MaxAge:   30,
		Compress: true,
	}

	fileWriter := zapcore.AddSync(&hook)

	// 对并发不安全的WriteSync加锁包装为并发安全的WriteSync
	consoleErrWriter := zapcore.Lock(os.Stderr)

	// 设置不同日志输出方式的 core Encoder
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// 创建自定义的core
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, level)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleErrWriter, level)
	core := zapcore.NewTee(fileCore, consoleCore)

	// 创建自定义的logger
	Log = zap.New(core, zap.AddCaller())
	return nil
}

func getLogFilepath(filepath string) (string, error) {

	ErrUnavailableLogFilepath := errors.New(UnavailableLogFile)
	if filepath == "" {
		return "", ErrUnavailableLogFilepath
	}

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", ErrUnavailableLogFilepath
	}
	f.Close()

	return filepath, nil
}

func getLogLevel(level string) (zapcore.Level, error) {

	logLevel := zap.InfoLevel

	switch level {
	case "debug":
		logLevel = zap.DebugLevel

	case "info":
		logLevel = zap.InfoLevel

	case "warning":
		logLevel = zap.WarnLevel

	case "error":
		logLevel = zap.ErrorLevel

	case "fatal":
		logLevel = zap.FatalLevel

	default:
		ErrUnavailableLogLevel := errors.New(UnavailableLogLevel)
		return logLevel, ErrUnavailableLogLevel
	}

	return logLevel, nil
}
