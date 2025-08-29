package handlers

import (
	"net/http"

	mwauth "github.com/Luawig/neoneuro/backend/pkg/middleware/auth"
	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	if cl, ok := mwauth.GetClaims(c); ok {
		c.JSON(http.StatusOK, gin.H{"sub": cl.Sub, "roles": cl.Roles, "scopes": cl.Scopes})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func AdminStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"uptime": "ok", "qps": 42})
}

func ListModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"models": []string{"gpt-x", "tts-y", "emo-z"}})
}
