package services

import (
	"CorsGame/internal/models"
	"CorsGame/internal/store/sqlite"
	"fmt"
	"strings"
)

type GameService struct {
	GameStore   *store.GameStore
	PlayerStore *store.PlayerStore
}

func NewGameService(gs *store.GameStore, ps *store.PlayerStore) *GameService {
	return &GameService{
		GameStore:   gs,
		PlayerStore: ps,
	}
}

func (g *GameService) GetGameState(gameID string) (*models.GameRoom, error) {
	return g.GameStore.GetGameState(gameID)
}

func (g *GameService) GetGameByID(gameID string) (*models.GameRoom, error) {
	return g.GameStore.GetGameByID(gameID)
}

func (g *GameService) CreateRoom(roomID string, creatorID int) error {
	return g.GameStore.CreateNewRoom(roomID, creatorID)
}

func (g *GameService) JoinRoom(roomID string, playerID int) error {
	return g.GameStore.JoinRoom(roomID, playerID)
}

func (g *GameService) CreatePlayer(username, password string) (int, error) {
	auth := NewAuthService()
	hashed, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}
	return g.PlayerStore.CreatePlayer(username, hashed)
}

func (g *GameService) ExecuteMove(gameID string, PlayerID int, CellIndex int) error {
	if CellIndex < 0 || CellIndex > 8 {
		return fmt.Errorf("cell index %d is out of bounds", CellIndex)
	}
	game, err := g.GameStore.GetGameByID(gameID)
	if err != nil {
		return err
	}
	if game.GameState != "active" {
		return fmt.Errorf("game not active")
	}
	if game.TurnID != PlayerID {
		return fmt.Errorf("Not your Turn ")
	}

	CookedBoard := g.CookBoard(game.Board)
	if CookedBoard[CellIndex] != "" {
		return fmt.Errorf("Cell is already occupied")
	}
	playerSymbol := "X"
	if PlayerID == game.PlayerOID {
		playerSymbol = "O"
	}
	CookedBoard[CellIndex] = playerSymbol

	winnerFlag := g.CheckWinner(CookedBoard)
	var WinnerID int
	tempState := game.GameState
	if winnerFlag != "" {
		tempState = "finished"
		if winnerFlag == "X" {
			WinnerID = game.PlayerXID
		} else {
			WinnerID = game.PlayerOID
		}
	} else if g.CheckDraw(CookedBoard) {
		tempState = "draw"
	}

	board := g.UnCookBoard(CookedBoard)
	err = g.GameStore.UpdateMove(gameID, PlayerID, board)
	if err != nil {
		return fmt.Errorf("failed to update Move %v", err)
	}
	return g.GameStore.UpdateGame(gameID, tempState, board, WinnerID)
}

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
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, line := range winLines {
		a, b, c := line[0], line[1], line[2]
		if board[a] != "-" && board[a] == board[b] && board[b] == board[c] {
			return board[a]
		}
	}
	return ""
}

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

func (g *GameService) UnCookBoard(board [9]string) string {
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
