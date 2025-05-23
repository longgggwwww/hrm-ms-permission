package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	userPb "github.com/huynhthanhthao/hrm_user_service/generated"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/perm"
)

type PermHandler struct {
	Client     *ent.Client
	UserClient *userPb.UserServiceClient
}

func NewPermHandler(client *ent.Client, userClient *userPb.UserServiceClient) *PermHandler {
	return &PermHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *PermHandler) RegisterRoutes(r *gin.Engine) {
	perms := r.Group("/perms")
	{
		perms.GET("", h.List)
		perms.GET(":id", h.Get)
	}
}

func (h *PermHandler) List(c *gin.Context) {
	log.Println("List all perms")
	perms, err := h.Client.Perm.Query().
		WithGroup().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, perms)
}

func (h *PermHandler) Get(c *gin.Context) {
	log.Println("Get perm by ID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid UUID format for role ID",
		})
		return
	}

	perm, err := h.Client.Perm.Query().
		Where(perm.IDEQ(id)).
		WithRoles().
		Only(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Perm not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, perm)
}
