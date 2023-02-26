package auth

import (
	"encoding/json"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	utilPkg "github.com/Adesubomi/magic-ayo-api/pkg/util"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
)

type Registry struct {
	RedisClient *redis.Client
}

func (r Registry) CreateUserSession(user *User) (*UserSession, error) {
	return nil, nil
}

// GetUserAuthSession get the current user auth session
// from the bearer token in the authorization header
func (r Registry) getUserAuthSession(tokenString string) (*UserSession, error) {

	userSession := &UserSession{}
	jwtClaims := &JwtTokenClaim{}
	_, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	if err != nil {
		logPkg.ReportError(err)
		return nil, err
	}

	storageKey := utilPkg.GetUserSessionStorageKey(jwtClaims.ID)
	userSessionResult, err := r.RedisClient.Get(storageKey).Result()
	if err != nil {
		logPkg.ReportError(err)
		return nil, err
	}

	err = json.Unmarshal([]byte(userSessionResult), &userSession)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}
