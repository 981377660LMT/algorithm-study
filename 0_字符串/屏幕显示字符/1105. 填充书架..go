// 1105. 填充书架
// https://leetcode.cn/problems/filling-bookcase-shelves/description/
// 摆放书的顺序与你整理好的顺序相同。
// 以这种方式布置书架，返回书架整体可能的最小高度。
// 1 <= books.length <= 1000
// 1 <= thicknessi <= shelfWidth <= 1000
// 1 <= heighti <= 1000

package main

const INF int = 1e18

func minHeightShelves(books [][]int, shelfWidth int) int {
	n := len(books)
	dp := make([]int, n+1) // dp[0]=0
	for i := 1; i <= n; i++ {
		w, h := 0, 0
		dp[i] = INF
		for j := i; j >= 1; j-- { // 把 j..i 放同一层
			w += books[j-1][0]
			if w > shelfWidth {
				break
			}
			h = max(h, books[j-1][1])
			if cand := dp[j-1] + h; cand < dp[i] {
				dp[i] = cand
			}
		}
	}
	return dp[n]
}
