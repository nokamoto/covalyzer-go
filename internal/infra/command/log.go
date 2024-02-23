package command

import (
	"bytes"
	"io"
	"log/slog"
)

type logWriter struct {
	file string
	buf  *bytes.Buffer
}

func newLogWriter(buf *bytes.Buffer) io.Writer {
	return &logWriter{
		file: "stdout",
		buf:  buf,
	}
}

func newErrorLogWriter() io.Writer {
	return &logWriter{
		file: "stderr",
	}
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	slog.Debug(string(p), "file", l.file)
	if l.buf != nil {
		return l.buf.Write(p)
	}
	return len(p), nil
}
