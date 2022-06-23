from math import factorial

# 每个人与除自己外的一个人握手，所以总共会有 num_people / 2 次握手。
# 将握手的人之间连线，请你返回连线不会相交的握手方案数。
# 由于结果可能会很大，请你返回答案 模 10^9+7 后的结果。


class Solution1:
    def numberOfWays(self, num_people: int) -> int:
        n = num_people // 2
        return (factorial(2 * n) // factorial(n) // factorial(n) // (n + 1)) % (10 ** 9 + 7)


class Solution:
    def numberOfWays(self, num_people: int) -> int:
        """选择一个小朋友和编号最小的人的握手，把剩余的人群分为两个偶数人群，这样就递归到了两个更小的子问题上。"""
        MOD = 10 ** 9 + 7
        n = num_people
        dp = [0 for _ in range(n + 1)]
        dp[0] = 1
        for i in range(2, n + 1, 2):
            for j in range(0, i, 2):
                dp[i] += dp[j] * dp[i - j - 2]
                dp[i] %= MOD
        return dp[-1]


print(Solution().numberOfWays(2))
print(Solution().numberOfWays(4))
print(Solution().numberOfWays(6))
print(Solution().numberOfWays(8))
# 1 2 5 14 42 142  卡特兰数
