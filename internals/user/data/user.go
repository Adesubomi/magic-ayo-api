package data

import (
	"errors"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	utilPkg "github.com/Adesubomi/magic-ayo-api/pkg/util"
	"github.com/go-sql-driver/mysql"
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
	passwordBcrypt, _ := utilPkg.BcryptHash(password)
	user := &authPkg.User{
		ID:       id,
		Password: passwordBcrypt,
	}

	result := r.DbClient.Create(user)
	var mysqlErr *mysql.MySQLError
	if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
		return nil, logPkg.DuplicateRecordError
	} else if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
