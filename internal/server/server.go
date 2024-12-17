package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p40pmn/assignment-breed/internal/breed"
)

type Server struct {
	breedSvc *breed.Service
}

func NewServer(_ context.Context, svc *breed.Service) (*Server, error) {
	return &Server{
		breedSvc: svc,
	}, nil
}

func (s *Server) Install(e *echo.Echo, mws ...echo.MiddlewareFunc) error {
	if e == nil {
		return errors.New("echo is nil")
	}

	v1 := e.Group("/v1")

	v1.POST("/breed-inquiry", s.listBreeds)

	return nil
}

func (s *Server) listBreeds(c echo.Context) error {
	req := new(breed.BreedQuery)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Request body must be a valid JSON.",
			"status":  "BINDING_ERROR",
			"code":    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()
	breeds, err := s.breedSvc.ListBreeds(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, breeds)
}
