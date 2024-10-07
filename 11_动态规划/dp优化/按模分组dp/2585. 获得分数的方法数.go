// 2585. 获得分数的方法数(分组背包)
// https://leetcode.cn/problems/number-of-ways-to-earn-points/description/
//
// 考试中有 n 种类型的题目。给你一个整数 target 和一个下标从 0 开始的二维整数数组 types ，
// 其中 types[i] = [counti, marksi] 表示第 i 种类型的题目有 counti 道，每道题目对应 marksi 分。
// !返回你在考试中恰好得到 target 分的方法数。由于答案可能很大，结果需要对 1e9 +7 取余。
// !注意，同类型题目无法区分。
// target<=1000
// n<=50
// counti<=50

// !O(n*target) 按模分组前缀和优化dp
// dp[i][j]表示前i种题目恰好得到j分的方法数
// !dp[j] = sum(dp[j-k*mark] for k in range(count+1) if j-k*mark>=0
//
// 将i按照模mark分组，令i' = div * mark + mod，则
// !dp[i] = sum(g[j]) (0<=j<i, i-j<=count)

package main

const MOD int = 1e9 + 7

func waysToReachTarget(target int, types [][]int) int {
	dp := make([]int, target+1)
	dp[0] = 1
	g := []int{}
	for _, t := range types {
		count, mark := t[0], t[1]
		for mod := 0; mod < mark; mod++ {
			g = g[:0]
			for pos := mod; pos <= target; pos += mark {
				g = append(g, dp[pos])
			}

			Q := NewSizedQueue(count)
			for j, pos := 0, mod; j < len(g); j, pos = j+1, pos+mark {
				dp[pos] = (dp[pos] + Q.Sum()) % MOD
				Q.Append(g[j])
			}
		}
	}

	return dp[target]
}

type SizedQueue struct {
	sum   int
	size  int
	queue []int
}

func NewSizedQueue(size int) *SizedQueue {
	return &SizedQueue{size: size}
}

func (q *SizedQueue) Append(v int) {
	q.queue = append(q.queue, v)
	q.sum += v
	if len(q.queue) > q.size {
		q.sum -= q.queue[0]
		q.queue = q.queue[1:]
	}
}

func (q *SizedQueue) PopLeft() int {
	res := q.queue[0]
	q.queue = q.queue[1:]
	q.sum -= res
	return res
}

func (q *SizedQueue) Sum() int {
	return q.sum
}

func (q *SizedQueue) Len() int {
	return len(q.queue)
}

func (q *SizedQueue) At(i int) int {
	return q.queue[i]
}
