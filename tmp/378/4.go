package main

func main() {

}

// 给你一个长度为 偶数 n ，下标从 0 开始的字符串 s 。

// 同时给你一个下标从 0 开始的二维整数数组 queries ，其中 queries[i] = [ai, bi, ci, di] 。

// 对于每个查询 i ，你需要执行以下操作：

// 将下标在范围 0 <= ai <= bi < n / 2 内的 子字符串 s[ai:bi] 中的字符重新排列。
// 将下标在范围 n / 2 <= ci <= di < n 内的 子字符串 s[ci:di] 中的字符重新排列。
// 对于每个查询，你的任务是判断执行操作后能否让 s 变成一个 回文串 。

// 每个查询与其他查询都是 独立的 。

// 请你返回一个下标从 0 开始的数组 answer ，如果第 i 个查询执行操作后，可以将 s 变为一个回文串，那么 answer[i] = true，否则为 false 。

// 子字符串 指的是一个字符串中一段连续的字符序列。
// s[x:y] 表示 s 中从下标 x 到 y 且两个端点 都包含 的子字符串。
func canMakePalindromeQueries(s string, queries [][]int) []bool {

}
