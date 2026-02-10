package services

import (
	"CorsGame/internal/store/sqlite"
	"fmt"
	"strings"
)

type GameService struct {
	GameStore   store.GameStore
	PlayerStore store.PlayerStore
}

func NewGameService(gs store.GameStore, ps store.PlayerStore) *GameService {
	return &GameService{
		GameStore:   gs,
		PlayerStore: ps,
	}
}

// this layer should do the business logic, or what is it called
// use the store functions to make the game logic backend-wired

func (g *GameService) ExecuteMove(gameID string, PlayerID int, CellIndex int) error {
	game, err := g.GameStore.GetGameByID(gameID)
	if err != nil {
		return err
	}
	if game.GameState != "active" {
		return fmt.Errorf("game not active")
	}
	if game.Turn_id != PlayerID {
		return fmt.Errorf("Not your Turn ")
	}

	CookedBoard := g.CookBoard(game.Board)
	if CookedBoard[CellIndex] != "" {
		return fmt.Errorf("Cell is already occupied")
	}
	playerSymbol := "X"
	if PlayerID == game.Player_o_id {
		playerSymbol = "O"
	}
	CookedBoard[CellIndex] = playerSymbol

	// win ? draw ? continue ?

	winnerFlag := g.CheckWinner(CookedBoard)
	var Winner_id int
	tempState := "active"
	if winnerFlag != "" {
		tempState = "finished"
		if winnerFlag == "X" {
			Winner_id = game.Player_x_id
		} else {
			Winner_id = game.Player_o_id
		}
	} else if g.CheckDraw(CookedBoard) {
		tempState = "draw"
	}

	board := g.UnCookBoard(CookedBoard)
	if tempState == "active" {
		err := g.GameStore.UpdateMove(gameID, PlayerID, board)
		if err != nil {
			return fmt.Errorf("failed to update Move %v", err)
		}
	}
	return g.GameStore.UpdateGame(gameID, tempState, board, Winner_id)

}

// function that checks if its draw
func (g *GameService) CheckDraw(board [9]string) bool {
	for _, cell := range board {
		if cell == "" {
			return false
		}
	}
	return true
}

func (g *GameService) CheckWinner(board [9]string) string {
	winLines := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}
	// _> thos are the possible winning indexes

	for _, line := range winLines {
		a, b, c := line[0], line[1], line[2]

		if board[a] == "" {
			break
		}
		if board[a] == board[b] && board[b] == board[c] {
			return board[a] // this here should return the SIMBOL (x or o)
			// the winner.
		}
	}
	return ""
}

// Cook board , uncook board functions to make board a 3*3 grid
func (g *GameService) CookBoard(dbString string) [9]string {

	var newBoard [9]string

	for i, char := range dbString {
		val := string(char)
		if val == "-" {
			newBoard[i] = ""
		} else {
			newBoard[i] = val
		}
	}
	return newBoard
}

// OLD Uncook board CODE, Research why its inefficient
// var dbString string
//	for _, val := range board {
//		if val == "" {
//			dbString += "-"
//		} else {
//			dbString += val
//		}
//	}
//
// return dbString

func (g *GameService) UnCookBoard(board [9]string) string {
	// im gonna use string.builder for its the GO way
	var builder strings.Builder
	builder.Grow(9)

	for _, val := range board {
		if val == "" {
			builder.WriteString("-")
		} else {
			builder.WriteString(val)
		}
	}
	return builder.String()

}
