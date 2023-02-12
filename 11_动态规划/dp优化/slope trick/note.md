## slope trick

https://codeforces.com/blog/entry/47821
https://leetcode.cn/problems/make-array-non-decreasing-or-non-increasing/solution/xie-lu-you-hua-by-kkxbb-ufrf/

slope trick 的 dp 形式
dp[i][j] = min{dp[i-1][k] + cost(i,j)} k 的范围取决于 i 和 j
即:`多个绝对值函数的部分区间的叠加`,是一个凸函数即
**f(x) = minf + ∑(lefti - x)+ + ∑(x - righti)+**

- 左右的斜率是对称的
- f(x)在[left0,right0]区间的函数值相等 (斜率为 0)

slope trick 的思想就是用一个数据结构管理这些`斜率的转折点`

朴素的 dp 是`O(nk)`的
slope trick 是`O(nlogn)`的

- 斜率变化+-1 时用堆维护，不为+-1 时用平衡树维护
  Generalized-Slope-Trick
  https://ei1333.github.io/library/structure/others/generalized-slope-trick.hpp
