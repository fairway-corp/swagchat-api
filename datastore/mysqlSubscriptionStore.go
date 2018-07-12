package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.database)
}

func (p *mysqlProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	return rdbInsertSubscription(p.database, room)
}

func (p *mysqlProvider) SelectSubscription(roomID, userID string, platform int32) (*model.Subscription, error) {
	return rdbSelectSubscription(p.database, roomID, userID, platform)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.database, roomID)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int32) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.database, userID, platform)
}

func (p *mysqlProvider) DeleteSubscription(subscription *model.Subscription) error {
	return rdbDeleteSubscription(p.database, subscription)
}
