"""
状压+期望dp
注意这题与J - ボール的区别:
J - ボール 是固定了排列顺序
而这道题任意排列顺序 即组合 所以用counter保存状态

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
# !由于 sum(nums)<50，那么状态数有一个上界是对 sum(state) 的正整数拆分
# !50 的整数拆分大约在 2e5的量级

# dp[state] = 1 + (dp[state1]  + dp[state2]  + ... + dp[statek]) / k
# !注意如果nextState和state一样的话 需要合并同类项dp[state]


from functools import lru_cache
from itertools import accumulate
from typing import List, Tuple


class Solution:
    def chipGame(self, nums: List[int], kind: int) -> float:
        """
        注意要正向dp,counter元组保存每种筹码盘子的个数
        放筹码从无到有 不能转化成吃寿司从有到无
        他们判断原地踏步状态的条件是反的 一个是吃0的盘子 一个是超出范围
        """

        @lru_cache(None)
        def dfs(cur: Tuple[int, ...], remain: int) -> float:
            """cur表示有cur[i]个盘子有i个筹码 remain表示还要选的筹码个数"""
            if remain == 0:
                return 0

            lis = list(cur)
            curSum = 0
            step, p = 1, 0
            # !选中哪种盘子的筹码 注意前缀和恰好等于target前缀和时 不能继续向右搬运(再搬运右边就超出限制了)
            for i, count in enumerate(cur):
                curSum += count
                if count == 0 or i + 1 >= len(lis) or curSum == preSum[i + 1]:
                    continue
                lis[i] -= 1
                lis[i + 1] += 1
                step += (dfs(tuple(lis), remain - 1) / kind) * count
                p += count / kind
                lis[i] += 1
                lis[i + 1] -= 1
            return step / p

        max_ = max(nums)
        target = [0] * (max_ + 1)  # 结束时的筹码盘子状态
        for num in nums:
            target[num] += 1
        target[0] += kind - len(nums)
        preSum = [0] + list(accumulate(target))

        init = [0] * (max_ + 1)
        init[0] = kind  # 开始时的筹码盘子状态 要向右搬运到target状态
        return dfs(tuple(init), sum(nums))

    def chipGame1(self, nums: List[int], kind: int) -> float:
        """7500ms"""

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

                # !把 **all** 换成 **for ... break ... else ...** 后快了 1000多ms
                for j in range(kind):
                    if nextState[j] > target[j]:
                        break
                else:
                    step += dfs(nextState, remain - 1) / kind
                    p += 1 / kind

                lis[i] -= 1

            return step / p

        target = sorted(nums + [0] * (kind - len(nums)))
        res = dfs(tuple([0] * kind), sum(nums))  # 正着推导 补齐0
        dfs.cache_clear()
        return res


print(Solution().chipGame(nums=[1, 1], kind=2))  # 3
print(Solution().chipGame(nums=[1, 2], kind=4))  # 3.833333
