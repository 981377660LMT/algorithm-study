https://mugen1337.github.io/procon/tips/AlienDP.hpp
AlienDp(wqs 二分)

`某个东西使用k次，最小化总花费`
需要高速化 dp[index][使用次数] 的 dp 时,可以使用 wqs 二分
不使用 k 次，而是每使用一次罚款 x 元,转化为 dp[index] ，x 二分搜索即可

---

https://github.com/EndlessCheng/codeforces-go/blob/13510980250d7b265d308ea4ff5b51eeb1b5e3e4/copypasta/dp.go#L2781
