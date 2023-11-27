package handlers

import (
	"net/http"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/handlers/response"
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
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	user.Profile.Author = createdUser.UserName
	_, err = h.services.Profile.Create(&user.Profile)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	expiresAt := time.Now().Add(120 * time.Second)
	accessToken, err := h.services.Session.Create(&model.Session{
		UserName:  createdUser.UserName,
		ExpiresAt: expiresAt,
		UserAgent: ctx.Request.UserAgent(),
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   accessToken,
		Expires: expiresAt,
	})

	response.NewDataResponse(ctx, http.StatusCreated, response.DataResponse{
		"message": "registration successfully completed",
	})
}

func (h *UserHandler) login(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	foundUser, err := h.services.User.Login(&user)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	expiresAt := time.Now().Add(120 * time.Second)
	accessToken, err := h.services.Session.Create(&model.Session{
		UserName:  foundUser.UserName,
		ExpiresAt: expiresAt,
		UserAgent: ctx.Request.UserAgent(),
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   accessToken,
		Expires: expiresAt,
	})

	response.NewDataResponse(ctx, http.StatusCreated, response.DataResponse{
		"message": "login successfully completed",
	})
}

func (h *UserHandler) logout(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	if err := h.services.AccessToken.RemoveByAccessToken(sessionToken.(string)); err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

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
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	response.NewDataResponse(ctx, http.StatusOK, response.DataResponse{
		"user": foundUser,
	})
}

func (h *UserHandler) getUsers(ctx *gin.Context) {
	limit := ctx.Query("limit")
	searchTerm := ctx.Query("search_term")

	users, err := h.services.User.GetUsers(searchTerm, limit)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	response.NewDataResponse(ctx, http.StatusOK, response.DataResponse{
		"users": users,
	})
}

func (h *UserHandler) sendSMSCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	session, fail := h.services.Session.GetByAccessToken(sessionToken.(string))
	if fail != nil {
		response.NewErrorResponse(ctx, fail.StatusCode, fail.Message)
		return
	}

	foundUser, err := h.services.User.GetByUserName(session.UserName)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	foundUser.SMSCode.GenerateConfirmationCode(foundUser.UserName)

	smsCode, err := h.services.SMSCode.SetSMSCode(foundUser.SMSCode, foundUser.UserName)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if err := h.services.SMSCode.SendSMSConfirmation(foundUser.PhoneNumber, smsCode.Code); err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	response.NewDataResponse(ctx, http.StatusOK, response.DataResponse{
		"message": "sms confirmation sended succefully",
	})
}

func (h *UserHandler) confirmSMSCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	foundSession, fail := h.services.Session.GetByAccessToken(sessionToken.(string))
	if fail != nil {
		ctx.JSON(fail.StatusCode, fail.Message)
		return
	}

	var smsCode model.SMSCode

	if err := ctx.BindJSON(&smsCode); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	confirmed, err := h.services.User.ConfirmSMSCode(foundSession.UserName, smsCode)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if !confirmed {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "failed to confirm sms code")
		return
	}

	response.NewDataResponse(ctx, http.StatusOK, response.DataResponse{
		"message": "sms code succefully confirmed",
	})
}

func (h *UserHandler) sendEmailCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")

	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	session, fail := h.services.Session.GetByAccessToken(sessionToken.(string))
	if fail != nil {
		response.NewErrorResponse(ctx, fail.StatusCode, fail.Message)
		return
	}

	foundUser, err := h.services.User.GetByUserName(session.UserName)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	foundUser.EmailCode.GenerateConfirmationCode(foundUser.UserName)

	emailCode, err := h.services.EmailCode.SetEmailCode(foundUser.EmailCode, foundUser.UserName)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if err := h.services.EmailCode.SendEmailCode(
		"../../templates/email_confirm.html",
		foundUser.Email,
		emailCode.Code,
	); err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	response.NewDataResponse(ctx, http.StatusOK, response.DataResponse{
		"message": "email confirmation sended succefully",
	})
}

func (h *UserHandler) confirmEmailCode(ctx *gin.Context) {
	sessionToken, ok := ctx.Get("session_token")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	session, fail := h.services.Session.GetByAccessToken(sessionToken.(string))
	if fail != nil {
		response.NewErrorResponse(ctx, fail.StatusCode, fail.Message)
		return
	}

	var emailCode model.EmailCode

	if err := ctx.BindJSON(&emailCode); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	confirmed, err := h.services.User.ConfirmEmailCode(session.UserName, emailCode)
	if err != nil {
		response.NewErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	if !confirmed {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "failed to confirm email code")
		return
	}

	response.NewDataResponse(ctx, http.StatusOK, response.DataResponse{
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
