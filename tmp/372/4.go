package main

import "math"

func main() {

}

// def min(a: int, b: int) -> int:
//     return a if a < b else b

// def max(a: int, b: int) -> int:
//     return a if a > b else b

// def createBlock(
//     n: int, blockSize: Optional[int] = None
// ) -> Tuple[int, List[int], List[int], List[int]]:
//     if blockSize is None:
//         blockSize = int(n**0.5) + 1
//     blockCount = 1 + (n // blockSize)
//     belong = [i // blockSize for i in range(n)]
//     blockStart = [i * blockSize for i in range(blockCount)]
//     blockEnd = [min((i + 1) * blockSize, n) for i in range(blockCount)]
//     return blockCount, belong, blockStart, blockEnd

// class Solution:
//     def leftmostBuildingQueries(self, heights: List[int], queries: List[List[int]]) -> List[int]:
//         def findRightNearestHigher(start: int, target: int) -> int:
//             """找到[start, n)中第一个严格大于target的下标.如果不存在,返回-1."""
//             bid = belong[start]
//             for i in range(start, blockEnd[bid]):  # 散块
//                 if heights[i] > target:
//                     return i
//             for i in range(bid + 1, blockCount):  # 整块
//                 if blockMax[i] > target:
//                     for j in range(blockStart[i], blockEnd[i]):
//                         if heights[j] > target:
//                             return j
//             return -1

// n = len(heights)
// blockCount, belong, blockStart, blockEnd = createBlock(len(heights), 2 * int(n**0.5) + 1)
// blockMax = [INF] * blockCount
// for i, v in enumerate(heights):
//
//	bid = belong[i]
//	blockMax[bid] = max(blockMax[bid], v)
//
// res = [-1] * len(queries)
// for qi, (alice, bob) in enumerate(queries):
//
//	if alice == bob:
//	    res[qi] = alice
//	    continue
//	if alice > bob:
//	    alice, bob = bob, alice
//	if heights[alice] < heights[bob]:
//	    res[qi] = bob
//	    continue
//	res[qi] = findRightNearestHigher(bob, heights[alice])
//
// return res

func leftmostBuildingQueries(heights []int, queries [][]int) []int {
	n := len(heights)
	block := CreateBlock(len(heights), int(math.Sqrt(float64(n)))+1)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockMax := make([]int, blockCount)
	for i, v := range heights {
		bid := belong[i]
		blockMax[bid] = max(blockMax[bid], v)
	}

	// 找到[start, n)中第一个严格大于target的下标.如果不存在,返回-1.
	findRightNearestHigher := func(start int, target int) int {
		bid := belong[start]
		for i := start; i < blockEnd[bid]; i++ { // 散块
			if heights[i] > target {
				return i
			}
		}
		for i := bid + 1; i < blockCount; i++ { // 整块
			if blockMax[i] > target {
				for j := blockStart[i]; j < blockEnd[i]; j++ {
					if heights[j] > target {
						return j
					}
				}
			}
		}
		return -1
	}

	res := make([]int, len(queries))
	for qi, query := range queries {
		alice, bob := query[0], query[1]
		if alice == bob {
			res[qi] = alice
			continue
		}
		if alice > bob {
			alice, bob = bob, alice
		}
		if heights[alice] < heights[bob] {
			res[qi] = bob
			continue
		}
		res[qi] = findRightNearestHigher(bob, heights[alice])
	}

	return res
}

func CreateBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
