package file

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
	"strings"
	"time"
)

type writerWrapper struct {
	Writer io.Writer
}

func (w *writerWrapper) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	if err != nil {
		fmt.Printf("%s writerWrapper.Write failure! err=%v\n", time.Now(), err)
	}
	return
}

func NewWriter(cfg *WriteConfig) io.Writer {
	consoleWriter := getConsoleWriter(cfg)
	if consoleWriter != nil {
		return consoleWriter
	}
	if cfg.Driver == "" && cfg.MaxFiles == 0 {
		cfg.Driver = "rotate"
	}
	switch strings.ToLower(cfg.Driver) {
	case "rotate", "rotatelogs":
		return &writerWrapper{
			getRotateWriter(cfg),
		}
	case "lumberjack":
		fallthrough
	default:
		return &writerWrapper{
			getLumberjackLogWriter(cfg),
		}
	}
}

func getConsoleWriter(cfg *WriteConfig) io.Writer {
	var file io.Writer
	switch strings.ToLower(cfg.FileName) {
	case "stdout":
		fallthrough
	case "/dev/stdout":
		file = os.Stdout
	case "stdin":
		fallthrough
	case "/dev/stdin":
		file = os.Stdin
	case "stderr":
		fallthrough
	case "/dev/stderr":
		file = os.Stderr
	default:
		return nil
	}

	return file
}

func getRotateWriter(cfg *WriteConfig) io.Writer {
	opts := make([]rotatelogs.Option, 0, 5)
	opts = append(opts, rotatelogs.WithRotationTime(time.Hour*24))
	opts = append(opts, rotatelogs.WithLinkName(cfg.FileName))
	if cfg.MaxAge > 0 {
		opts = append(opts, rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour))
	} else if cfg.MaxFiles > 0 {
		opts = append(opts, rotatelogs.WithRotationCount(uint(cfg.MaxFiles)))
	}
	if cfg.MaxSize > 0 {
		opts = append(opts, rotatelogs.WithRotationSize(int64(cfg.MaxSize)))
	}
	if cfg.Compress {
		opts = append(opts, rotatelogs.WithHandler(rotatelogs.HandlerFunc(func(event rotatelogs.Event) {
			if event.Type() != rotatelogs.FileRotatedEventType {
				return
			}
			rotateCompress(event.(*rotatelogs.FileRotatedEvent).PreviousFile())
		})))
	}
	logf, err := rotatelogs.New(
		genRotateFile(cfg.FileName),
		opts...,
	)
	if err != nil {
		panic(err)
	}

	return logf
}

func rotateCompress(file string) {
	src := file
	dst := fmt.Sprintf("%s.gz", file)
	_ = compressLogFile(src, dst)
}

func genRotateFile(fileName string) string {
	sep := "."
	arr := strings.Split(fileName, sep)
	target := len(arr) - 2
	arr[target] = arr[target] + "-%Y-%m-%d"
	return strings.Join(arr, sep)
}

func getLumberjackLogWriter(cfg *WriteConfig) io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxFiles,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}
	lumberjackLogRotate(lumberJackLogger)
	return lumberJackLogger
}

func lumberjackLogRotate(lumberJackLogger *lumberjack.Logger) {
	go func() {
		for {
			y, m, d := time.Now().Date()
			tomorrow := time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
			<-time.After(tomorrow.Sub(time.Now()) - 1*time.Nanosecond)
			_ = lumberJackLogger.Rotate()
		}
	}()
	if info, err := os.Stat(lumberJackLogger.Filename); !os.IsNotExist(err) {
		y, m, d := time.Now().Date()
		today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
		if info.ModTime().Before(today) {
			// 最后修改时间在今天之前，则需要切割
			_ = lumberJackLogger.Rotate()
		}
	}
}
