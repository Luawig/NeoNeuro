package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Luawig/neoneuro/backend/internal/handlers"
	"github.com/Luawig/neoneuro/backend/pkg/config"
	"github.com/Luawig/neoneuro/backend/pkg/logger"
	mw "github.com/Luawig/neoneuro/backend/pkg/middleware"
	mwauth "github.com/Luawig/neoneuro/backend/pkg/middleware/auth"
)

func NewEngine(cfg config.Config) *gin.Engine {
	// 单一模式，不区分 prod/dev
	logger.Init(cfg)

	r := gin.New()
	r.Use(mw.RequestID(), mw.Recovery(), mw.AccessLog(), cors.Default())

	// public
	r.GET("/api/v1/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// auth group
	auth := r.Group("/api/v1")
	auth.Use(mwauth.JWTAuth(cfg))
	{
		auth.GET("/me", handlers.Me)
		auth.GET("/admin/stats", mwauth.RequireRoles("admin"), handlers.AdminStats)
		auth.GET("/models", mwauth.RequireScopes("models:read"), handlers.ListModels)
	}

	return r
}
