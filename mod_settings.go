package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type mod_settings struct {
	cfg *Config
}

func (m *mod_settings) init(config *Config) {
	// Set config
	m.cfg = config
}

func (m *mod_settings) getall(c echo.Context) error {
	type Data struct {
		AdminUser     string `json:"AdminUser"`
		ScoopEnabled  bool   `json:"ScoopEnabled"`
		ScoopDir      string `json:"ScoopDir"`
		ScoopInitRepo string `json:"ScoopInitRepo"`
		AptlyEnabled  bool   `json:"repos"`
		AptlyLocal    bool   `json:"AptlyLocal"`
		AptlyURL      string `json:"AptlyURL"`
	}
	data := &Data{
		AdminUser:     m.cfg.AdminUser,
		ScoopEnabled:  m.cfg.ScoopEnabled,
		ScoopDir:      m.cfg.ScoopDir,
		ScoopInitRepo: m.cfg.ScoopInitRepo,
		AptlyEnabled:  m.cfg.AptlyEnabled,
		AptlyLocal:    m.cfg.AptlyLocal,
		AptlyURL:      m.cfg.AptlyURL,
	}
	return c.JSON(http.StatusOK, data)
}

func (m *mod_settings) init_route(loc *echo.Group) {
	loc.POST("/set_auth", func(c echo.Context) error {
		user := c.QueryParam("username")
		if user == "" {
			return c.String(http.StatusInternalServerError, "No username param.")
		}
		password := c.QueryParam("password")
		if password == "" {
			return c.String(http.StatusInternalServerError, "No password param.")
		}
		hasher := sha256.New()
		hasher.Write([]byte(password))
		sha256 := fmt.Sprintf("%x", (hasher.Sum(nil)))
		m.cfg.set_auth(user, sha256)
		return c.String(http.StatusOK, "ok")
	})
	loc.POST("/set_scoop", func(c echo.Context) error {
		ScoopEnabled := c.QueryParam("ScoopEnabled")
		if ScoopEnabled == "" {
			return c.String(http.StatusInternalServerError, "No ScoopEnabled param.")
		}
		if ScoopEnabled != "true" && ScoopEnabled != "false" {
			return c.String(http.StatusInternalServerError, "Incorrect ScoopInitRepo param.")
		}
		ScoopDir := c.QueryParam("ScoopDir")
		if ScoopDir == "" {
			return c.String(http.StatusInternalServerError, "No ScoopDir param.")
		}
		ScoopInitRepo := c.QueryParam("ScoopInitRepo")
		if ScoopInitRepo == "" {
			return c.String(http.StatusInternalServerError, "No ScoopInitRepo param.")
		}
		m.cfg.set_scoop(ScoopEnabled, ScoopDir, ScoopInitRepo)
		return c.String(http.StatusOK, "ok")
	})
	loc.GET("/get_all", m.getall)
}
