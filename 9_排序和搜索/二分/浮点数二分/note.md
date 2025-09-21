# 浮点数二分

https://leetcode.cn/problems/separate-squares-i/solutions/3076424/zheng-shu-er-fen-pythonjavacgo-by-endles-8yn5/
推荐的写法是固定一个循环次数，因为浮点数有舍入误差，可能算出的 mid 和 left 相等，此时 left=mid 不会更新 left，导致死循环。

```go
// > x 的下一个浮点数
func PrevFloat64(x float64) float64 {
	return math.Nextafter(x, -math.MaxFloat64)
}

// < x 的下一个浮点数
func NextFloat64(x float64) float64 {
	return math.Nextafter(x, math.MaxFloat64)
}
```

---

技巧：避免中途浮点数运算
https://atcoder.jp/contests/abc424/editorial/13878
