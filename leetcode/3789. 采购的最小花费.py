# 3789. 采购的最小花费
# 给你五个整数 cost1, cost2, costBoth, need1 和 need2。
#
# 有三种类型的物品可以购买：
#
# 类型 1 的物品花费 cost1，并仅满足类型 1 的需求 1 个单位。
# 类型 2 的物品花费 cost2，并仅满足类型 2 的需求 1 个单位。
# 类型 3 的物品花费 costBoth，同时满足类型 1 和类型 2 的需求各 1 个单位。
# 你需要购买足够的物品，使得：
#
# 满足类型 1 的总需求 至少 为 need1。
# 满足类型 2 的总需求 至少 为 need2。
# 返回满足这些需求的 最小 可能总花费。


class Solution:
    def minimumCost(self, cost1: int, cost2: int, costBoth: int, need1: int, need2: int) -> int:
        min_ = min(need1, need2)
        return (
            min_ * min(costBoth, cost1 + cost2)
            + (need1 - min_) * min(cost1, costBoth)
            + (need2 - min_) * min(cost2, costBoth)
        )
