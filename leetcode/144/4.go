package main

import (
	"math"
)

type Position struct {
	i int
	j int
}

type Move struct {
	di int
	dj int
}

var moves = map[int][]Move{
	1: {
		{1, 1},
	},
	2: {
		{1, -1},
		{1, 0},
		{1, 1},
	},
	3: {
		{-1, 1},
		{0, 1},
		{1, 1},
	},
}

func computeMaxPath(n int, kidNum int, grid [][]int) (int, []Position) {
	var startI, startJ int
	switch kidNum {
	case 2:
		startI, startJ = 0, n-1
	case 3:
		startI, startJ = n-1, 0
	default:
		return 0, nil
	}

	type State struct {
		i       int
		j       int
		balance int
	}

	var maxBalance int
	if kidNum == 2 {
		maxBalance = n - 1
	} else if kidNum == 3 {
		maxBalance = n - 1
	}

	dpCurrent := make([][]int, 2)
	for k := 0; k < 2; k++ {
		dpCurrent[k] = make([]int, 2*maxBalance+1)
		for b := range dpCurrent[k] {
			dpCurrent[k][b] = -math.MaxInt32
		}
	}
	initialSum := grid[startI][startJ]
	dpCurrent[0][maxBalance] = initialSum

	for step := 0; step < n-1; step++ {
		nextStep := (step + 1) % 2
		for b := range dpCurrent[nextStep] {
			dpCurrent[nextStep][b] = -math.MaxInt32
		}

		for balance := 0; balance < 2*maxBalance+1; balance++ {
			currentSum := dpCurrent[step%2][balance]
			if currentSum == -math.MaxInt32 {
				continue
			}
			for _, move := range moves[kidNum] {
				ni := startI + (step+1)*move.di
				nj := startJ + (step+1)*move.dj

			}
		}
		dpCurrent = dpNext
	}

	dp := make([][][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([][]int, n)
		for j := 0; j < n; j++ {
			dp[i][j] = make([]int, 2*maxBalance+1)
			for b := 0; b < 2*maxBalance+1; b++ {
				dp[i][j][b] = -math.MaxInt32
			}
		}
	}
	dp[startI][startJ][maxBalance] = grid[startI][startJ]

	for step := 0; step < n-1; step++ {
		newDP := make([][][]int, n)
		for i := 0; i < n; i++ {
			newDP[i] = make([][]int, n)
			for j := 0; j < n; j++ {
				newDP[i][j] = make([]int, 2*maxBalance+1)
				for b := 0; b < 2*maxBalance+1; b++ {
					newDP[i][j][b] = -math.MaxInt32
				}
			}
		}

		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				for b := 0; b < 2*maxBalance+1; b++ {
					if dp[i][j][b] == -math.MaxInt32 {
						continue
					}
					for _, move := range moves[kidNum] {
						ni := i + move.di
						nj := j + move.dj
						if ni < 0 || ni >= n || nj < 0 || nj >= n {
							continue
						}
						newBalance := b
						if kidNum == 2 {
							newBalance += move.dj
						} else if kidNum == 3 {
							newBalance += -move.di
						}
						if newBalance < 0 || newBalance > 2*maxBalance {
							continue
						}
						newSum := dp[i][j][b] + grid[ni][nj]
						if newSum > newDP[ni][nj][newBalance] {
							newDP[ni][nj][newBalance] = newSum
						}
					}
				}
			}
		}
		dp = newDP
	}

	maxSum := -math.MaxInt32
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for b := 0; b < 2*maxBalance+1; b++ {
				if (kidNum == 2 && b == maxBalance) ||
					(kidNum == 3 && b == maxBalance) {
					if dp[i][j][b] > maxSum {
						maxSum = dp[i][j][b]
					}
				}
			}
		}
	}

	if maxSum == -math.MaxInt32 {
		return 0
	}
	return maxSum
}

func maxCollectedFruits(fruits [][]int) int {
	n := len(fruits)
	if n == 0 {
		return 0
	}

	kid1Sum := 0
	gridCopy := make([][]int, n)
	for i := 0; i < n; i++ {
		gridCopy[i] = make([]int, n)
		copy(gridCopy[i], fruits[i])
		kid1Sum += gridCopy[i][i]
		gridCopy[i][i] = 0
	}

	orderings := [][]int{
		{2, 3},
		{3, 2},
	}

	maxTotal := 0

	for _, order := range orderings {
		currentSum := kid1Sum
		tempGrid := make([][]int, n)
		for i := 0; i < n; i++ {
			tempGrid[i] = make([]int, n)
			copy(tempGrid[i], gridCopy[i])
		}

		paths := make([][][]int, 4)
		for _, kid := range order {
			sum := computeMaxPath(n, kid, tempGrid)
			currentSum += sum
		}

		if currentSum > maxTotal {
			maxTotal = currentSum
		}
	}

	return maxTotal
}
