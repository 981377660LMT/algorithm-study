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
        """背包 O(V*m)
        dp[index][value] 表示前index种配料选取的总价值能否达到value
        最后加上基料进行验证
        """

        dp = set([0])
        for topping in toppingCosts:
            ndp = set()
            for pre in dp:
                ndp |= {pre, pre + topping, pre + topping * 2}
            dp = ndp

        res = INF
        for base in baseCosts:
            for num in dp:
                cost = base + num
                # 找到最接近target的值
                if (
                    abs(cost - target) < abs(res - target)
                    or abs(cost - target) == abs(res - target)
                    and cost < res
                ):
                    res = cost
        return res


print(Solution().closestCost(baseCosts=[1, 7], toppingCosts=[3, 4], target=10))
