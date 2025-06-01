package logger

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(lv logrus.Level) {
	Log = logrus.New()
	Log.SetLevel(lv)
	ConfigLocalFilesystemLogger("./log/", "std.log", time.Hour*24*7, time.Hour*24, lv)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration, lv logrus.Level) {
	baseLogPath := path.Join(logPath, logFileName)
	normalWriter, err := rotatelogs.New(
		path.Join(logPath, "%Y%m%d/", "%Y%m%d%H%M%S."+logFileName),
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔LLL
	)

	if err != nil {
		Log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	baseLogPath1 := path.Join(logPath, "Error_"+logFileName)
	errorWriter, err1 := rotatelogs.New(
		path.Join(logPath, "%Y%m%d/", "Error_%Y%m%d%H%M%S."+logFileName),
		rotatelogs.WithLinkName(baseLogPath1),     // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔LLL
	)
	if err1 != nil {
		Log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	filenameHook := NewHook()
	filenameHook.Field = "\nline"
	Log.AddHook(filenameHook)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: normalWriter, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  normalWriter,
		logrus.WarnLevel:  normalWriter,
		logrus.ErrorLevel: errorWriter,
		logrus.FatalLevel: errorWriter,
		logrus.PanicLevel: errorWriter,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000", FullTimestamp: true})
	//注意上面这个 and &amp;符号被转义了
	Log.AddHook(lfHook)

	lfHookErr := lfshook.NewHook(lfshook.WriterMap{
		logrus.ErrorLevel: normalWriter,
		logrus.FatalLevel: normalWriter,
		logrus.PanicLevel: normalWriter,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000", FullTimestamp: true})
	Log.AddHook(lfHookErr)
}
