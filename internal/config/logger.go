package config

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func GetRequestLoggerConfig() middleware.RequestLoggerConfig {
    return middleware.RequestLoggerConfig{
        LogStatus:   true,
        LogURI:      true,
        LogMethod:   true,
        LogLatency:  true,
        HandleError: true,
        LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
            attrs := []slog.Attr{
                slog.String("method", v.Method),
                slog.String("uri", v.URI),
                slog.Int("status", v.Status),
            }

            if v.Error != nil {
                attrs = append(attrs, slog.String("err", v.Error.Error()))
                slog.LogAttrs(c.Request().Context(), slog.LevelError, "request failed", attrs...)
            } else {
                attrs = append(attrs, slog.Duration("latency", v.Latency))
                slog.LogAttrs(c.Request().Context(), slog.LevelInfo, "request processed", attrs...)
            }

            return nil
        },
    }
}