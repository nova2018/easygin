package logger

import "go.uber.org/zap/zapcore"

type CoreFactory func(name string) zapcore.Core

var (
	listFactory = make([]CoreFactory, 0)
)

func RegisterCoreFactory(factory CoreFactory) {
	listFactory = append(listFactory, factory)
}

func makeCores(name string) []zapcore.Core {
	listCore := make([]zapcore.Core, 0, len(listFactory))
	for _, f := range listFactory {
		c := f(name)
		if c != nil {
			listCore = append(listCore, c)
		}
	}
	return listCore
}
