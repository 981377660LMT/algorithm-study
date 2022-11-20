# 给你 n 个项目。对于每个项目 i ，它都有一个纯利润 profits[i] ，和启动该项目需要的最小资本 capital[i] 。
# 最初，你的资本为 w 。当你完成一个项目时，你将获得纯利润，且利润将被添加到你的总资本中。
# 总而言之，从给定项目中选择 最多 k 个不同项目的列表，以 最大化最终资本 ，并输出最终可获得的最多资本。
# IPO


from typing import List
from sortedcontainers import SortedList


class Solution:
    def findMaximizedCapital(self, k: int, w: int, profits: List[int], capital: List[int]) -> int:
        sl1 = SortedList((c, p) for c, p in zip(capital, profits))  # 启动成本升序
        sl2 = SortedList(key=lambda x: -x)  # 利润降序

        curCapital = w
        for _ in range(k):
            while sl1 and sl1[0][0] <= curCapital:
                sl2.add(sl1.pop(0)[1])
            if not sl2:
                break
            curCapital += sl2.pop(0)
        return curCapital


print(Solution().findMaximizedCapital(2, 0, [1, 2, 3], [0, 1, 1]))
