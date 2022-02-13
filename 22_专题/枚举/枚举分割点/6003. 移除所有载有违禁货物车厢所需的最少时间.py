# 移除左右1；移除中间2
# 返回移除所有载有违禁货物车厢所需要的 最少 单位时间数

# 因为题目是两端，所以需要左右dp+枚举分割点
# 记清dp[i]的含义!
# n<=2**10^5
class Solution:
    def minimumTime(self, s: str) -> int:
        if '1' not in s:
            return 0
        n = len(s)

        # 从左向右移除
        dp1 = [0] * n
        if s[0] == '1':
            dp1[0] = 1
        for i in range(1, n):
            if s[i] == '0':
                dp1[i] = dp1[i - 1]
            else:
                dp1[i] = min(dp1[i - 1] + 2, i + 1)

        # 从右向左移除
        dp2 = [0] * n
        if s[-1] == '1':
            dp2[-1] = 1
        for i in range(n - 2, -1, -1):
            if s[i] == '0':
                dp2[i] = dp2[i + 1]
            else:
                dp2[i] = min(dp2[i + 1] + 2, n - i)

        # 枚举分割点
        res = n
        for i in range(n - 1):
            res = min(res, dp1[i] + dp2[i + 1])
        return res


print(Solution().minimumTime("1100101"))
