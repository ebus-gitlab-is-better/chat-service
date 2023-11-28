package biz

import (
	"context"
	"time"
)

type Message struct {
	ID         uint `gorm:"primarykey"`
	ChatID     uint
	Message    string
	SenderID   string
	SenderName string
	SentAt     time.Time
}

type Chat struct {
	ID         uint
	UserID     string
	ReceiverID string
	Message    []Message
}

type ChatRepo interface {
	CreateMessage(context.Context, *Message) error
	GetChatHistory(context.Context, uint, uint) ([]*Message, error)
	CreateChat(context.Context, string, string) error
	GetChat(context.Context, string, string) (*Chat, error)
	DeleteChat(context.Context, uint) error
	GetListChat(context.Context, string) ([]*Chat, error)
}

type ChatUseCase struct {
	repo ChatRepo
}

func NewChatUseCase(repo ChatRepo) *ChatUseCase {
	return &ChatUseCase{
		repo: repo,
	}
}

func (uc *ChatUseCase) CreateMessage(ctx context.Context, msg *Message) error {
	return uc.repo.CreateMessage(ctx, msg)
}

func (uc *ChatUseCase) GetChatHistory(ctx context.Context, chatId, offset uint) ([]*Message, error) {
	return uc.repo.GetChatHistory(ctx, chatId, offset)
}

func (uc *ChatUseCase) CreateChat(ctx context.Context, userId, receiverId string) error {
	return uc.repo.CreateChat(ctx, userId, receiverId)
}

func (uc *ChatUseCase) GetListChat(ctx context.Context, userID string) ([]*Chat, error) {
	return uc.repo.GetListChat(ctx, userID)
}

func (uc *ChatUseCase) DeleteChat(ctx context.Context, name string, userID uint) error {
	return uc.repo.DeleteChat(ctx, userID)
}

func (uc *ChatUseCase) GetChat(ctx context.Context, userID string, receiverID string) (*Chat, error) {
	return uc.repo.GetChat(ctx, userID, receiverID)
}
