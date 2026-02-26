package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type GetResponse struct {
	Value               string  `json:"value"`
	HasTTL              bool    `json:"has_ttl"`
	TTLRemainingSeconds float64 `json:"ttl_remaining_seconds,omitempty"`
}

func (s *Server) HandlePut(c *echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")
	ttlStr := c.QueryParam("ttl")

	if ttlStr != "" {
		ttlSeconds, err := strconv.ParseFloat(ttlStr, 64)
		if err != nil || ttlSeconds <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ttl must be a positive number of seconds"})
		}
		ttl := time.Duration(ttlSeconds * float64(time.Second))
		if err := s.Store.PutWithTTL(key, value, ttl); err != nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
	} else {
		if err := s.Store.Put(key, value); err != nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{"msg": "ok"})
}

func (s *Server) HandleGet(c *echo.Context) error {
	key := c.Param("key")

	value, err := s.Store.Get(key)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	remaining, hasTTL, _ := s.Store.TTLRemaining(key)
	resp := GetResponse{
		Value:  value,
		HasTTL: hasTTL,
	}
	if hasTTL {
		resp.TTLRemainingSeconds = remaining.Seconds()
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) HandleGetTTL(c *echo.Context) error {
	key := c.Param("key")

	remaining, hasTTL, err := s.Store.TTLRemaining(key)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if !hasTTL {
		return c.JSON(http.StatusOK, map[string]any{
			"key":     key,
			"has_ttl": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"key":                   key,
		"has_ttl":               true,
		"ttl_remaining_seconds": remaining.Seconds(),
	})
}

func (s *Server) HandleUpdate(c *echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	if err := s.Store.Update(key, value); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{key: value})
}

func (s *Server) HandleDelete(c *echo.Context) error {
	key := c.Param("key")

	if err := s.Store.Delete(key); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}

func (s *Server) HandleSetTTL(c *echo.Context) error {
	key := c.Param("key")
	secondsStr := c.Param("seconds")

	seconds, err := strconv.ParseFloat(secondsStr, 64)
	if err != nil || seconds <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "seconds must be a positive number"})
	}

	ttl := time.Duration(seconds * float64(time.Second))
	if err := s.Store.SetTTL(key, ttl); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}
