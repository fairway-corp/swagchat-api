package datastore

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

func rdbCreateRoomUserStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.RoomUser{}, tableNameRoomUser)
	tableMap.SetUniqueTogether("room_id", "user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating roomUser. %v.", err))
		return
	}
}

func rdbInsertRoomUsers(db string, roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		err := errors.Wrap(err, "An error occurred while recreating roomUser")
		logger.Error(err.Error())
		return err
	}

	opt := insertRoomUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.beforeClean {
		query := fmt.Sprintf("DELETE FROM %s WHERE room_id=:roomId;", tableNameRoomUser)
		params := map[string]interface{}{"roomId": roomUsers[0].RoomID}
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			err := errors.Wrap(err, "An error occurred while recreating roomUser")
			logger.Error(err.Error())
			return err
		}
	}

	for _, roomUser := range roomUsers {
		err = trans.Insert(roomUser)
		if err != nil {
			trans.Rollback()
			err := errors.Wrap(err, "An error occurred while recreating roomUser")
			logger.Error(err.Error())
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err := errors.Wrap(err, "An error occurred while recreating roomUser")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbSelectRoomUsers(db string, opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	replica := RdbStore(db).replica()

	opt := selectRoomUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var roomUsers []*model.RoomUser
	var userIDsQuery string
	var userIDsParams map[string]interface{}
	var roomIDParams map[string]interface{}

	if opt.roomID != "" {
		roomIDParams = map[string]interface{}{"roomId": opt.roomID}
	}
	if opt.userIDs != nil {
		userIDsQuery, userIDsParams = utils.MakePrepareForInExpression(opt.userIDs)
	}
	params := make(map[string]interface{}, len(userIDsParams)+len(roomIDParams))
	params = utils.MergeMap(params, userIDsParams, roomIDParams)

	query := fmt.Sprintf("SELECT * FROM %s WHERE", tableNameRoomUser)
	if opt.roomID != "" {
		query = fmt.Sprintf("%s room_id=:roomId", query)
	}
	if opt.roomID != "" && opt.userIDs != nil {
		query = fmt.Sprintf("%s AND ", query)
	}
	if opt.userIDs != nil {
		query = fmt.Sprintf("%s user_id IN (%s)", query, userIDsQuery)
	}

	var err error
	if params == nil {
		_, err = replica.Select(&roomUsers, query)
	} else {
		_, err = replica.Select(&roomUsers, query, params)
	}
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting room users")
		logger.Error(err.Error())
		return nil, err
	}

	return roomUsers, nil
}

func rdbSelectRoomUser(db, roomID, userID string) (*model.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*model.RoomUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND user_id=:userId;", tableNameRoomUser)
	params := map[string]interface{}{
		"roomId": roomID,
		"userId": userID,
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting roomUser")
		logger.Error(err.Error())
		return nil, err
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectRoomUserOfOneOnOne(db, myUserID, opponentUserID string) (*model.RoomUser, error) {
	replica := RdbStore(db).replica()

	var roomUsers []*model.RoomUser
	query := fmt.Sprintf(`SELECT * FROM %s
WHERE room_id IN (
	SELECT room_id FROM %s WHERE type=:type AND user_id=:myUserId
) AND user_id=:opponentUserId;`, tableNameRoomUser, tableNameRoom)
	params := map[string]interface{}{
		"type":           scpb.RoomType_OneOnOne,
		"myUserId":       myUserID,
		"opponentUserId": opponentUserID,
	}
	_, err := replica.Select(&roomUsers, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting roomUser for OneOnOne")
		logger.Error(err.Error())
		return nil, err
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectUserIDsOfRoomUser(db string, roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	replica := RdbStore(db).replica()

	opt := selectUserIDsOfRoomUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var userIDs []string

	var query string
	var params map[string]interface{}
	if opt.roleIDs == nil {
		query = fmt.Sprintf("SELECT ru.user_id FROM %s AS ru LEFT JOIN %s AS u ON ru.user_id = u.user_id WHERE ru.room_id=:roomId;", tableNameRoomUser, tableNameUser)
		params = map[string]interface{}{
			"roomId": roomID,
		}
	} else {
		roleIDsQuery, pms := utils.MakePrepareForInExpression(opt.roleIDs)
		params = pms
		query = fmt.Sprintf("SELECT ru.user_id FROM %s AS ru LEFT JOIN %s AS ur ON ru.user_id = ur.user_id WHERE ru.room_id=:roomId AND ur.role_id IN (%s) GROUP BY ru.user_id;", tableNameRoomUser, tableNameUserRole, roleIDsQuery)
		params["roomId"] = roomID
	}

	_, err := replica.Select(&userIDs, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return userIDs, nil
}

func rdbUpdateRoomUser(db string, ru *model.RoomUser) error {
	master := RdbStore(db).master()

	query := fmt.Sprintf("UPDATE %s SET unread_count=:unreadCount WHERE room_id=:roomId AND user_id=:userId;", tableNameRoomUser)
	params := map[string]interface{}{
		"roomId":      ru.RoomID,
		"userId":      ru.UserID,
		"unreadCount": ru.UnreadCount,
	}
	_, err := master.Exec(query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while updating room user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbDeleteRoomUsers(db, roomID string, userIDs []string) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		err := errors.Wrap(err, "An error occurred while deleting room users")
		logger.Error(err.Error())
		return err
	}

	var query string
	var params map[string]interface{}
	if userIDs == nil {
		query = fmt.Sprintf("DELETE FROM %s WHERE room_id=:roomId;", tableNameRoomUser)
		params = map[string]interface{}{"roomId": roomID}
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}

		query = fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE room_id=:roomId;", tableNameSubscription)
		params = map[string]interface{}{
			"roomId":  roomID,
			"deleted": time.Now().Unix(),
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}
	} else {
		var userIdsQuery string
		userIdsQuery, params = utils.MakePrepareForInExpression(userIDs)
		query = fmt.Sprintf("DELETE FROM %s WHERE room_id=:roomId AND user_id IN (%s);", tableNameRoomUser, userIdsQuery)
		params["roomId"] = roomID
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}

		query = fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE room_id=:roomId AND user_id IN (%s);", tableNameSubscription, userIdsQuery)
		params["deleted"] = time.Now().Unix()
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err := errors.Wrap(err, "An error occurred while deleting room users")
		logger.Error(err.Error())
		return err
	}

	return nil
}

// func rdbUpdateRoomUser(db string, roomUser *model.RoomUser) (*model.RoomUser, error) {
// 	master := RdbStore(db).master()
// 	trans, err := master.Begin()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "An error occurred while transaction beginning")
// 	}

// 	updateQuery := ""
// 	params := map[string]interface{}{
// 		"roomId": roomUser.RoomID,
// 		"userId": roomUser.UserID,
// 	}
// 	if roomUser.UnreadCount != nil {
// 		params["unreadCount"] = roomUser.UnreadCount
// 		updateQuery = "unread_count=:unreadCount"
// 	}
// 	// if roomUser.MetaData != nil {
// 	// 	params["metaData"] = roomUser.MetaData
// 	// 	if updateQuery == "" {
// 	// 		updateQuery = "meta_data=:metaData"
// 	// 	} else {
// 	// 		updateQuery = utils.AppendStrings(updateQuery, ",", "meta_data=:metaData")
// 	// 	}
// 	// }
// 	if updateQuery != "" {
// 		query := utils.AppendStrings("UPDATE ", tableNameRoomUser, " SET "+updateQuery+" WHERE room_id=:roomId AND user_id=:userId;")
// 		_, err = trans.Exec(query, params)
// 		if err != nil {
// 			trans.Rollback()
// 			return nil, errors.Wrap(err, "An error occurred while updating room's users")
// 		}

// 		if roomUser.UnreadCount != nil {
// 			query = utils.AppendStrings("UPDATE ", tableNameUser,
// 				" SET unread_count=(SELECT SUM(unread_count) FROM ", tableNameRoomUser,
// 				" WHERE user_id=:userId1) WHERE user_id=:userId2;")
// 			params = map[string]interface{}{
// 				"userId1": roomUser.UserID,
// 				"userId2": roomUser.UserID,
// 			}
// 			_, err = trans.Exec(query, params)
// 			if err != nil {
// 				trans.Rollback()
// 				return nil, errors.Wrap(err, "An error occurred while updating user unread count")
// 			}
// 		}
// 	}

// 	err = trans.Commit()
// 	if err != nil {
// 		trans.Rollback()
// 		return nil, errors.New("An error occurred while commit updating room's user")
// 	}

// 	return roomUser, nil
// }
