package handlers

import (
	"net/http"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	services *service.Service
}

func NewPostHandler(services *service.Service) *PostHandler {
	return &PostHandler{services}
}

func (h *PostHandler) createPost(ctx *gin.Context) {
	var post *model.Post
	if err := ctx.BindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, "failed to bind JSON: "+err.Error())
		return
	}

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

	createdPost, err := h.services.Post.Create(post, foundUser.UserName)
	if err != nil {
		newErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	newDataResponse(ctx, http.StatusCreated, dataResponse{
		"message": "post created successfully",
		"post":    createdPost,
	})
}
