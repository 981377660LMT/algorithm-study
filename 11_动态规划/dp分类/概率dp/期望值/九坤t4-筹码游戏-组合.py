"""
状压+期望dp
注意这题与J - ボール的区别:
J - ボール 是固定了排列顺序
而这道题任意排列顺序 即组合

纸箱中共有 kind 种面值的筹码。现给定九坤取出筹码的最终目标为 nums
nums[i] 表示第 i 堆筹码的数量。
假设每种面值的筹码都有无限多个，且九坤总是遵循最优策略，
使得他达成目标的操作次数最小化。
请返回九坤达成目标的情况下，需要取出筹码次数的期望值。

1 <= kind <= 50
1 <= nums.length <= kind
sum(nums[0],nums[1],...,nums[n]) <= 50
"""

# !状压+概率dp 答案只与已经选择的筹码状态有关
# !由于 sum(nums)<50，那么状态数有一个上界是对 sum(state) 的整数拆分
# !50 的整数拆分大约在 2e5的量级

# dp[state] = 1 + (dp[state1]  + dp[state2]  + ... + dp[statek]) / k
# !注意如果nextState和state一样的话 需要合并同类项dp[state]


from functools import lru_cache
from typing import List, Tuple


class Solution:
    def chipGame(self, nums: List[int], kind: int) -> float:
        @lru_cache(None)
        def dfs(state: Tuple[int, ...], remain: int) -> float:
            """state表示当前筹码的`组合` remain表示还要取的筹码数量"""
            if remain == 0:
                return 0

            lis = list(state)
            step, p = 1, 0
            for i in range(kind):
                lis[i] += 1
                nextState = tuple(sorted(lis))
                if all(nextState[j] <= target[j] for j in range(kind)):
                    step += dfs(nextState, remain - 1) / kind
                    p += 1 / kind
                lis[i] -= 1

            return step / p

        target = sorted(nums + [0] * (kind - len(nums)))
        res = dfs(tuple([0] * kind), sum(nums))  # 正着推导 补齐0
        dfs.cache_clear()
        return res


print(Solution().chipGame(nums=[1, 1], kind=2))
print(Solution().chipGame(nums=[1, 2], kind=4))
print(Solution().chipGame(nums=[1, 1, 1], kind=3))  # 5.5
print(Solution().chipGame(nums=[3], kind=1))  # 3.0
