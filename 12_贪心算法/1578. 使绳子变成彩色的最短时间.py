# 不希望两个连续的气球涂着相同的颜色
# 返回 Bob 使绳子变成 彩色 需要的 最少时间 。
from typing import List


class Solution:
    def minCost(self, colors: str, neededTime: List[int]) -> int:
        res = 0
        index = 0

        while index < len(colors):
            curColor = colors[index]
            cost = 0
            maxNeed = 0

            # 相同的里面留下哪个最大的
            while index < len(colors) and colors[index] == curColor:
                cost += neededTime[index]
                maxNeed = max(maxNeed, neededTime[index])
                index += 1

            res += cost - maxNeed

        return res


print(Solution().minCost(colors="aabaa", neededTime=[1, 2, 3, 4, 1]))
