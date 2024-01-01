// https://github.com/shanzi/algo-ds/blob/master/mayersdiff/distance.go
// https://cloud.tencent.com/developer/article/1724029
// 时间复杂度: 期望为O(M+N+D^2)，最坏情况为为O((M+N)D)

package main

// https://leetcode.cn/problems/edit-distance/
func minDistance(word1 string, word2 string) int {
	return EditingDistance(word1, word2)
}

type Str = string

// Myers算法求编辑距离(并不准确).
func EditingDistance(a, b Str) int {
	n1 := int32(len(a))
	n2 := int32(len(b))
	max_ := n1 + n2

	V := make([]int32, max_+2)
	for i := int32(0); i < int32(len(V)); i++ {
		V[i] = max_
	}

	setV(V, 1, 0)
	for d := int32(0); d <= max_; d++ {
		var x int32
		for k := -(d - 2*upzero(d-n1)); k <= d-2*upzero(d-n2); k += 2 {
			if k == -d || (k != d && getV(V, k-1) < getV(V, k+1)) {
				x = getV(V, k+1)
			} else {
				x = getV(V, k-1) + 1
			}
			y := x - k
			for x < n1 && y < n2 && a[x] == b[y] {
				x += 1
				y += 1
			}

			setV(V, k, x)

			if x == n1 && y == n2 {
				return int(d)
			}
		}
	}
	return int(max_)
}

func upzero(a int32) int32 {
	if a < 0 {
		return 0
	}
	return a
}

func getV(V []int32, idx int32) int32 {
	if idx >= 0 {
		return V[idx]
	} else {
		return V[int32(len(V))+idx]
	}
}

func setV(V []int32, idx int32, val int32) {
	if idx >= 0 {
		V[idx] = val
	} else {
		V[int32(len(V))+idx] = val
	}
}
