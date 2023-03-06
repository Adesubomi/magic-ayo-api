package auth

import (
	"encoding/json"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	utilPkg "github.com/Adesubomi/magic-ayo-api/pkg/util"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"time"
)

type Registry struct {
	RedisClient *redis.Client
}

func (r Registry) CreateUserSession(user *User) (*UserSession, error) {
	userSession := new(UserSession)

	if user == nil || user.ID == "" || user.Password == "" {
		return userSession, logPkg.UnknownError
	}

	expiresAtDuration := time.Hour * 48
	expiresAtTime := time.Now().Add(expiresAtDuration)
	claims := JwtTokenClaim{
		ID:        user.ID,
		ExpiresAt: expiresAtTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.ToMap())

	// Sign and get the complete encoded token as a string
	// using the user password as secret. This is so
	// that the token will be invalidated if
	// the user's password changes
	tokenString, err := token.SignedString([]byte(""))

	if err != nil {
		return userSession, err
	}

	userSession = &UserSession{
		Token:     tokenString,
		User:      *user,
		ExpiresAt: expiresAtTime.Unix(),
	}

	sessionClaimString, err := json.Marshal(userSession)
	if err != nil {
		return nil, err
	}

	//modify
	storageKey := utilPkg.GetUserSessionStorageKey(user.ID)
	r.RedisClient.Set(storageKey, sessionClaimString, expiresAtDuration)
	return userSession, nil
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
