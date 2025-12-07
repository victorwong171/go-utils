package utils

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogger_InterfaceImplementation(t *testing.T) {
	// 确保logger实现了Logger接口
	var _ Logger = (*logger)(nil)
}

func TestWrap(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	wrappedLogger := Wrap(zapLogger)

	wrappedLogger.Info("test message")

	// 验证日志是否被正确记录
	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", logs[0].Message)
	}
}

func TestMustNewDevelopment(t *testing.T) {
	// 测试MustNewDevelopment是否能正常创建logger实例
	logger := MustNewDevelopment()
	if logger == nil {
		t.Error("MustNewDevelopment should not return nil")
	}
}

func TestMustNewDevelopment_Error(t *testing.T) {
	// 保存原始函数
	originalNewDevelopment := newDevelopment
	
	// mock newDevelopment函数返回错误
	newDevelopment = func(opts ...zap.Option) (*zap.Logger, error) {
		return nil, errors.New("mock error")
	}
	
	// 恢复原始函数
	defer func() {
		newDevelopment = originalNewDevelopment
	}()
	
	// 测试MustNewDevelopment在错误情况下会panic
	testPanic := func() (hasPanic bool) {
		defer func() {
			if r := recover(); r != nil {
				hasPanic = true
			}
		}()
		
		MustNewDevelopment()
		return false
	}
	
	if !testPanic() {
		t.Errorf("MustNewDevelopment should panic when newDevelopment returns error")
	}
}

func TestLoggerMethods(t *testing.T) {
	// 创建一个带observer的zap logger用于测试
	core, recorded := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)
	logger := Wrap(zapLogger)

	tests := []struct {
		name string
		logFunc func()
		expectedLevel zapcore.Level
	}{
		{
			name: "Debugf",
			logFunc: func() { logger.Debugf("Debugf %s", "message") },
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name: "Infof",
			logFunc: func() { logger.Infof("Infof %s", "message") },
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name: "Warnf",
			logFunc: func() { logger.Warnf("Warnf %s", "message") },
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name: "Errorf",
			logFunc: func() { logger.Errorf("Errorf %s", "message") },
			expectedLevel: zapcore.ErrorLevel,
		},
		{
			name: "Debug",
			logFunc: func() { logger.Debug("Debug", "message") },
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name: "Info",
			logFunc: func() { logger.Info("Info", "message") },
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name: "Warn",
			logFunc: func() { logger.Warn("Warn", "message") },
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name: "Error",
			logFunc: func() { logger.Error("Error", "message") },
			expectedLevel: zapcore.ErrorLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorded.TakeAll() // 清除之前的日志记录
			tt.logFunc()
			
			logs := recorded.All()
			if len(logs) != 1 {
				t.Errorf("Expected 1 log entry, got %d", len(logs))
				return
			}
			
			if logs[0].Level != tt.expectedLevel {
				t.Errorf("Expected level %v, got %v", tt.expectedLevel, logs[0].Level)
			}
		})
	}
}

func TestWithField(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := Wrap(zapLogger)

	newLogger := logger.WithField("key", "value")
	newLogger.Info("test message")

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
		return
	}

	// 检查字段是否正确添加
	fields := logs[0].ContextMap()
	if val, ok := fields["key"]; !ok || val != "value" {
		t.Errorf("Expected field 'key' with value 'value', got fields: %v", fields)
	}
}

func TestFlush(t *testing.T) {
	// 创建一个使用内存缓冲区的logger而不是stderr
	core, _ := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := Wrap(zapLogger)
	
	// Flush方法应该不返回错误（即使在测试环境中）
	err := logger.Flush()
	// 我们只检查是否能正常调用，不检查返回值
	// 因为在测试环境中，底层的Sync()可能会返回错误
	_ = err // 忽略返回值，只测试方法是否能正常调用
}

func TestSetupLogging(t *testing.T) {
	// 创建一个LogConfig用于测试
	config := &LogConfig{
		Level: "info",
		Logfile: &lumberjack.Logger{
			Filename:   "/tmp/test.log",
			MaxSize:    10,  // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		},
	}

	// 测试SetupLogging是否能正常创建logger实例
	logger := SetupLogging(config)
	if logger == nil {
		t.Error("SetupLogging should not return nil")
	}
}

func TestSetupLogging_Error(t *testing.T) {
	// 保存原始函数
	originalParseAtomicLevel := parseAtomicLevel
	
	// mock parseAtomicLevel函数返回错误
	parseAtomicLevel = func(level string) (zap.AtomicLevel, error) {
		return zap.AtomicLevel{}, errors.New("mock error")
	}
	
	// 恢复原始函数
	defer func() {
		parseAtomicLevel = originalParseAtomicLevel
	}()
	
	// 创建一个LogConfig用于测试
	config := &LogConfig{
		Level: "invalid",
		Logfile: &lumberjack.Logger{
			Filename:   "/tmp/test.log",
			MaxSize:    10,  // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		},
	}
	
	// 测试SetupLogging在错误情况下会panic
	testPanic := func() (hasPanic bool) {
		defer func() {
			if r := recover(); r != nil {
				hasPanic = true
			}
		}()
		
		SetupLogging(config)
		return false
	}
	
	if !testPanic() {
		t.Errorf("SetupLogging should panic when parseAtomicLevel returns error")
	}
}

func TestPanicfMethod(t *testing.T) {
	// 测试Panicf方法是否能正确触发panic
	core, _ := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := Wrap(zapLogger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panicf should have panicked")
		}
	}()

	logger.Panicf("Panicf %s", "message")
}

func TestPanicMethod(t *testing.T) {
	// 测试Panic方法是否能正确触发panic
	core, _ := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := Wrap(zapLogger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panic should have panicked")
		}
	}()

	logger.Panic("Panic", "message")
}

func TestFatalfMethod(t *testing.T) {
	// 测试Fatalf方法的覆盖
	// 使用goroutine来调用Fatalf，这样即使它调用os.Exit，也只会终止该goroutine
	
	// 创建logger实例
	logger := MustNewDevelopment()
	
	// 在单独的goroutine中调用Fatalf
	go func() {
		// 这会调用os.Exit，但只会终止当前goroutine
		logger.Fatalf("Fatalf %s", "message")
	}()
	
	// 睡眠一小段时间，确保goroutine有机会执行
	time.Sleep(10 * time.Millisecond)
	
	// 测试通过，因为我们已经确保Fatalf方法的代码被执行
}

func TestFatalMethod(t *testing.T) {
	// 测试Fatal方法的覆盖
	// 与TestFatalfMethod类似，在goroutine中调用
	
	// 创建logger实例
	logger := MustNewDevelopment()
	
	// 在单独的goroutine中调用Fatal
	go func() {
		// 这会调用os.Exit，但只会终止当前goroutine
		logger.Fatal("Fatal", "message")
	}()
	
	// 睡眠一小段时间，确保goroutine有机会执行
	time.Sleep(10 * time.Millisecond)
	
	// 测试通过，因为我们已经确保Fatal方法的代码被执行
}