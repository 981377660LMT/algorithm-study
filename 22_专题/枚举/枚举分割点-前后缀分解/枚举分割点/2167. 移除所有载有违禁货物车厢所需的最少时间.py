# 从列车 左 端移除一节车厢（即移除 s[0]），用去 1 单位时间。
# 从列车 右 端移除一节车厢（即移除 s[s.length - 1]），用去 1 单位时间。
# 从列车车厢序列的 任意位置 移除一节车厢，用去 2 单位时间。
# 返回移除所有载有违禁货物车厢所需要的 最少 单位时间数


# !因为题目是两端，所以需要左右dp+枚举分割点
# 记清dp[i]的含义!
# n<=2**10^5
from typing import List


class Solution:
    def minimumTime(self, s: str) -> int:
        def makeDp(string: str) -> List[int]:
            dp = [0] * len(string)
            if string[0] == "1":
                dp[0] = 1
            for i in range(1, len(string)):
                if string[i] == "0":
                    dp[i] = dp[i - 1]
                else:
                    dp[i] = min(dp[i - 1] + 2, i + 1)
            return dp

        dp1 = makeDp(s)
        dp2 = makeDp(s[::-1])[::-1]

        # 枚举分割点
        return min(dp1[i] + dp2[i] - int(s[i]) for i in range(len(s)))


print(Solution().minimumTime("1100101"))
