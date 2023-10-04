# 返回将 instructions 中所有元素依次插入 nums 后的 总最小代价
# 每一次插入操作的 代价 是以下两者的 较小值 ：
# nums 中 严格小于  instructions[i] 的数字数目。
# nums 中 严格大于  instructions[i] 的数字数目。


from typing import List
from sortedcontainers import SortedList


class Solution:
    def createSortedArray(self, instructions: List[int]) -> int:
        res = 0
        sl = SortedList()
        for num in instructions:
            smaller = sl.bisect_left(num)
            bigger = len(sl) - sl.bisect_right(num)
            res += min(smaller, bigger)
            sl.add(num)
        return res % (10**9 + 7)
