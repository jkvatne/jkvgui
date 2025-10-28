package main

import (
	"context"
	"io"
	"log"
	"log/slog"
)

type PatHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PatHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PatHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	timeStr := r.Time.Format("[15:05:05.000]")
	h.l.Println(timeStr, level, r.Message, fields)
	return nil
}

func NewPatHandler(
	out io.Writer,
	opts PatHandlerOptions,
) *PatHandler {
	h := &PatHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}
