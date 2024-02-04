from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个  n x 2 的二维数组 points ，它表示二维平面上的一些点坐标，其中 points[i] = [xi, yi] 。

# 我们定义 x 轴的正方向为 右 （x 轴递增的方向），x 轴的负方向为 左 （x 轴递减的方向）。类似的，我们定义 y 轴的正方向为 上 （y 轴递增的方向），y 轴的负方向为 下 （y 轴递减的方向）。

# 你需要安排这 n 个人的站位，这 n 个人中包括 liupengsay 和小羊肖恩 。你需要确保每个点处 恰好 有 一个 人。同时，liupengsay 想跟小羊肖恩单独玩耍，所以 liupengsay 会以 liupengsay 的坐标为 左上角 ，小羊肖恩的坐标为 右下角 建立一个矩形的围栏（注意，围栏可能 不 包含任何区域，也就是说围栏可能是一条线段）。如果围栏的 内部 或者 边缘 上有任何其他人，liupengsay 都会难过。

# 请你在确保 liupengsay 不会 难过的前提下，返回 liupengsay 和小羊肖恩可以选择的 点对 数目。

# 注意，liupengsay 建立的围栏必须确保 liupengsay 的位置是矩形的左上角，小羊肖恩的位置是矩形的右下角。比方说，以 (1, 1) ，(1, 3) ，(3, 1) 和 (3, 3) 为矩形的四个角，给定下图的两个输入，liupengsay 都不能建立围栏，原因如下：


# 图一中，liupengsay 在 (3, 3) 且小羊肖恩在 (1, 1) ，liupengsay 的位置不是左上角且小羊肖恩的位置不是右下角。
# 图二中，liupengsay 在 (1, 3) 且小羊肖恩在 (1, 1) ，小羊肖恩的位置不是在围栏的右下角。


class Solution:
    def numberOfPairs(self, points: List[List[int]]) -> int:
        ...
