from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# !展示的区域为正三角形，这片区域可以拆分为若干个子区域，每个子区域都是边长为 1 的小三角形，其中第 i 行有 2i - 1 个小三角形。
# !初始情况下，区域中的所有位置都为沙地，你需要指定一些子区域种植沙柳树成为绿地，以达到转化整片区域为绿地的最终目的，规则如下：

# !1.若两个子区域共用一条边，则视为相邻；
# !2.若至少有两片绿地与同一片沙地相邻，则这片沙地也会转化为绿地
# !3.转化为绿地的区域会影响其相邻的沙地

# !现要将一片边长为 size 的沙地全部转化为绿地，
# !请找到任意一种初始指定 最少 数量子区域种植沙柳的方案，并返回所有初始种植沙柳树的绿地坐标。


class Solution:
    def sandyLandManagement(self, size: int) -> List[List[int]]:
        def getEdgePos(size: int, level: int) -> List[List[int]]:
            if level == 1:
                return [[1, 1]]
            return [[level, 1], [level, level * 2 - 1]]

        res = []
        for i in range(1, size + 1):
            res.extend(getEdgePos(size, i))
        return res


print(Solution().sandyLandManagement(3))
print(Solution().sandyLandManagement(2))
print(Solution().sandyLandManagement(1))
