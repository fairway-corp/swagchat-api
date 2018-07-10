package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
	"google.golang.org/grpc"
)

func webhookRoom(ctx context.Context, room *model.Room) {
	webhooks, err := datastore.Provider(ctx).SelectWebhooks(model.WebhookEventTypeRoom, datastore.WithRoomID(datastore.RoomIDAll))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(webhooks) == 0 {
		return
	}

	pbRoom := &protobuf.Room{
		Workspace: ctx.Value(utils.CtxWorkspace).(string),
		RoomId:    room.RoomID,
	}

	for _, webhook := range webhooks {
		pbRoom.WebhookToken = webhook.Token

		switch webhook.Protocol {
		case model.WebhookProtocolHTTP:
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(pbRoom)

			resp, err := http.Post(
				webhook.Endpoint,
				"application/json",
				buf,
			)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookRoom] Post failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookRoom] Response body read failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logger.Error(fmt.Sprintf("[HTTP][WebhookRoom] Status code is not 200. Endpoint=[%s] StatusCode[%d].", webhook.Endpoint, resp.StatusCode))
				continue
			}
		case model.WebhookProtocolGRPC:
			conn, err := grpc.Dial(webhook.Endpoint, grpc.WithInsecure())
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookRoom] Connect failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			defer conn.Close()

			c := protobuf.NewChatOutgoingClient(conn)
			_, err = c.PostWebhookRoom(context.Background(), pbRoom)
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookRoom] Response body read failure. GRPC Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
		}
	}
}

func webhookMessage(ctx context.Context, message *model.Message, user *model.User) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(message.RoomID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	webhooks, err := datastore.Provider(ctx).SelectWebhooks(model.WebhookEventTypeMessage, datastore.WithRoomID(datastore.RoomIDAll))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(webhooks) == 0 {
		return
	}

	// Only support text message
	if message.Type != model.MessageTypeText {
		return
	}

	var p model.PayloadText
	json.Unmarshal(message.Payload, &p)

	pbMessage := &protobuf.Message{
		Workspace: ctx.Value(utils.CtxWorkspace).(string),
		UserIds:   userIDs,
		RoomId:    message.RoomID,
		UserId:    message.UserID,
		Type:      message.Type,
		Payload: &protobuf.MessagePayload{
			Text: p.Text,
		},
	}

	for _, webhook := range webhooks {
		matchRole := false
		for _, v := range user.Roles {
			if v == webhook.RoleID {
				matchRole = true
			}
		}

		if !matchRole {
			continue
		}

		pbMessage.WebhookToken = webhook.Token

		switch webhook.Protocol {
		case model.WebhookProtocolHTTP:
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(pbMessage)

			resp, err := http.Post(
				webhook.Endpoint,
				"application/json",
				buf,
			)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookMessage] Post failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookMessage] Response body read failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logger.Error(fmt.Sprintf("[HTTP][WebhookMessage] Status code is not 200. Endpoint=[%s] StatusCode[%d]", webhook.Endpoint, resp.StatusCode))
				continue
			}
		case model.WebhookProtocolGRPC:
			conn, err := grpc.Dial(webhook.Endpoint, grpc.WithInsecure())
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookMessage] Connect failure. Endpoint=[%s]. %v", webhook.Endpoint, err))
				continue
			}
			defer conn.Close()

			c := protobuf.NewChatOutgoingClient(conn)
			_, err = c.PostWebhookMessage(context.Background(), pbMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookMessage] Response body read failure. GRPC Endpoint=[%s]. %v", webhook.Endpoint, err))
				continue
			}
		}
	}
}