package logger

import (
	"forever.love/global"
	viper2 "forever.love/initialize/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

/*
Zap提供了两种类型的日志记录器：SugaredLogger 和 Logger。
在性能很好但不是很关键的上下文中，使用 SugaredLogger。它比其他结构化日志记录包快4-10倍，并且支持结构化和 printf 风格的日志记录。
SugaredLogger 是在 Logger 的基础上，提供了更大、更便利的参数格式化方式。但是也可能会带来风险，也会造成额外的内存分配。
在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
建议：如果对参数格式化方式无太大要求，强烈建议使用 Logger 而不是 SugaredLogger。
*/

var Log *zap.SugaredLogger

func LogSetup(logLevel zapcore.Level) *zap.SugaredLogger {
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		CallerKey:   "file",
		EncodeLevel: zapcore.CapitalColorLevelEncoder, // 大写编码器带颜色
		//EncodeCaller: zapcore.ShortCallerEncoder, 		//采用短文件路径编码输出
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
			//enc.AppendString(t.Format("15:04:05"))
		}, // 输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	// 实现判断日志等级的 interface，自定义日志级别
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && lvl >= logLevel
	})
	//zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel。
	var core zapcore.Core
	// 如果配置文件的路径存在，就把日志写入到这里（用于线上）
	if viper2.LoggerConf.Console {
		// 获取 debug、info、error 日志文件的 io.Writer
		debugWriter := getLogWriter(viper2.LoggerConf.DebugLogPath + ".log")
		infoWriter := getLogWriter(viper2.LoggerConf.InfoLogPath + ".log")
		errorWriter := getLogWriter(viper2.LoggerConf.ErrorLogPath + ".log")
		// 最后创建具体的 Logger，NewConsoleEncoder 是非结构化输出
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(errorWriter), errorLevel),
		)
	} else {
		// 将日志输出到控制台
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(os.Stdout), logLevel)
	}
	//通过调用主logger的. Sugar()方法来获取一个SugaredLogger
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel), zap.AddCallerSkip(1)).Sugar()
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel), zap.AddCallerSkip(1)).Sugar()
	return Log
}

// SugaredLogger 和 rotatelogs 实现日志按天分割
//func getLogWriter(filename string) io.Writer {
//	rotateLogger, err := rotatelogs.New(
//		strings.Replace(filename, ".log", "", -1) + "-%Y-%m-%d.log",
//		// WithMaxAge（日志有效时长） 和 WithRotationCount（日志保存个数） 只能设置其中一个
//		rotatelogs.WithMaxAge(time.Hour*24*30),
//		//rotatelogs.WithRotationCount(30),
//		// 日志按天分割
//		rotatelogs.WithRotationTime(time.Hour*24),
//	)
//	if err != nil {
//		panic(err)
//	}
//	return rotateLogger
//}

// SugaredLogger 和 lumberjack 实现日志按大小分割
func getLogWriter(filename string) io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,   // 日志文件最大大小（以MB为单位）
		MaxBackups: 10,    // 保留旧文件最大个数
		MaxAge:     1,     // 保留旧文件最大天数
		Compress:   false, // 是否压缩、归档旧文件
		LocalTime:  true,
	}
	return lumberJackLogger
}

func GetApiRoutePathAndLine(skip int) (apiRoute, fileName, lineStr string) {
	goroutineID := runtime.Goid()
	apiRouteA, isExist := global.GoIDMap.Load(goroutineID)
	if isExist == false {
		apiRoute = ""
	} else {
		apiRoute = apiRouteA.(string)
	}
	_, file, lineInt, ok := runtime.Caller(skip)
	if !ok {
		Log.Error("get info failed")
		return
	}
	lineStr = strconv.Itoa(lineInt)
	fileName = path.Base(file)
	return
}

func Debug(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.Debugf("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}

func Info(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.Info(args...)
}

func Infof(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.Infof("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}

func Warn(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.Warnf("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}

func Error(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.Errorf("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}

func DPanic(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.DPanicf("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}

func Panic(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.Panicf("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}

func Fatal(args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	apiRoutePath := "，apiRoute：" + apiRoute + "，caller：" + fileName + ":" + line
	args = append(args, apiRoutePath)
	Log.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	apiRoute, fileName, line := GetApiRoutePathAndLine(2)
	Log.Fatalf("apiRoute："+apiRoute+"，caller："+fileName+":"+line+" "+template, args...)
}
