package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"path"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	cfg *Config
	api *echo.Group
	e   *echo.Echo
}

// Generate token for auth from curl
func (s *Server) login(c echo.Context) error {
	user := c.FormValue("username")
	pass := sha256.Sum256([]byte(c.FormValue("password")))

	if subtle.ConstantTimeCompare([]byte(user), []byte(s.cfg.AdminUser)) != 1 ||
		subtle.ConstantTimeCompare([]byte(hex.EncodeToString(pass[:])), []byte(s.cfg.AdminPassword)) != 1 {
		return echo.ErrUnauthorized
	}

	// Create token with claims
	token := jwt.New(jwt.SigningMethodHS256)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// Generate token for auth from browser
func (s *Server) login_from_page(c echo.Context) error {
	user := c.QueryParam("username")
	pass := sha256.Sum256([]byte(c.QueryParam("password")))

	if subtle.ConstantTimeCompare([]byte(user), []byte(s.cfg.AdminUser)) != 1 ||
		subtle.ConstantTimeCompare([]byte(hex.EncodeToString(pass[:])), []byte(s.cfg.AdminPassword)) != 1 {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    t,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusTemporaryRedirect, "/admin")
}

func (s *Server) init(config *Config) {
	// Set config
	s.cfg = config
	s.e = echo.New()

	// Middleware
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	// All static files and login page are public
	// We serve them to correctly show login page
	s.e.Static("/static", path.Join(s.cfg.ClientPath, "static"))
	s.e.File("/favicon.ico", path.Join(s.cfg.ClientPath, "favicon.ico"))
	s.e.Static("/login", path.Join(s.cfg.ClientPath, "login"))

	// Login route
	s.e.POST("/api/login", s.login)          // for REST API
	s.e.GET("/api/login", s.login_from_page) // for Browser

	// Redirect if unauthenticated from root to login page
	s.e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/login")
	})

	// Protected group
	r := s.e.Group("/")
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:Authorization",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusPermanentRedirect, "/login")
		},
	}))

	// Map admin
	r.Static("admin", path.Join(s.cfg.ClientPath, "admin"))
	// Map /api route
	s.api = r.Group("api/")
}

func (s *Server) start() {
	s.e.Logger.Fatal(s.e.Start(":" + s.cfg.Port))
}
