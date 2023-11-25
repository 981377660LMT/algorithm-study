from typing import List, Set, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个网格图，由 n + 2 条 横线段 和 m + 2 条 竖线段 组成，一开始所有区域均为 1 x 1 的单元格。

# 所有线段的编号从 1 开始。

# 给你两个整数 n 和 m 。

# 同时给你两个整数数组 hBars 和 vBars 。

# hBars 包含区间 [2, n + 1] 内 互不相同 的横线段编号。
# vBars 包含 [2, m + 1] 内 互不相同的 竖线段编号。
# 如果满足以下条件之一，你可以 移除 两个数组中的部分线段：


# 如果移除的是横线段，它必须是 hBars 中的值。
# 如果移除的是竖线段，它必须是 vBars 中的值。
# 请你返回移除一些线段后（可能不移除任何线段），剩余网格图中 最大正方形 空洞的面积，正方形空洞的意思是正方形 内部 不含有任何线段。


def max2(a, b):
    return a if a > b else b


class Solution:
    def maximizeSquareHoleArea(self, n: int, m: int, hBars: List[int], vBars: List[int]) -> int:
        hBars.sort()
        vBars.sort()

        arr1 = list(set(sorted([0] + [v - 1 for v in hBars] + [n + 1])))
        arr2 = list(set(sorted([0] + [v - 1 for v in vBars] + [m + 1])))

        res = 0
        # 枚举左上角
        for i1 in range(len(arr1)):
            for j1 in range(len(arr2)):
                curJ = j1
                curX, curY = 0, 0
                for i2 in range(i1 + 1, len(arr1)):
                    curX += arr1[i2] - arr1[i2 - 1]
                    while curJ + 1 < len(arr2) and curY < curX:
                        curJ += 1
                        curY += arr2[curJ] - arr2[curJ - 1]
                    if curX == curY:
                        res = max2(res, curX * curY)
                        print(arr1[i1], arr2[j1], arr1[i2], arr2[curJ], res)

        return res


# 2
# 4
# [3,2]
# [4,2]

print(Solution().maximizeSquareHoleArea(2, 4, [3, 2], [4, 2]))
