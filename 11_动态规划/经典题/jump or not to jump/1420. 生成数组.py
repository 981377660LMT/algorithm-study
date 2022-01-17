MOD = int(1e9 + 7)
from functools import lru_cache
from itertools import product

# 描述的算法用于找出正整数数组中最大的元素。
# 请你生成一个具有下述属性的数组 arr ：

# arr 中有 n 个整数。
# 1 <= arr[i] <= m 其中 (0 <= i < n) 。
# 将上面提到的算法应用于 arr ，search_cost 的值等于 k(jump ) 。
# 返回上述条件下生成数组 arr 的 方法数
# 1 <= n <= 50
# 1 <= m <= 100
# 0 <= k <= n

# dp[i][j][k]代表符合长度为n,跳跃数为j,最大值确定为k这三个条件的数组个数


# summary:
# jump or not jump
class Solution:
    def numOfArrays(self, n: int, m: int, k: int) -> int:
        dp = [[[0 for _ in range(m + 1)] for _ in range(k + 1)] for _ in range(n + 1)]

        for maxNum in range(1, m + 1):
            dp[1][1][maxNum] = 1

        for length, jump, maxNum in product(range(1, n + 1), range(1, k + 1), range(1, m + 1)):
            # dont jump: curVal is smaller than or equal to k.
            dp[length][jump][maxNum] += dp[length - 1][jump][maxNum] * maxNum
            # jump:curVal is bigger than preVal.
            dp[length][jump][maxNum] += sum(dp[length - 1][jump - 1][1:maxNum])

        return sum(dp[n][k][1:]) % int(1e9 + 7)


print(Solution().numOfArrays(n=2, m=3, k=1))
# 输出：6
# 解释：可能的数组分别为 [1, 1], [2, 1], [2, 2], [3, 1], [3, 2] [3, 3]
