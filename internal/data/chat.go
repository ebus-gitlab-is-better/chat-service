package data

import (
	"chat-service/internal/biz"
	"context"
	"encoding/json"
	"errors"
	"time"

	gocent "github.com/centrifugal/gocent/v3"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Chat struct {
	ID         uint `gorm:"primarykey"`
	UserID     string
	ReceiverID string
	Message    []Message
}

type Message struct {
	ID         uint `gorm:"primarykey"`
	ChatID     uint
	Message    string
	SenderID   string
	SenderName string
	SentAt     time.Time
}

func (m Message) modelToResponse() *biz.Message {
	return &biz.Message{
		ID: m.ID,
	}
}

type chatRepo struct {
	data   *Data
	logger *log.Helper
	cent   *gocent.Client
}

func NewChatRepo(data *Data, logger log.Logger, cent *gocent.Client) biz.ChatRepo {
	return &chatRepo{data: data, logger: log.NewHelper(logger), cent: cent}
}

// Create implements biz.ChatRepo.
func (r *chatRepo) CreateMessage(_ context.Context, msg *biz.Message, channel string) error {
	msgDB := Message{}
	msgDB.SenderID = msg.SenderID
	msgDB.SentAt = time.Now()
	msgDB.SenderName = msg.SenderName
	msgDB.ChatID = msg.ChatID
	msgDB.Message = msg.Message
	if err := r.data.db.Create(&msgDB).Error; err != nil {
		return err
	}
	data, _ := json.Marshal(msg)
	_, err := r.cent.Publish(context.TODO(), channel, data)
	if err != nil {
		return err
	}
	return nil
}

// GetHistory implements biz.ChatRepo.
func (r *chatRepo) GetChatHistory(_ context.Context, chatId uint, offset uint) ([]*biz.Message, error) {
	var messagesDB []Message
	if err := r.data.db.
		Where(&Message{ChatID: chatId}).
		Scopes(paginate(uint32(offset), uint32(100))).
		Find(&messagesDB).Error; err != nil {
		return nil, err
	}
	messages := make([]*biz.Message, 0)
	for _, msg := range messagesDB {
		messages = append(messages, msg.modelToResponse())
	}
	return messages, nil
}

// CreateChannel implements biz.ChatRepo.
func (r *chatRepo) CreateChat(_ context.Context, userID string, receiverID string) error {
	var ch Chat
	res := r.data.db.Where(&Chat{UserID: userID, ReceiverID: receiverID}).Find(&ch)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 0 {
		return errors.New("CHANNEL_IS_EXISTS")
	}
	ch.UserID = userID
	ch.ReceiverID = receiverID
	if err := r.data.db.Create(&ch).Error; err != nil {
		return err
	}
	ch.UserID = receiverID
	ch.ReceiverID = userID
	if err := r.data.db.Create(&ch).Error; err != nil {
		return err
	}
	return nil
}

// GetListChannel implements biz.ChatRepo.
func (r *chatRepo) GetListChat(_ context.Context, userID string) ([]*biz.Chat, error) {
	var channelsDB []Chat
	if err := r.data.db.Preload("Message", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.sent_at DESC").Limit(1)
	}).
		Where(&Chat{UserID: userID}).
		Find(&channelsDB).Error; err != nil {
		return nil, err
	}
	channels := make([]*biz.Chat, 0)
	for _, ch := range channelsDB {
		chBiz := &biz.Chat{
			ID:         ch.ID,
			UserID:     ch.UserID,
			ReceiverID: ch.ReceiverID,
		}
		if len(ch.Message) > 0 {
			chBiz.Message = []biz.Message{*ch.Message[0].modelToResponse()}
		}
		channels = append(channels, chBiz)
	}
	return channels, nil
}

// DeleteChannel implements biz.ChatRepo.
func (r *chatRepo) DeleteChat(_ context.Context, id uint) error {
	if err := r.data.db.Delete(&Chat{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

// GetChat implements biz.ChatRepo.
func (r *chatRepo) GetChat(_ context.Context, userID string, receieverID string) (*biz.Chat, error) {
	var channelsDB Chat
	if err := r.data.db.
		Where(&Chat{UserID: userID, ReceiverID: receieverID}).
		Find(&channelsDB).Error; err != nil {
		return nil, err
	}
	return &biz.Chat{
		UserID:     channelsDB.UserID,
		ReceiverID: channelsDB.ReceiverID,
	}, nil
}

// GetChatById implements biz.ChatRepo.
func (r *chatRepo) GetChatById(_ context.Context, id uint) (*biz.Chat, error) {
	var channelsDB Chat
	if err := r.data.db.
		Where(&Chat{ID: id}).
		Find(&channelsDB).Error; err != nil {
		return nil, err
	}
	return &biz.Chat{
		UserID:     channelsDB.UserID,
		ReceiverID: channelsDB.ReceiverID,
	}, nil
}
