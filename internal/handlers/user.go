package handlers

import (
	"net/http"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	services *service.Service
}

func NewUserHandler(services *service.Service) *UserHandler {
	return &UserHandler{services}
}

func (h *UserHandler) registration(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	createdUser, tx, err := h.services.User.Registration(&user)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	user.Profile.Author = createdUser.UserName
	_, err = h.services.Profile.Create(&user.Profile)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	expiresAt := time.Now().Add(120 * time.Second)
	sessionToken := h.services.Session.CreateSession(&model.Session{
		UserName:  createdUser.UserName,
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

func (h *UserHandler) login(ctx *gin.Context) {
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

func (h *UserHandler) logout(ctx *gin.Context) {
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

func (h *UserHandler) getUserByUserName(ctx *gin.Context) {
	userName := ctx.Param("user_name")

	foundUser, err := h.services.User.GetByUserName((model.UserName(userName)))
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"user": foundUser,
	})
}

func (h *UserHandler) getUsers(ctx *gin.Context) {
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

func (h *UserHandler) sendSMSCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	userName, ok := h.services.Session.GetUserName(sessionToken.(string))
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	foundUser, err := h.services.User.GetByUserName((model.UserName(userName)))
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	foundUser.SMSCode.GenerateConfirmationCode(foundUser.UserName)

	smsCode, err := h.services.SMSCode.SetSMSCode(foundUser.SMSCode, foundUser.UserName)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if err := h.services.SMSCode.SendSMSConfirmation(foundUser.PhoneNumber, smsCode.Code); err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"message": "sms confirmation sended succefully",
	})
}

func (h *UserHandler) confirmSMSCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	userName, ok := h.services.Session.GetUserName(sessionToken.(string))
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var smsCode model.SMSCode

	if err := ctx.BindJSON(&smsCode); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	confirmed, err := h.services.User.ConfirmSMSCode(model.UserName(userName), smsCode)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if !confirmed {
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to confirm sms code")
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"message": "sms code succefully confirmed",
	})
}

func (h *UserHandler) sendEmailCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")

	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	userName, ok := h.services.Session.GetUserName(sessionToken.(string))
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	foundUser, err := h.services.User.GetByUserName((model.UserName(userName)))
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	foundUser.EmailCode.GenerateConfirmationCode(foundUser.UserName)

	emailCode, err := h.services.EmailCode.SetEmailCode(foundUser.EmailCode, foundUser.UserName)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if err := h.services.EmailCode.SendEmailCode(
		"../../templates/email_confirm.html",
		foundUser.Email,
		emailCode.Code,
	); err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"message": "email confirmation sended succefully",
	})
}

func (h *UserHandler) confirmEmailCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	userName, ok := h.services.Session.GetUserName(sessionToken.(string))
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var emailCode model.EmailCode

	if err := ctx.BindJSON(&emailCode); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	confirmed, err := h.services.User.ConfirmEmailCode(model.UserName(userName), emailCode)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if !confirmed {
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to confirm email code")
		return
	}

	newDataResponse(ctx, http.StatusOK, dataResponse{
		"message": "email code succefully confirmed",
	})
}

func (h *UserHandler) updatePasswordByEmailCode(ctx *gin.Context) {

}

func (h *UserHandler) updatePasswordBySMSCode(ctx *gin.Context) {

}

func (h *UserHandler) updateEmail(ctx *gin.Context) {

}

func (h *UserHandler) updatePhone(ctx *gin.Context) {

}

func (h *UserHandler) updateFullName(ctx *gin.Context) {

}
