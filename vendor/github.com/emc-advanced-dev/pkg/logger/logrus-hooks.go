package logger

import "github.com/Sirupsen/logrus"

type LoggerNameHook struct {
	Name string
}

func (h *LoggerNameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LoggerNameHook) Fire(entry *logrus.Entry) error {
	entry = entry.WithField("origin", h.Name)
	return nil
}
