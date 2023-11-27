package middleware

import (
	"net/http"

	"github.com/MarselBissengaliyev/soilem/internal/handlers/response"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) Authenticate(ctx *gin.Context) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := ctx.Request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			response.NewErrorResponse(ctx, http.StatusUnauthorized, "session_token cookie is not found: "+err.Error())
			return
		}

		response.NewErrorResponse(ctx, http.StatusUnauthorized, "Session token is not set: "+err.Error())
		return
	}

	sessionToken := c.Value
	foundSession, fail := m.services.Session.GetByAccessToken(sessionToken)
	if fail != nil {
		response.NewErrorResponse(ctx, fail.StatusCode, fail.Message)
		return
	}

	if foundSession.IsExpired() {
		response.NewErrorResponse(ctx, http.StatusUnauthorized, "session token is expired")
		return
	}

	ctx.Set("session_token", sessionToken)

	ctx.Next()
}
