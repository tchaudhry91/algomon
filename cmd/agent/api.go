package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tchaudhry91/algoprom/store"
)

type APIServer struct {
	e      *echo.Echo
	db     *store.BoltStore
	logger *slog.Logger
}

func NewAPIServer(db *store.BoltStore, logger *slog.Logger) *APIServer {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())
	server := APIServer{
		e:      e,
		db:     db,
		logger: logger,
	}
	server.Routes()
	server.logger.Info("Registered Routes!")
	return &server
}

func (s *APIServer) Mux() http.Handler {
	return s.e
}

func (s *APIServer) Routes() {
	s.e.GET("/api/v1/checks", s.getChecksStatus)
}

func (s *APIServer) getChecksStatus(c echo.Context) error {
	data, err := s.db.GetChecksStatus(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}
