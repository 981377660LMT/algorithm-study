# 1 <= capacities.length <= 200
# 1 <= capacities[i] <= 200
# 0 <= k <= sum(capacities)

# 两种类新的安检室：栈类型和队列类型
# 若可以任意设定每个安检室的类型，请问有多少种设定安检室类型的方案可以使得编号 k 的观众第一个通过最后一个安检室入场。
# 观众不可主动离开安检室，只有当安检室容纳人数达到上限，且又有新观众需要进入时，才可根据安检室的类型选择一位观众离开；
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def securityCheck(self, capacities: List[int], k: int) -> int:
        """
        我们可以对栈填充 capacities-1 个人，这样可以将栈转换成一个大小为 1 的队列
        问题转化为选出一些安检室，使得安检室 capacities-1 之和恰好等于k，即01背包求方案数
        """
        dp = [0] * (k + 1)
        dp[0] = 1
        for i in range(len(capacities)):
            cost = capacities[i] - 1
            for j in range(k, cost - 1, -1):
                dp[j] += dp[j - cost]
                dp[j] %= MOD
        return dp[-1]


print(Solution().securityCheck(capacities=[2, 2, 3], k=2))
