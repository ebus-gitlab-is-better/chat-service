package route

import (
	"chat-service/internal/biz"
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

type ChatRoute struct {
	uc *biz.ChatUseCase
}

func NewChatRoute(uc *biz.ChatUseCase) *ChatRoute {
	return &ChatRoute{uc: uc}
}

func (r *ChatRoute) Register(router *gin.RouterGroup) {
	router.GET("/:id/history", r.GetHistory)
	router.POST("/", r.CreateChat)
	router.GET("/", r.GetChats)
	// router.DELETE("/chats/:id", )
}

type CreateChatDTO struct {
	UserID     string `json:"user_id"`
	ReceiverID string `json:"receiver_id"`
}

// @Summary	Create chat
// @Accept		json
// @Produce	json
// @Tags		chat
// @Param		dto	body	route.CreateChatDTO	true	"dto"
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/chats/ [post]
func (r *ChatRoute) CreateChat(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := CreateChatDTO{}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	err = r.uc.CreateChat(context.TODO(), dto.UserID, dto.ReceiverID)
	if err != nil {
		c.JSON(500, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
}

type MessageDTO struct {
	Messages []*biz.Message `json:"messages"`
}

type ChatDTO struct {
	biz.Chat
}

// @Summary	Get chats
// @Accept		json
// @Produce	json
// @Tags		chat
// @Param		id	path	int	true	"Chat ID"	Format(uint64)
//
//	@Success	200	{object}	route.ChatDTO
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/chats/ [get]
func (r *ChatRoute) GetChats(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.Status(400)
		return
	}
	userKeycloak, ok := user.(*gocloak.UserInfo)
	if !ok {
		c.Status(400)
		return
	}
	chats, err := r.uc.GetListChat(context.TODO(), *userKeycloak.Sub)
	if err != nil {
		c.JSON(500, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, &ChatsDTO{
		Chats: chats,
	})
}

type ChatsDTO struct {
	Chats []*biz.Chat `json:"chats"`
}

// @Summary	Get chat history
// @Accept		json
// @Produce	json
// @Tags		chat
//
// @Param		id	path	int	true	"Chat ID"	Format(uint64)
//
//	@Param		page			query	int		false	"offset"	Format(uint64)
//
//	@Success	200	{object}	route.MessageDTO
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/chats/ [get]
func (r *ChatRoute) GetHistory(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}

	offset := c.DefaultQuery("offset", "1")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": "offset not found",
		})
		return
	}

	messages, err := r.uc.GetChatHistory(context.TODO(), uint(idUint), uint(offsetInt))
	if err != nil {
		c.JSON(500, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, MessageDTO{
		Messages: messages,
	})
}
