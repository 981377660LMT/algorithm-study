# 网页瀑布流
# 当有数据块需要加载时，优先加载在高度最短的那一列；
# 若存在多个高度相同且最短的情况，则加载在其中最靠左的那一列
# 已知当前网页共分割为 num 列，该网页有若干数据块可以加载，block[i] 表示第 i 个数据块的高度。
# 当页面按顺序加载完所有的数据块后，请返回高度最大的那一列的高度。

from typing import List
from sortedcontainers import SortedList


class Solution:
    def getLengthOfWaterfallFlow(self, num: int, block: List[int]) -> int:
        sl = SortedList([(0, i) for i in range(num)])
        for b in block:
            preH, preI = sl.pop(0)
            sl.add((preH + b, preI))
        return sl[-1][0]
