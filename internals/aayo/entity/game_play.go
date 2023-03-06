package entity

import (
	"encoding/json"
	"fmt"
	utilPkg "github.com/Adesubomi/magic-ayo-api/pkg/util"
	"strconv"
	"strings"
	"time"
)

type GameMovesList []int

func (gML GameMovesList) ToString() string {
	movesAsStringList := make([]string, 0)
	for _, moveAsInt := range gML {
		if moveAsInt != 0 {
			movesAsStringList = append(
				movesAsStringList,
				fmt.Sprintf("%v", moveAsInt),
			)
		}
	}
	return strings.Join(movesAsStringList, ",")
}

type GamePlay struct {
	ID        string    `json:"id"`
	UserId    string    `json:"userId"`
	Moves     string    `json:"moves"` // csv
	CreatedAt time.Time `gorm:"created_at,autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"updated_at,autoUpdateTime" json:"updatedAt"`
}

func NewGamePlay(identifier string) *GamePlay {
	return &GamePlay{
		ID:     utilPkg.GenerateModelID(),
		UserId: identifier,
	}
}

func FromString(i string) (*GamePlay, error) {
	gamePlay := &GamePlay{}
	return gamePlay, json.Unmarshal([]byte(i), gamePlay)
}

func (g *GamePlay) AddMove(move int) error {
	gameState, err := g.GetLatestGameState()
	if err != nil {
		return err
	}

	err = gameState.MakeMove(move, false)
	if err != nil {
		return err
	}

	movesList := append(g.GetMovesList(), move)
	g.Moves = movesList.ToString()
	fmt.Println("<...........> moves changed to . . .")
	fmt.Println(movesList.ToString())
	return nil
}

func (g *GamePlay) GetMovesList() GameMovesList {
	var moves GameMovesList
	movesStr := strings.Split(g.Moves, ",")
	for _, moveStr := range movesStr {
		x, err := strconv.Atoi(moveStr)
		if err == nil && x != 0 {
			moves = append(moves, x)
		}
	}
	return moves
}

func (g *GamePlay) GetGameState(move int) (*GameState, error) {
	gameState := InitialGameState()
	return gameState, nil
}

func (g *GamePlay) GetLatestGameState() (*GameState, error) {
	moves := g.GetMovesList()
	gameState := InitialGameState()

	err := gameState.MakeMove(moves[0], true)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	break
	//}

	return gameState, err
}
