package tlog

//func Init() {
//
//}
//package tlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var Logger *zap.Logger
var log *zap.SugaredLogger

// Init 日志文件初始化
func Init() {
	file, _ := os.Create("./server.log")

	// 利用io.MultiWriter支持文件和终端两个输出目标

	ws := io.MultiWriter(file, os.Stdout)

	ztarget := zapcore.AddSync(ws)

	zenCoder := getEncoder()

	zCore := zapcore.NewCore(zenCoder, ztarget, zapcore.DebugLevel)

	// 当我们不是直接使用初始化好的logger实例记录日志，而是将其包装成一个函数等，此时日录日志的函数调用链会增加，
	// 想要获得准确的调用信息就需要通过AddCallerSkip函数来跳过
	Logger = zap.New(zCore, zap.AddCaller(), zap.AddCallerSkip(1))

	log = Logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	// 修改默认的日志配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	log.Panicf(template, args...)
}
