package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *sqliteProvider) InsertUserRoles(urs []*model.UserRole) error {
	return rdbInsertUserRoles(p.database, urs)
}

func (p *sqliteProvider) SelectUserRole(userID string, roleID int32) (*model.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *sqliteProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *sqliteProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
	return rdbDeleteUserRoles(p.database, opts...)
}
