package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tchaudhry91/algomon/store"
)

type APIServer struct {
	e      *echo.Echo
	db     *store.BoltStore
	config *Config
	logger *slog.Logger
}

func NewAPIServer(db *store.BoltStore, config *Config, logger *slog.Logger) *APIServer {
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
	e.Use(middleware.CORS())
	server := APIServer{
		e:      e,
		db:     db,
		logger: logger,
		config: config,
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
	s.e.GET("/api/v1/checks/:name", s.getNamedCheck)
	s.e.GET("/api/v1/checks/:name/failures", s.getNamedCheckFailures)
}

func (s *APIServer) getChecksStatus(c echo.Context) error {
	data, err := s.db.GetChecksStatus(c.Request().Context())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.JSON(http.StatusOK, []string{})
		}
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (s *APIServer) getNamedCheck(c echo.Context) error {
	name := c.Param("name")
	data, err := s.db.GetNamedCheck(c.Request().Context(), name, 5)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "check not found"})
		}
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (s *APIServer) getNamedCheckFailures(c echo.Context) error {
	name := c.Param("name")
	data, err := s.db.GetNamedCheckFailures(c.Request().Context(), name, 5)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "check not found"})
		}
		return err
	}
	return c.JSON(http.StatusOK, data)
}
