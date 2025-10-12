```go
// “平方剩余核”（square-free kernel）：把数中偶数次出现的质因子全部抵消后剩下的乘积.
// 若两个数 a、b 的平方剩余核相同 (sf[a]==sf[b])，则 a·b 为完全平方数.
func GetSquareFreeKernel(n int, getPrimeFactors func(int) map[int]int) int {
	res := 1
	factors := getPrimeFactors(n)
	for p, c := range factors {
		if c&1 == 1 {
			res *= p
		}
	}
	return res
}
```
