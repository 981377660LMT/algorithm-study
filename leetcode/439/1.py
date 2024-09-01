from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个字符串 coordinate1 和 coordinate2，代表 8 x 8 国际象棋棋盘上的两个方格的坐标。


# 以下是棋盘的参考图。
# 如果这两个方格颜色相同，返回 true，否则返回 false。


# 坐标总是表示有效的棋盘方格。坐标的格式总是先字母（表示列），再数字（表示行）。
class Solution:
    def checkTwoChessboards(self, coordinate1: str, coordinate2: str) -> bool:
        def check(coordinate: str) -> bool:
            return (ord(coordinate[0]) - ord("a") + int(coordinate[1])) % 2 == 0

        return check(coordinate1) == check(coordinate2)
