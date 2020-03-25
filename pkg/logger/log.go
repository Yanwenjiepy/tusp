package logger

import (
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	UnavailableLogLevel = "unavailable log level"

	UnavailableLogFile = "unavailable log file path"

	UnavailableLogFileMaxSize = "unavailable log file max size"

	UnavailableLogFileMaxAge = "unavailable log file max age"

	UnavailableLogFileMaxBackups = "unavailable log file max backups"

	UnavailableCompressFlag = "unavailable compress flag"

	UnavailableLocalTimeFlag = "unavailable local time flag"
)

var Log *zap.Logger

func InitLog() error {

	logFilePathInput := config.ProjectConfig.LogPath
	logLevelInput := config.ProjectConfig.LogLevel
	logFileMaxSizeInput := config.ProjectConfig.LogFileMaxSize
	logFileMaxAgeInput := config.ProjectConfig.LogFileMaxAge
	logFileMaxBackupsInput := config.ProjectConfig.LogFileMaxBackups
	isUseLocalTimeInput := config.ProjectConfig.LocalTime
	isUseCompressInput := config.ProjectConfig.Compress

	logFilepath, err := getLogFilepath(logFilePathInput)
	if err != nil {
		return err
	}

	level, err := getLogLevel(logLevelInput)
	if err != nil {
		return err
	}

	maxSize, err := getLogFileMaxSize(logFileMaxSizeInput)
	if err != nil {
		return err
	}

	maxAge, err := getLogFileMaxAge(logFileMaxAgeInput)
	if err != nil {
		return err
	}

	maxBackups, err := getLogFileMaxBackups(logFileMaxBackupsInput)
	if err != nil {
		return err
	}

	localTime, err := isUseLocalTime(isUseLocalTimeInput)
	if err != nil {
		return err
	}

	compress, err := isCompressLogFile(isUseCompressInput)
	if err != nil {
		return err
	}

	fileLogger := lumberjack.Logger{
		Filename:   logFilepath,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  localTime,
		Compress:   compress,
	}

	fileWriter := zapcore.AddSync(&fileLogger)
	consoleErrWriter := zapcore.Lock(os.Stderr)

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	fileCore := zapcore.NewCore(fileEncoder, fileWriter, level)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleErrWriter, level)
	core := zapcore.NewTee(fileCore, consoleCore)

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

func getLogFileMaxSize(size int) (int, error) {

	if size < 0 {
		ErrUnavailableLogFileMaxSize := errors.New(UnavailableLogFileMaxSize)
		return 0, ErrUnavailableLogFileMaxSize
	}

	return size, nil
}

func getLogFileMaxAge(age int) (int, error) {

	if age < 0 {
		ErrUnavailableLogFileMaxAge := errors.New(UnavailableLogFileMaxAge)
		return 0, ErrUnavailableLogFileMaxAge
	}

	return age, nil
}

func getLogFileMaxBackups(backups int) (int, error) {

	if backups < 0 {
		ErrUnavailableLogFileMaxBackups := errors.New(UnavailableLogFileMaxBackups)
		return 0, ErrUnavailableLogFileMaxBackups
	}

	return backups, nil
}

func isUseLocalTime(isLocalTime int) (bool, error) {
	switch isLocalTime {

	case 0:
		return false, nil

	case 1:
		return true, nil

	default:
		ErrUnavailableLocalTimeFlag := errors.New(UnavailableLocalTimeFlag)
		return false, ErrUnavailableLocalTimeFlag
	}
}

func isCompressLogFile(isCompress int) (bool, error) {
	switch isCompress {

	case 0:
		return false, nil

	case 1:
		return true, nil

	default:
		ErrUnavailableCompressFlag := errors.New(UnavailableCompressFlag)
		return false, ErrUnavailableCompressFlag
	}
}
