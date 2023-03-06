package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Adesubomi/magic-ayo-api/internals/aayo/entity"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	"github.com/go-redis/redis"
)

func (r Repo) StartGame(userID string) (*entity.GamePlay, error) {
	gamePlay := entity.NewGamePlay(userID)

	storageKey := r.GetGameStorageKey(userID)
	fmt.Println(storageKey)
	result := r.RedisClient.Get(storageKey)
	if result.Err() != redis.Nil {
		existingGamePlay, err := entity.FromString(result.Val())

		if err == nil {
			gamePlay = existingGamePlay
		}
	}

	gamePlayRaw, err := json.Marshal(gamePlay)
	if err != nil {
		logPkg.ReportError(err)
	}

	r.RedisClient.Set(
		storageKey,
		gamePlayRaw,
		0)

	return gamePlay, nil
}

func (r Repo) GetActiveGame(userID string) (*entity.GamePlay, error) {
	storageKey := r.GetGameStorageKey(userID)
	result := r.RedisClient.Get(storageKey)
	if result.Err() != nil && result.Err() == redis.Nil {
		return nil, logPkg.RecordNotFoundError
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	existingGamePlay, err := entity.FromString(result.Val())
	if err != nil {
		return nil, err
	}

	return existingGamePlay, nil
}

func (r Repo) MakeMove(game *entity.GamePlay, pot int) error {
	err := game.AddMove(pot)

	fmt.Println("<< . . . >> did you add ?", game.Moves)
	if err != nil {
		return err
	}

	return r.updateGameToRedis(game)
}

func (r Repo) AbortGame(userID string) error {
	activeGame, err := r.GetActiveGame(userID)
	if err != nil && errors.Is(err, logPkg.RecordNotFoundError) {
		return err
	} else if activeGame == nil {
		return nil
	}

	err = r.storeGameToHistory(activeGame)
	if err != nil {
		return err
	}

	storageKey := r.GetGameStorageKey(userID)
	return r.RedisClient.Del(storageKey).Err()
}

func (r Repo) updateGameToRedis(gamePlay *entity.GamePlay) error {
	gamePlayRaw, err := json.Marshal(gamePlay)

	if err != nil {
		logPkg.ReportError(err)
	}

	storageKey := r.GetGameStorageKey(gamePlay.UserId)
	fmt.Println("<< . . . >> storage key ?", gamePlay.UserId)
	return r.RedisClient.Set(
		storageKey,
		gamePlayRaw,
		0).Err()
}

func (r Repo) storeGameToHistory(gamePlay *entity.GamePlay) error {
	if gamePlay != nil {
		result := r.DbClient.Create(gamePlay)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
