package middleware

import (
	"auth_service/utils"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"time"
)

type ILog interface {
	Start() fiber.Handler
}

type log struct {
	logInstance *slog.Logger
}

func NewLogMiddleware(logInstance *slog.Logger) ILog {
	return &log{
		logInstance: logInstance,
	}
}

func (l *log) Start() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()
		err := ctx.Next()
		latency := time.Since(start)

		span := trace.SpanFromContext(ctx.UserContext())
		spanCtx := span.SpanContext()

		attrs := []slog.Attr{
			slog.String("method", ctx.Method()),
			slog.String("path", ctx.Path()),
			slog.Int("status", ctx.Response().StatusCode()),
			slog.String("latency", utils.FormatLatency(latency)),
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		}
		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			span.SetAttributes(attribute.String("error", err.Error()))
		}
		l.logInstance.LogAttrs(
			ctx.UserContext(),
			slog.LevelInfo,
			"incoming request",
			attrs...,
		)
		return err
	}
}
