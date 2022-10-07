from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


# 1 <= shape.length <= 50
# 1 <= shape[i].length <= 50
# shape[i][j] 仅为 'l'、'r' 或 '.'
# 'l'表示向左倾斜的隔板（即从左上到右下）；
# 'r'表示向右倾斜的隔板（即从左下到右上）；
# '.' 表示此位置没有隔板

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def reservoir(self, shape: List[str]) -> int:
        ROW, COL = len(shape), len(shape[0])
