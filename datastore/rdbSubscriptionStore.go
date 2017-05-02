package datastore

import (
	"log"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateSubscriptionStore() {
	tableMap := dbMap.AddTableWithName(models.Subscription{}, TABLE_NAME_SUBSCRIPTION)
	tableMap.SetUniqueTogether("room_id", "user_id", "platform")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertSubscription(subscription *models.Subscription) StoreResult {
	result := StoreResult{}
	if err := dbMap.Insert(subscription); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating subscription item.", err)
	}
	result.Data = subscription
	return result
}

func RdbSelectSubscription(roomId, userId string, platform int) StoreResult {
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;")
	params := map[string]interface{}{
		"roomId":   roomId,
		"userId":   userId,
		"platform": platform,
	}
	if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription item.", err)
	}
	if len(subscriptions) == 1 {
		result.Data = subscriptions[0]
	}
	return result
}

func RdbSelectSubscriptionsByRoomId(roomId string) StoreResult {
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbSelectSubscriptionsByUserId(userId string) StoreResult {
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{
		"userId": userId,
	}
	if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbSelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND platform=:platform AND deleted=0;")
	params := map[string]interface{}{
		"roomId":   roomId,
		"platform": platform,
	}
	if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbSelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND platform=:platform AND deleted=0;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbDeleteSubscription(subscription *models.Subscription) StoreResult {
	result := StoreResult{}
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"roomId":   subscription.RoomId,
		"userId":   subscription.UserId,
		"platform": subscription.Platform,
	}
	_, err := dbMap.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting subscription item.", err)
	}
	return result
}
