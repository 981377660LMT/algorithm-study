https://leetcode-cn.com/problems/non-negative-integers-without-consecutive-ones/solution/shu-wei-dpmo-ban-ji-jie-fa-by-initness-let3/
对于「数位 DP」题，都存在「`询问 [a, b]（a 和 b 均为正整数，且 a < b）区间内符合条件的数值个数为多少`」的一般形式，通常我们需要实现一个查询 **[0, x] 有多少**合法数值的函数 int dp(int x)，然后应用前缀和求解出 [a, b] 的个数：dp(b)−dp(a−1)。

这类题的特征是 `N在10^9量级` 要根据数位而不是遍历数值来做

`357. 计算各个位数不同的数字个数`
`600. 不含连续1的非负整数.ts`
`902. 最大为 N 的数字组合.ts`
`1012. 至少有 1 位重复的数字.py`
`1067. 范围内的数字计数`

https://www.acwing.com/blog/content/7944/

DFS 做法才是数位 DP 的正解

数位 DP 问题一般给定一个区间 [L,R]，问区间满足的条件的数有多少个。数字范围很大。
可以利用前缀和来求解答案

数位 dfs 模板

1. 把数字拆成 nums
2. `dfs(pos,isLimit,*rest)` 开始 dfs(len(nums),True,...)
3. 结束条件为 pos==len(nums)
4. 每次选下一位数先确定 up `up = nums[pos] if isLimit else 9`

假设数字 x 位数为 a1⋯an
两个必带参数：
pos:pos: 表示数字的位数
isLimit: 可以填数的限制（无限制的话(isLimit=True) 0∼9 随便填，否则只能填到 up

四个选带参数：
pre:表示上一个数是多少
hasLeadingZero:前导零是否存在，表示是否开始选数了
`注意：如果题目要求正整数，那么 dfs 终点不能有 leadingZero`
curSum:搜索到当前所有数字之和
count: 某个数字出现的次数
`经典计数问题`

```Python
from functools import lru_cache


@lru_cache(None)
def cal(upper: int, queryDigit: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, count: int) -> int:
        """当前在第pos位，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return count

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), count)
            else:
                res += dfs(pos + 1, False, (isLimit and cur == up), count + (cur == queryDigit))
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, True, 0)


class Solution:
    def digitsCount(self, d: int, low: int, high: int) -> int:
        # 1 <= low <= high <= 2×10^8
        return cal(high, d) - cal(low - 1, d)
```
