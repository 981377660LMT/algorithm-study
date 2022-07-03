# 2 <= n <= 1000
# 1 <= delay < forget <= n
from collections import defaultdict

MOD = int(1e9 + 7)


class Solution:
    def peopleAwareOfSecret1(self, n: int, delay: int, forget: int) -> int:
        """dp

        怎么定义状态
        dp[i] 表示 第i天新增的传播人数 不要定义成知道秘密的总人数
        转移:第i天的新增人数来源于(i-forget,i-delay]区间的新增人数
        答案:最后-forget天之和
        """
        dp = defaultdict(int, {1: 1})
        for pre in range(1, n + 1):
            for cur in range(pre + delay, pre + forget):
                dp[cur] += dp[pre]
        res = 0
        for cur in range(n - forget + 1, n + 1):
            res += dp[cur]
            res %= MOD
        return res

    def peopleAwareOfSecret2(self, n: int, delay: int, forget: int) -> int:
        """前缀和优化dp 思路sum(dp[pre+delay:pre+forget])]) 这一段可以前缀和优化"""
        dp = defaultdict(int, {1: 1})
        dpSum = defaultdict(int, {1: 1})
        for day in range(2, n + 1):
            dp[day] = (dpSum[day - delay] - dpSum[day - forget]) % MOD
            dpSum[day] = (dpSum[day - 1] + dp[day]) % MOD
        return (dpSum[n] - dpSum[n - forget]) % MOD

    def peopleAwareOfSecret3(self, n: int, delay: int, forget: int) -> int:
        """哈希表模拟"""
        toBorn = defaultdict(int, {1 + delay: 1})  # 传播者将在这天新增
        toDie = defaultdict(int, {1 + forget: 1})  # 传播者将在这天死去
        alive = 1  # 当前人数
        spread = 0  # 传播者人数
        for day in range(1, n + 1):
            spread += toBorn[day] - toDie[day]
            newBorn = spread
            toBorn[day + delay] += newBorn
            toDie[day + forget] += newBorn
            alive += newBorn - toDie[day]
            alive %= MOD
        return alive


print(Solution().peopleAwareOfSecret1(n=6, delay=2, forget=4))
print(Solution().peopleAwareOfSecret2(n=6, delay=2, forget=4))
print(Solution().peopleAwareOfSecret3(n=6, delay=2, forget=4))
