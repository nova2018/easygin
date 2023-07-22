package logger

import (
	"fmt"
	"github.com/nova2018/easygin/config"
	"github.com/nova2018/easygin/file"
	"github.com/nova2018/easygin/utils"
	"github.com/nova2018/gologger/logger"
	"github.com/nova2018/gologger/zap"
	zap2 "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	loggerConsole      = "console"
	loggerDefault      = "default"
	loggerConsoleError = "consoleError"
)

func Init() {
	v := config.Sub("logger")
	if v == nil {
		panic("no logger config")
	}
	var cfg map[string][]*file.WriteConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	// init console
	setLogger(loggerConsole, initConsoleLogger())
	setLogger(loggerConsoleError, initConsoleErrorLogger())

	for name, c := range cfg {
		listCore := makeCores(name)
		l := newLogger(c, listCore...)
		setLogger(name, l)
	}
}

func newLogger(listCfg []*file.WriteConfig, cores ...zapcore.Core) logger.Logger {
	encodeConfig := zap.EncoderConfig()
	var err error
	listCore := make([]zapcore.Core, 0, len(listCfg)+1)
	for _, c := range listCfg {
		var lv zapcore.Level
		err = lv.Set(c.Level)
		if err != nil {
			panic(err)
		}
		write := getLogWriter(c)
		var encode zapcore.Encoder
		if strings.ToLower(c.OutFormat) == "json" {
			encode = zapcore.NewJSONEncoder(encodeConfig)
		} else {
			encode = zapcore.NewConsoleEncoder(encodeConfig)
		}
		core := zapcore.NewCore(encode, write, lv)
		listCore = append(listCore, core)
	}

	listCore = append(listCore, newLoggerAlarm())
	if len(cores) > 0 {
		listCore = append(listCore, cores...)
	}

	tee := zapcore.NewTee(listCore...)

	listOpt := []zap2.Option{
		zap2.AddCaller(),
	}
	if !utils.IsProd() {
		listOpt = append(listOpt, zap2.Development())
	}
	zapLogger := zap2.New(tee, listOpt...)

	return logger.NewLogger(zapLogger)
}

func newCliLogger() logger.Logger {
	return logger.NewLogger(zap2.NewExample())
}

func initConsoleLogger() logger.Logger {
	return newLogger([]*file.WriteConfig{{
		Level:    "debug",
		FileName: "stdout",
	}})
}

func initConsoleErrorLogger() logger.Logger {
	return newLogger([]*file.WriteConfig{{
		Level:    "debug",
		FileName: "stderr",
	}})
}

func getLogWriter(cfg *file.WriteConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(file.NewWriter(cfg))
}

func setLogger(name string, l logger.Logger) {
	fmt.Printf("logger init success! name=[%s]\n", name)
	_mapLogger[name] = l
}

func Default() logger.Logger {
	l := Logger(loggerDefault)
	if l == nil {
		return newCliLogger()
	}
	return l
}

func Console() logger.Logger {
	return Logger(loggerConsole)
}

func ConsoleError() logger.Logger {
	return Logger(loggerConsoleError)
}

func Logger(name string) logger.Logger {
	if v, ok := _mapLogger[name]; ok {
		return v
	}
	return nil
}

var (
	_mapLogger = map[string]logger.Logger{}
)

func Sync() {
	for _, l := range _mapLogger {
		_ = l.Sync()
	}
}
