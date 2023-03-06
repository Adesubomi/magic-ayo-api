package entity

import (
	"fmt"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
)

type GameState struct {
	State        [12]int `json:"state"`
	Capture0     int     `json:"capture0"`
	Capture1     int     `json:"capture1"`
	_stoneInHand int
}

func InitialGameState() *GameState {
	s := [12]int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}

	return &GameState{
		State:    s,
		Capture0: 0,
		Capture1: 0,
	}
}

func (g *GameState) MakeMove(potNumber int, player1 bool) error {
	if potNumber < 1 || potNumber > 12 {
		return logPkg.InvalidPotNumber
	}

	gameState := g.State

	sourcePotIndex := (potNumber + 11) % 12
	destinationPotIndex := potNumber % 12
	if gameState[sourcePotIndex] == 0 {
		return logPkg.EmptyPot
	}

	if g._stoneInHand == 0 {
		g._stoneInHand = gameState[sourcePotIndex]
		gameState[sourcePotIndex] = 0
	}

	gameState[destinationPotIndex] += 1
	g._stoneInHand--
	g.State = gameState

	fmt.Println(g._stoneInHand, "->", gameState)
	fmt.Println(" - - - - - - - - ")

	if g._stoneInHand > 0 {
		_ = g.MakeMove(destinationPotIndex+1, player1)
	} else if g._stoneInHand == 0 {
		if gameState[destinationPotIndex] > 3 {
			_ = g.MakeMove(destinationPotIndex+1, player1)
		} else if gameState[destinationPotIndex] == 2 ||
			gameState[destinationPotIndex] == 3 {
			g.collectEarnings()
		}
	}

	g.State = gameState
	fmt.Println("state is:", g.State)
	return nil
}

func (g *GameState) collectEarnings() {

}
