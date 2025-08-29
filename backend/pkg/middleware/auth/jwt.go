package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/Luawig/neoneuro/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub    string   `json:"sub"`
	Roles  []string `json:"roles,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

const CtxClaimsKey = "auth_claims"

// JWTAuth validates Bearer token (HS256) and injects Claims into context.
func JWTAuth(cfg config.Config) gin.HandlerFunc {
	secret := cfg.JWT.Secret
	if strings.ToUpper(cfg.JWT.Alg) != "HS256" {
		panic("only HS256 implemented in this scaffold")
	}
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimSpace(h[7:])
		var cl Claims
		parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))
		tok, err := parser.ParseWithClaims(raw, &cl, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// iss/aud/exp
		now := time.Now()
		if cl.ExpiresAt != nil && now.After(cl.ExpiresAt.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}
		if iss := cfg.JWT.Issuer; iss != "" && cl.Issuer != iss {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad issuer"})
			return
		}
		if aud := cfg.JWT.Audience; aud != "" && !containsAudience(cl.Audience, aud) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad audience"})
			return
		}
		c.Set(CtxClaimsKey, cl)
		c.Next()
	}
}

func containsAudience(aud jwt.ClaimStrings, want string) bool {
	for _, v := range aud {
		if v == want {
			return true
		}
	}
	return false
}

func GetClaims(c *gin.Context) (Claims, bool) {
	v, ok := c.Get(CtxClaimsKey)
	if !ok {
		return Claims{}, false
	}
	cl, _ := v.(Claims)
	return cl, true
}
