package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

// CreateRoomUsers creates room users
func CreateRoomUsers(ctx context.Context, req *model.CreateRoomUsersRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start  CreateRoomUsers. Request[%#v]", req))

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to create room users."
		return errRes
	}

	req.Room = room

	userIDs, errRes := getExistUserIDs(ctx, req.UserIDs)
	if errRes != nil {
		errRes.Message = "Failed to create room users."
		return errRes
	}
	req.UserIDs = userIDs

	errRes = req.Validate()
	if errRes != nil {
		return errRes
	}

	// if room.NotificationTopicID == "" {
	// 	notificationTopicID, pd := createTopicOld(room.RoomID)
	// 	if pd != nil {
	// 		return pd
	// 	}

	// 	room.NotificationTopicID = notificationTopicID
	// 	room.Modified = time.Now().Unix()
	// 	_, err := datastore.Provider(ctx).UpdateRoom(room)
	// 	if err != nil {
	// 		pd := &model.ProblemDetail{
	// 			Message: "Failed to create room users.",
	// 			Status:  http.StatusInternalServerError,
	// 			Error:   err,
	// 		}
	// 		return pd
	// 	}
	// }

	roomUsers := req.GenerateRoomUsers()
	err := datastore.Provider(ctx).InsertRoomUsers(roomUsers)
	if err != nil {
		return model.NewErrorResponse("Failed to create room users.", nil, http.StatusInternalServerError, err)
	}

	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, req.RoomID)

	logger.Info("Finish CreateRoomUsers.")
	return nil
}

func GetUserIDsOfRoomUser(ctx context.Context, req *model.GetUserIdsOfRoomUserRequest) (*scpb.UserIds, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  SelectUserIDsOfRoomUser. Request[%#v]", req))

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
		req.RoomID,
		datastore.SelectUserIDsOfRoomUserOptionWithRoleIDs(req.RoleIDs),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get userIds.", nil, http.StatusInternalServerError, err)
	}

	logger.Info("Finish SelectUserIDsOfRoomUser.")
	return &scpb.UserIds{
		UserIDs: userIDs,
	}, nil
}

// UpdateRoomUser updates room user
func UpdateRoomUser(ctx context.Context, req *model.UpdateRoomUserRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start UpdateRoomUser. Request[%#v]", req))

	ru, errRes := confirmRoomUserExist(ctx, req.RoomID, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to update room user."
		return errRes
	}

	ru.UpdateRoomUser(req)

	err := datastore.Provider(ctx).UpdateRoomUser(ru)
	if err != nil {
		return model.NewErrorResponse("Failed to update room user.", nil, http.StatusInternalServerError, err)
	}

	// var p json.RawMessage
	// err = json.Unmarshal([]byte("{}"), &p)
	// m := &model.Message{
	// 	RoomID:    roomUser.RoomID,
	// 	UserID:    roomUser.UserID,
	// 	Type:      model.MessageTypeUpdateRoomUser,
	// 	EventName: "message",
	// 	Payload:   p,
	// }
	// rtmPublish(ctx, m)

	logger.Info("Finish UpdateRoomUser.")
	return nil
}

// DeleteRoomUsers deletes room users
func DeleteRoomUsers(ctx context.Context, req *model.DeleteRoomUsersRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start DeleteRoomUsers. Request[%#v]", req))

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to delete room users."
		return errRes
	}

	req.Room = room

	err := datastore.Provider(ctx).DeleteRoomUsers(req.RoomID, req.UserIDs)
	if err != nil {
		return model.NewErrorResponse("Failed to delete room users.", nil, http.StatusInternalServerError, err)
	}

	go func() {
		rus, err := datastore.Provider(ctx).SelectRoomUsers(
			datastore.SelectRoomUsersOptionWithRoomID(req.RoomID),
			datastore.SelectRoomUsersOptionWithUserIDs(req.UserIDs),
		)
		if err != nil {
			logger.Error(err.Error())
		}

		unsubscribeByRoomUsers(ctx, rus)
	}()

	logger.Info(fmt.Sprintf("Finish DeleteRoomUsers."))
	return nil
}
