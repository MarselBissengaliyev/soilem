package handlers

import (
	"net/http"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) registration(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	foundUser, err := h.services.User.Registration(&user)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if err = h.services.Twilo.SendSMSConfirmation(foundUser); err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	expiresAt := time.Now().Add(120 * time.Second)
	sessionToken := h.services.Session.CreateSession(&model.Session{
		UserName:  foundUser.UserName,
		Expiry:    expiresAt,
		UserAgent: ctx.Request.UserAgent(),
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	newDataResponse(ctx, http.StatusCreated, dataResponse{
		"message": "registration successfully completed",
	})
}

func (h *Handler) login(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	foundUser, err := h.services.User.Login(&user)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	expiresAt := time.Now().Add(120 * time.Second)
	sessionToken := h.services.Session.CreateSession(&model.Session{
		UserName:  foundUser.UserName,
		Expiry:    expiresAt,
		UserAgent: ctx.Request.UserAgent(),
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	newDataResponse(ctx, http.StatusCreated, dataResponse{
		"message": "login successfully completed",
	})
}

func (h *Handler) logout(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	h.services.Session.RemoveSession(sessionToken.(string))

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

func (h *Handler) getUserByUserName(ctx *gin.Context) {
	userName := ctx.Param("user_name")

	foundUser, err := h.services.User.GetUserByUserName((model.UserName(userName)))
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"user": foundUser,
	})
}

func (h *Handler) getUsers(ctx *gin.Context) {
	limit := ctx.Query("limit")
	searchTerm := ctx.Query("search_term")

	users, err := h.services.User.GetUsers(searchTerm, limit)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"users": users,
	})
}
