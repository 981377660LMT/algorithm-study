from typing import List
from sortedcontainers import SortedList


class Solution:
    def getLengthOfWaterfallFlow(self, num: int, block: List[int]) -> int:
        sl = SortedList([(0, i) for i in range(num)])  # 可以用堆代替
        for b in block:
            preH, preI = sl.pop(0)
            sl.add((preH + b, preI))
        return sl[-1][0]
