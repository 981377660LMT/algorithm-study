class Solution:
    def closestCost(self, baseCosts: list[int], toppingCosts: list[int], target: int) -> int:
        n, m = len(baseCosts), len(toppingCosts)
        res = baseCosts[0]

        for base in baseCosts:
            for state in range(3 ** m):
                cost = base
                for i in range(m):
                    prefix = state // (3 ** i)
                    mod = prefix % 3
                    if mod == 1:
                        cost += toppingCosts[i]
                    elif mod == 2:
                        cost += toppingCosts[i] * 2
                    if cost > target:
                        break

                if abs(cost - target) < abs(res - target) or (
                    abs(cost - target) == abs(res - target) and cost < res
                ):
                    res = cost

        return res


print(Solution().closestCost(baseCosts=[1, 7], toppingCosts=[3, 4], target=10))
