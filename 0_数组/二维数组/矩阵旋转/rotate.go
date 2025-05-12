package main

import (
	"fmt"
)

// 顺时针旋转矩阵90度.
func RotateRight[T any](matrix [][]T) [][]T {
	if len(matrix) == 0 {
		return matrix
	}

	rows, cols := len(matrix), len(matrix[0])
	result := make([][]T, cols)
	for i := range result {
		result[i] = make([]T, rows)
	}

	for i := range cols {
		for j := range rows {
			result[i][j] = matrix[rows-1-j][i]
		}
	}

	return result
}

// 逆时针旋转矩阵90度.
func RotateLeft[T any](matrix [][]T) [][]T {
	if len(matrix) == 0 {
		return matrix
	}

	rows, cols := len(matrix), len(matrix[0])
	result := make([][]T, cols)
	for i := range result {
		result[i] = make([]T, rows)
	}

	for i := range cols {
		for j := range rows {
			result[i][j] = matrix[j][cols-1-i]
		}
	}

	return result
}

// 矩阵转置.
func Transpose[T any](matrix [][]T) [][]T {
	if len(matrix) == 0 {
		return matrix
	}

	rows, cols := len(matrix), len(matrix[0])
	result := make([][]T, cols)
	for i := range result {
		result[i] = make([]T, rows)
	}

	for i := range rows {
		for j := range cols {
			result[j][i] = matrix[i][j]
		}
	}

	return result
}

// 检查两个矩阵是否相等
func equal[T comparable](a, b [][]T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func main() {
	// 原矩阵:     右旋转:     左旋转:
	// 1 2 3      7 4 1      3 6 9
	// 4 5 6  ->  8 5 2  or  2 5 8
	// 7 8 9      9 6 3      1 4 7
	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	expected1 := [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}

	expected2 := [][]int{
		{3, 6, 9},
		{2, 5, 8},
		{1, 4, 7},
	}

	expected3 := [][]int{
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
	}

	rotated1 := RotateRight(grid)
	rotated2 := RotateLeft(grid)
	transposed := Transpose(grid)

	fmt.Println("RotateRight correct:", equal(rotated1, expected1))
	fmt.Println("RotateLeft correct:", equal(rotated2, expected2))
	fmt.Println("Transpose correct:", equal(transposed, expected3))

	// 打印矩阵以便可视化结果
	fmt.Println("\n原矩阵:")
	printMatrix(grid)

	fmt.Println("\n顺时针旋转:")
	printMatrix(rotated1)

	fmt.Println("\n逆时针旋转:")
	printMatrix(rotated2)

	fmt.Println("\n矩阵转置:")
	printMatrix(transposed)
}

// 打印矩阵
func printMatrix[T any](matrix [][]T) {
	for _, row := range matrix {
		fmt.Println(row)
	}
}
