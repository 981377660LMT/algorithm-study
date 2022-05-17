from typing import List
from sortedcontainers import SortedList


class Solution:
    def getSkyline(self, buildings: List[List[int]]) -> List[List[int]]:
        events = []
        # 横坐标相同时，左端点进在前，右端点出在后
        for left, right, height in buildings:
            events.append((left, -height, 0))
            events.append((right, height, 1))
        events.sort()

        preY = 0
        sl = SortedList([preY])  # 最底下的点，哨兵元素
        res = []
        for x, y, flag in events:
            if flag == 0:
                sl.add(-y)
            else:
                sl.discard(y)

            top = sl[-1]
            if top != preY:
                res.append([x, top])

            preY = top

        return res

