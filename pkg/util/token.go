package util

import (
	"crypto/rand"
	"fmt"
	"github.com/aidarkhanov/nanoid"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

func randomStringGenerator(length int8, charSet string, preventConfusion bool) string {
	characters := "~-!abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if charSet == "alpha" || charSet == "alphabet" || charSet == "alphabets" || charSet == "alphabetic" {
		characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	} else if charSet == "num" || charSet == "numeric" || charSet == "numbers" {
		characters = "1234567890"
	} else if charSet == "alphaNum" || charSet == "alphaNumeric" {
		characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	}

	if preventConfusion {
		characters = strings.Replace(characters, "o", "", -1)
		characters = strings.Replace(characters, "O", "", -1)
		characters = strings.Replace(characters, "0", "", -1)
		characters = strings.Replace(characters, "i", "", -1)
		characters = strings.Replace(characters, "l", "", -1)
		characters = strings.Replace(characters, "1", "", -1)
	}

	value, err := nanoid.Format(
		func(step int) ([]byte, error) {
			buffer := make([]byte, step)
			if _, err := rand.Read(buffer); err != nil {
				return nil, err
			}
			return buffer, nil
		},
		characters,
		int(length),
	)

	if err != nil {
		panic(err)
	}

	return value
}

// GenerateOTP
// channel: email | phone
func GenerateOTP(length int8, channel string) string {
	return randomStringGenerator(length, "num", true)
}

func GenerateModelID() string {
	return randomStringGenerator(12, "alphaNumeric", false)
}

func GenerateRandomID(length int8) string {
	return randomStringGenerator(length, "alphaNumeric", false)
}

func GetRedisStorageKey(identifier string, channel string, purpose string) string {
	// e.g. "forgot-password:email:daniel@example.com"
	return fmt.Sprintf("%v:%v:%v", purpose, channel, identifier)
}

func GetUserSessionStorageKey(userId string) string {
	return "session-token:" + userId
}

func GetTransactionStorageKey(identifier string) string {
	return fmt.Sprintf("transaction-%v", identifier)
}

func GetBearerTokenFromAuthorizationHeader(ctx *fiber.Ctx) string {
	authorization := ctx.GetReqHeaders()["Authorization"]
	parts := strings.Split(authorization, " ")
	if len(parts) != 2 {
		return ""
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func sessionRefStorageKey(identifier string) string {
	return GetRedisStorageKey(
		identifier,
		"wallet-transaction",
		"session-lock",
	)
}

func GenerateSessionRef(redisClient *redis.Client, identifier string) (string, error) {
	sessionRef := GenerateModelID()
	sessionRefKey := sessionRefStorageKey(identifier)

	redisClient.Set(
		sessionRefKey,
		sessionRef,
		3*time.Minute)

	return sessionRef, nil
}

func VerifySessionRef(redisClient *redis.Client, identifier string, sessionRef string) bool {
	sessionRefKey := sessionRefStorageKey(identifier)
	result := redisClient.Get(sessionRefKey)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
		}
		return false
	}

	value := result.Val()
	return sessionRef == value
}

func DeleteSessionRef(redisClient *redis.Client, identifier string) error {
	sessionRefKey := sessionRefStorageKey(identifier)
	return redisClient.Del(sessionRefKey).Err()
}
