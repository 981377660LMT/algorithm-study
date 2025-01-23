package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc386_f()
}

func abc386_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int
	fmt.Fscan(in, &k)
	var s string
	fmt.Fscan(in, &s)
	var t string
	fmt.Fscan(in, &t)

	d := EditDistanceK([]byte(s), []byte(t), k+1)
	if d <= k {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}

// https://leetcode.cn/problems/edit-distance/
func minDistance(word1 string, word2 string) int {
	return EditDistanceK([]byte(word1), []byte(word2), 1<<31-1)
}

// 计算两个数组的编辑距离，允许最多k次操作.超过k次操作返回k.
// !时间复杂度O(k*min(n,m)), 空间复杂度O(k).
// TODO:
// 线性长度复杂度的编辑距离算法，时间复杂度 O(N+M+K^2).
// https://atcoder.jp/contests/abc386/editorial/11701
// https://rsk0315.hatenablog.com/entry/2022/12/30/180216 入力長に対して線形時間 + 編集距離に対して二乗時間で LCS を求めるよ
func EditDistanceK[T comparable](a []T, b []T, k int) int {
	aSize := len(a)
	bSize := len(b)
	if aSize > bSize {
		a, b = b, a
		aSize, bSize = bSize, aSize
	}

	for aSize > 0 && a[aSize-1] == b[bSize-1] {
		aSize--
		bSize--
	}

	prefix := 0
	for prefix < aSize && a[prefix] == b[prefix] {
		prefix++
	}
	if prefix > 0 {
		a = a[prefix:]
		b = b[prefix:]
		aSize -= prefix
		bSize -= prefix
	}

	if aSize == 0 {
		return min(k, bSize)
	}

	k = min(bSize, k)
	sizeD := bSize - aSize
	if sizeD > k {
		return k
	}

	zeroK := min(k, aSize)/2 + 2
	arraySize := sizeD + zeroK*2 + 2

	currentRow := make([]int, arraySize)
	nextRow := make([]int, arraySize)
	for i := range currentRow {
		currentRow[i] = -1
		nextRow[i] = -1
	}

	i := 0
	conditionRow := sizeD + zeroK
	endMax := conditionRow * 2

	for {
		i++
		currentRow, nextRow = nextRow, currentRow

		var start int
		var nextCell int

		if i <= zeroK {
			start = -i + 1
			nextCell = i - 2
		} else {
			start = i - zeroK*2 + 1
			nextCell = currentRow[zeroK+start]
		}

		end := 0
		if i <= conditionRow {
			end = i
			nextRow[zeroK+i] = -1
		} else {
			end = endMax - i
		}

		previousCell := -1
		currentCell := -1
		for q, rowIndex := start, start+zeroK; q < end; q, rowIndex = q+1, rowIndex+1 {
			previousCell, currentCell = currentCell, nextCell
			if rowIndex+1 < len(currentRow) {
				nextCell = currentRow[rowIndex+1]
			} else {
				nextCell = -1
			}

			t := max(currentCell+1, max(previousCell, nextCell+1))

			for t < aSize && t+q < bSize && a[t] == b[t+q] {
				t++
			}

			if rowIndex < len(nextRow) {
				nextRow[rowIndex] = t
			}
		}

		if nextRow[conditionRow] >= aSize || i > k {
			break
		}
	}

	return i - 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
