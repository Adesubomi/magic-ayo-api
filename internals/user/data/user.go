package data

import (
	"errors"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	"gorm.io/gorm"
)

func (r UserRepo) FindUser(id string) (*authPkg.User, error) {
	user := &authPkg.User{}

	result := r.DbClient.
		Where(&authPkg.User{ID: id}).
		First(user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, logPkg.RecordNotFoundError
	} else if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r UserRepo) CreateAccount(id, password string) (*authPkg.User, error) {
	user := &authPkg.User{
		ID:       id,
		Password: password,
	}

	result := r.DbClient.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
