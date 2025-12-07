package utils

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 定义日志记录器接口
type Logger interface {
	WithField(key string, value interface{}) Logger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	// Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	// Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	// Print(args ...interface{})
	Warn(args ...interface{})
	// Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Flush() error
}

// isTesting 检查当前是否在测试环境中
var isTesting = func() bool {
	for _, arg := range os.Args {
		if arg == "-test.run" || arg == "-test.v" || arg == "-test.cover" {
			return true
		}
	}
	return testing.Testing()
}()

// 定义函数类型以便于测试时进行mock
type (
	newDevelopmentFunc func(...zap.Option) (*zap.Logger, error)
	parseAtomicLevelFunc func(string) (zap.AtomicLevel, error)
)

// 默认实现
var (
	newDevelopment newDevelopmentFunc = func(opts ...zap.Option) (*zap.Logger, error) {
		return zap.NewDevelopment(opts...)
	}
	parseAtomicLevel parseAtomicLevelFunc = zap.ParseAtomicLevel
)

type logger struct {
	l *zap.SugaredLogger
}

func Wrap(l *zap.Logger) *logger {
	l = l.WithOptions(zap.AddCallerSkip(1))
	return &logger{l: l.Sugar()}
}

func MustNewDevelopment() *logger {
	l, err := newDevelopment()
	if err != nil {
		panic("slog_zap: zap.NewDevelopment(): " + err.Error())
	}
	return Wrap(l)
}

func (s *logger) WithField(key string, value interface{}) Logger {
	return &logger{l: s.l.With(key, value)}
}

func (s *logger) Debugf(format string, args ...interface{}) {
	s.l.Debugf(format, args...)
}

func (s *logger) Infof(format string, args ...interface{}) {
	s.l.Infof(format, args...)
}

func (s *logger) Warnf(format string, args ...interface{}) {
	s.l.Warnf(format, args...)
}

func (s *logger) Errorf(format string, args ...interface{}) {
	s.l.Errorf(format, args...)
}

func (s *logger) Fatalf(format string, args ...interface{}) {
	if isTesting {
		// 在测试环境中，使用Errorf代替Fatalf，避免调用os.Exit
		s.l.Errorf(format, args...)
	} else {
		s.l.Fatalf(format, args...)
	}
}

func (s *logger) Panicf(format string, args ...interface{}) {
	s.l.Panicf(format, args...)
}

func (s *logger) Debug(args ...interface{}) {
	s.l.Debug(args...)
}

func (s *logger) Info(args ...interface{}) {
	s.l.Info(args...)
}

func (s *logger) Warn(args ...interface{}) {
	s.l.Warn(args...)
}

func (s *logger) Error(args ...interface{}) {
	s.l.Error(args...)
}

func (s *logger) Fatal(args ...interface{}) {
	if isTesting {
		// 在测试环境中，使用Error代替Fatal，避免调用os.Exit
		s.l.Error(args...)
	} else {
		s.l.Fatal(args...)
	}
}

func (s *logger) Panic(args ...interface{}) {
	s.l.Panic(args...)
}

func (s *logger) Flush() error {
	return s.l.Sync()
}

/*
日志配置
level: 日志级别，可选值有：debug, info, warn, error, fatal, panic
logfile: 日志文件配置，包含以下字段：

	filename: 日志文件名
	maxsize: 每个日志文件的最大大小（MB）
	maxage: 保留旧文件的最大天数
	maxbackups: 保留旧文件的最大数量
	localtime: 是否使用本地时间戳
	compress: 是否压缩旧文件
*/
type LogConfig struct {
	Level   string             `yaml:"level"`
	Logfile *lumberjack.Logger `yaml:"logfile"`
}

// SetupLogging 配置日志记录器
func SetupLogging(lc *LogConfig) Logger {
	w := zapcore.AddSync(lc.Logfile)
	level, err := parseAtomicLevel(lc.Level)
	if err != nil {
		panic(err)
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		&level,
	)
	logger := zap.New(core, zap.AddCaller())
	return Wrap(logger)
}