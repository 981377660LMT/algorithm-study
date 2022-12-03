from typing import List
from functools import lru_cache

# 目前共有 n 种冰激凌基料和 m 种配料可供选购。而制作甜点需要遵循以下几条规则：

# 必须选择 一种 冰激凌基料。
# 可以添加 一种或多种 配料，也可以不添加任何配料。
# 每种类型的配料 最多两份 。

# 你希望自己做的甜点总成本尽可能接近目标价格 target 。
# 返回最接近 target 的甜点成本。`如果有多种方案，返回 成本相对较低 的一种。`(元组排序法)

# 1 <= n, m <= 10
# Time: O(n * 2^m).


class Solution:
    def closestCost(self, baseCosts: List[int], toppingCosts: List[int], target: int) -> int:
        # 可以用两次
        toppingCosts *= 2
        self.res = baseCosts[0]

        @lru_cache(None)
        def dfs(index: int, cost: int) -> None:
            if cost > target:
                if cost - target < abs(self.res - target):
                    self.res = cost
                return
            if target - cost <= abs(self.res - target):
                self.res = cost
            if index == len(toppingCosts):
                return

            dfs(index + 1, cost + toppingCosts[index])
            dfs(index + 1, cost)

        for base in baseCosts:
            dfs(0, base)

        return self.res


print(Solution().closestCost(baseCosts=[1, 7], toppingCosts=[3, 4], target=10))
# 输出：10
# 解释：考虑下面的方案组合（所有下标均从 0 开始）：
# - 选择 1 号基料：成本 7
# - 选择 1 份 0 号配料：成本 1 x 3 = 3
# - 选择 0 份 1 号配料：成本 0 x 4 = 0
# 总成本：7 + 3 + 0 = 10 。
