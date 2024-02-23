package log

import (
	"io"
	"log/slog"
)

type LogWriter struct {
	file string
}

func NewLogWriter() io.Writer {
	return &LogWriter{
		file: "stdout",
	}
}

func NewErrorLogWriter() io.Writer {
	return &LogWriter{
		file: "stderr",
	}
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	slog.Debug(string(p), "file", l.file)
	return len(p), nil
}
