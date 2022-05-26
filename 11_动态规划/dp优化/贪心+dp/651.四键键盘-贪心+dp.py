# 假设你有一个特殊的键盘包含下面的按键：
# A：在屏幕上打印一个 'A'。
# Ctrl-A：选中整个屏幕。
# Ctrl-C：复制选中区域到缓冲区。
# Ctrl-V：将缓冲区内容输出到上次输入的结束位置，并显示在屏幕上。
# 现在，你可以 最多 按键 n 次（使用上述四种按键），返回屏幕上最多可以显示 'A' 的个数 。
#
from functools import lru_cache


@lru_cache(None)
def dfs(remain: int) -> int:
    if remain <= 0:
        return 0 if remain == 0 else -int(1e20)

    res = dfs(remain - 1) + 1  # 打字
    for count in range(remain):  # 贪心，到最终执行几步复制
        res = max(res, dfs(remain - count - 1) * count)
    return res


class Solution:
    def maxA(self, n: int) -> int:
        """肯定是从某个时间开始不断按v"""
        return dfs(n)

    def maxA2(self, n: int) -> int:
        """肯定是从某个时间开始不断按v"""
        dp = list(range(n + 1))
        for i in range(1, n + 1):
            for j in range(2, i):
                dp[i] = max(dp[i], dp[j - 2] * (i - (j + 1) + 2))

        return dp[-1]


print(Solution().maxA(7))
print(Solution().maxA2(11))
