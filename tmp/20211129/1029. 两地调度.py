from typing import List

# 返回将每个人都飞到 a 、b 中某座城市的最低费用，要求每个城市都有 n 人抵达。

# 我们将每位顾客按照飞往A的金额-飞往B的金额差值来进行排序。
# 排在前面的肯定应该飞往A,而排在后面的人就理所应当的要飞往B了


class Solution:
    def twoCitySchedCost(self, costs: List[List[int]]) -> int:
        costs.sort(key=lambda pair: pair[0] - pair[1])
        return sum(v[0] if i < len(costs) / 2 else v[1] for i, v in enumerate(costs))


print(Solution().twoCitySchedCost(costs=[[10, 20], [30, 200], [400, 50], [30, 20]]))
