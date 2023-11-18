package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Authenticate(ctx *gin.Context) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := ctx.Request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			ctx.Status(http.StatusUnauthorized)
			return
		}

		ctx.Status(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value
	_, ok := m.services.Session.GetSession(sessionToken)
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	

	ctx.Set("session_token", sessionToken)

	ctx.Next()
}


