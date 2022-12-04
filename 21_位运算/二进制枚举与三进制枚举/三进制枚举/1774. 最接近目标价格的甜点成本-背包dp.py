# 目前共有 n 种冰激凌基料和 m 种配料可供选购。而制作甜点需要遵循以下几条规则：

# 必须选择 一种 冰激凌基料。
# 可以添加 一种或多种 配料，也可以不添加任何配料。
# 每种类型的配料 最多两份 。

# 返回最接近 target 的甜点成本。如果有多种方案，返回 成本相对较低 的一种。
# n,m<=10
# 1 <= baseCosts[i], toppingCosts[i] <= 1e4
# 1 <= target <= 1e4

INF = int(1e18)


class Solution:
    def closestCost(self, baseCosts: list[int], toppingCosts: list[int], target: int) -> int:
        """dfs代替三进制枚举会更快 并且更好剪枝"""

        def dfs(index: int, curCost: int) -> None:
            nonlocal res
            if index == len(toppingCosts):
                if abs(curCost - target) < abs(res - target) or (
                    abs(curCost - target) == abs(res - target) and curCost < res
                ):
                    res = curCost
                return
            dfs(index + 1, curCost)
            dfs(index + 1, curCost + toppingCosts[index])
            dfs(index + 1, curCost + toppingCosts[index] * 2)

        res = INF
        for baseCost in baseCosts:
            dfs(0, baseCost)
        return res


print(Solution().closestCost(baseCosts=[1, 7], toppingCosts=[3, 4], target=10))
