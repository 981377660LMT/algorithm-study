package main

import (
	"fmt"
)

type Board comparable

type Move any

type Nimber int64

type State[B Board, M Move] struct {
	Game any
	Move M
}

type ImpartialGame[B Board, M Move] interface {
	NextStates(board B) []State[B, M]
	IsSplittable() bool
}

type ImpartialGameSolver[B Board, M Move] struct {
	memo map[B]Nimber
	game ImpartialGame[B, M]
}

func NewImpartialGameSolver[B Board, M Move](game ImpartialGame[B, M]) *ImpartialGameSolver[B, M] {
	return &ImpartialGameSolver[B, M]{
		memo: make(map[B]Nimber),
		game: game,
	}
}

func (s *ImpartialGameSolver[B, M]) getGameNimber(game any) Nimber {
	switch g := game.(type) {
	case B:
		return s.GetNimber(g)
	case []B:
		var n Nimber = 0
		for _, board := range g {
			n ^= s.GetNimber(board)
		}
		return n
	default:
		panic("invalid game type")
	}
}

func (s *ImpartialGameSolver[B, M]) GetNimber(board B) Nimber {
	if nimber, ok := s.memo[board]; ok {
		return nimber
	}

	nimber := s.calculateNimber(board)
	s.memo[board] = nimber
	return nimber
}

func (s *ImpartialGameSolver[B, M]) calculateNimber(board B) Nimber {
	nextStates := s.game.NextStates(board)
	if len(nextStates) == 0 {
		return 0
	}

	nextNimbers := make(map[Nimber]struct{})
	for _, state := range nextStates {
		nextNimbers[s.getGameNimber(state.Game)] = struct{}{}
	}

	var mex Nimber = 0
	for {
		if _, exists := nextNimbers[mex]; !exists {
			return mex
		}
		mex++
	}
}

func (s *ImpartialGameSolver[B, M]) GetBestMove(game any) any {
	totalNimber := s.getGameNimber(game)
	if totalNimber == 0 {
		panic("No best move in a losing position.")
	}

	if !s.game.IsSplittable() {
		board := game.(B)
		for _, state := range s.game.NextStates(board) {
			if s.getGameNimber(state.Game) == 0 {
				return state.Move
			}
		}
	} else {
		boards := game.([]B)
		for i, subBoard := range boards {
			subNimber := s.GetNimber(subBoard)
			targetNimber := totalNimber ^ subNimber
			if targetNimber < subNimber {
				for _, state := range s.game.NextStates(subBoard) {
					if s.getGameNimber(state.Game) == targetNimber {
						return struct {
							Index int
							Move  M
						}{i, state.Move}
					}
				}
			}
		}
	}

	panic("Error in GetBestMove: logic failed.")
}

type NimMove struct {
	Index int
	Count int
}

type NimGame struct{}

func (ng *NimGame) IsSplittable() bool {
	return true
}

func (ng *NimGame) NextStates(board int) []State[int, NimMove] {
	states := []State[int, NimMove]{}
	for i := 1; i <= board; i++ {
		states = append(states, State[int, NimMove]{
			Game: board - i,
			Move: NimMove{Count: i},
		})
	}
	return states
}

func main() {
	nimGame := &NimGame{}
	solver := NewImpartialGameSolver(nimGame)

	initialBoard := []int{3, 4, 5}

	nimber := solver.getGameNimber(initialBoard)
	fmt.Printf("Initial board %v has Nimber: %d\n", initialBoard, nimber)

	if nimber == 0 {
		fmt.Println("It's a losing position (P-position).")
	} else {
		fmt.Println("It's a winning position (N-position).")
		bestMove := solver.GetBestMove(initialBoard)
		fmt.Printf("Best move is: %+v\n", bestMove)
	}
}
