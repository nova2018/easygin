package logger

import (
	"go.uber.org/zap/zapcore"
)

type alarmCore struct {
	fields []zapcore.Field
}

func (a *alarmCore) Enabled(level zapcore.Level) bool {
	return level >= zapcore.ErrorLevel
}

func (a *alarmCore) With(fields []zapcore.Field) zapcore.Core {
	x := a.clone()
	x.fields = append(x.fields, fields...)
	return x
}

func (a *alarmCore) clone() *alarmCore {
	core := newLoggerAlarm()
	core.fields = a.fields
	return core
}

func (a *alarmCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if a.Enabled(entry.Level) {
		ce.AddCore(entry, a)
	}
	return ce
}

func (a *alarmCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// @todo
	//encode := zapcore.NewMapObjectEncoder()
	//for _, f := range a.fields {
	//	f.AddTo(encode)
	//}
	//for _, f := range fields {
	//	f.AddTo(encode)
	//}
	//b, _ := json.Marshal(encode.Fields)
	//
	//messageBody := strings.SplitN(entry.Message, zap.FieldSeparator, 3)
	//traceId := goutils.Substring(messageBody[1], 1, -1)
	//message := goutils.Substring(messageBody[2], 1, -1)
	//
	//event.Dispatch(alarm.NewSimpleEvent(
	//	fmt.Sprintf("日志报警 - %s %s", entry.Level.CapitalString(), entry.LoggerName),
	//	fmt.Sprintf("日志内容：%s", message),
	//	"",
	//	string(b),
	//	traceId,
	//	"default",
	//))

	return nil
}

func (a *alarmCore) Sync() error {
	return nil
}

func newLoggerAlarm() *alarmCore {
	return &alarmCore{
		fields: make([]zapcore.Field, 0),
	}
}
