# 移除左右1；移除中间2
# 返回移除所有载有违禁货物车厢所需要的 最少 单位时间数

# !因为题目是两端，所以需要左右dp+枚举分割点
# 记清dp[i]的含义!
# n<=2**10^5
from typing import List


class Solution:
    def minimumTime(self, s: str) -> int:
        def getDp(string: str) -> List[int]:
            dp = [0] * len(string)
            if string[0] == '1':
                dp[0] = 1
            for i in range(1, len(string)):
                if string[i] == '0':
                    dp[i] = dp[i - 1]
                else:
                    dp[i] = min(dp[i - 1] + 2, i + 1)
            return dp

        dp1 = getDp(s)
        dp2 = getDp(s[::-1])[::-1]

        # 枚举分割点
        return min(dp1[i] + dp2[i] - int(s[i]) for i in range(len(s)))


print(Solution().minimumTime("1100101"))
