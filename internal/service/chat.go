package service

import (
	proxyproto "chat-service/api/centrifugo"
	"chat-service/internal/biz"
	"context"
	"fmt"
)

type ChatService struct {
	uc *biz.ChatUseCase
	proxyproto.UnimplementedCentrifugoProxyServer
}

func NewChatService(uc *biz.ChatUseCase) *ChatService {
	return &ChatService{uc: uc}
}
func (s *ChatService) Publish(ctx context.Context, req *proxyproto.PublishRequest) (*proxyproto.PublishResponse, error) {
	fmt.Print(req)
	// parts := strings.Split(strings.TrimPrefix(req.Channel, "dialog#"), ",")
	// if len(parts) != 2 {
	// 	return nil, errors.New("ERROR_PUBLISH")
	// }
	// userID := parts[0]
	// recieverID := parts[1]
	// ch, err := s.uc.GetChat(ctx, userID, recieverID)
	// if err != nil {
	// 	return nil, errors.New("DB_ERROR")
	// }
	// s.uc.CreateMessage(ctx, &biz.Message{
	// 	ChatID: ch.ID,
	// 	Message: req.con,
	// })
	return &proxyproto.PublishResponse{}, nil
}
