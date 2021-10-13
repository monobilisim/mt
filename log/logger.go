package log

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"os"
)

type Params struct {
	Level string
	File  string
}

type Logger struct {
	*Params
	*logrus.Logger
}

func NewLogger(params *Params) (l *Logger) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
		PadLevelText:              true,
	})

	l = &Logger{
		Params: params,
		Logger: logger,
	}

	level := logrus.InfoLevel
	if params.Level == "debug" {
		level = logrus.DebugLevel
	} else if params.Level == "info" {
		level = logrus.InfoLevel
	} else if params.Level == "warn" {
		level = logrus.WarnLevel
	} else if params.Level == "error" {
		level = logrus.ErrorLevel
	}
	l.Logger.SetLevel(level)

	if params.File != "" {
		file, err := os.OpenFile(params.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Error(err)
			return
		}
		l.Logger.AddHook(&writer.Hook{
			Writer:    file,
			LogLevels: levelsAbove(level),
		})
	}

	return
}

func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args)
}

func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args)
}

func (l *Logger) Panic(args ...interface{}) {
	l.Logger.Panic(args)
}

func (l *Logger) DebugWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Debug(args)
}

func (l *Logger) InfoWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Info(args)
}

func (l *Logger) WarnWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Warn(args)
}

func (l *Logger) ErrorWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Error(args)
}

func (l *Logger) FatalWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Fatal(args)
}

func (l *Logger) PanicWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Panic(args)
}

func levelsAbove(level logrus.Level) []logrus.Level {
	levels := make([]logrus.Level, 0)

	levels = append(levels, logrus.PanicLevel)

	if level == logrus.FatalLevel {
		levels = append(levels, logrus.FatalLevel)
	}

	if level == logrus.ErrorLevel {
		levels = append(levels, logrus.FatalLevel)
		levels = append(levels, logrus.ErrorLevel)
	}

	if level == logrus.WarnLevel {
		levels = append(levels, logrus.ErrorLevel)
		levels = append(levels, logrus.WarnLevel)

	}

	if level == logrus.InfoLevel {
		levels = append(levels, logrus.ErrorLevel)
		levels = append(levels, logrus.WarnLevel)
		levels = append(levels, logrus.InfoLevel)
	}

	if level == logrus.DebugLevel {
		levels = append(levels, logrus.ErrorLevel)
		levels = append(levels, logrus.WarnLevel)
		levels = append(levels, logrus.DebugLevel)
	}

	return levels
}

func Fatal(args ...interface{})  {
	logrus.Fatal(args)
}
