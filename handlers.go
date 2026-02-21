package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s *Server) HandlePut(c *echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	err := s.Store.Put(key, value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"msg": "ok"})
}

func (s *Server) HandleGet(c *echo.Context) error {
	key := c.Param("key")

	value, err := s.Store.Get(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{key: value})
}

func (s *Server) HandleUpdate(c *echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	err := s.Store.Update(key, value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{key: value})
}

func (s *Server) HandleDelete(c *echo.Context) error {
	key := c.Param("key")

	err := s.Store.Delete(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}
