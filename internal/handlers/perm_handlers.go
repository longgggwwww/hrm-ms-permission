package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userPb "github.com/huynhthanhthao/hrm_user_service/generated"
	"github.com/longgggwwww/hrm-ms-permission/ent"
)

type PermHandler struct {
	Client     *ent.Client
	UserClient *userPb.UserServiceClient
}

func NewPermHandler(c *ent.Client, user *userPb.UserServiceClient) *PermHandler {
	return &PermHandler{
		Client:     c,
		UserClient: user,
	}
}

func (h *PermHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.List)
}

func (h *PermHandler) List(c *gin.Context) {
	out, err := h.Client.Perm.Query().
		WithGroup().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}
