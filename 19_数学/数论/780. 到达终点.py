import sys

sys.setrecursionlimit(99999999)
from functools import lru_cache

# 从点 (x, y) 可以转换到 (x, x+y)  或者 (x+y, y)。

# It looks like a tree, if you know one node, it's always easy to get the root because you have exact one way to get the parent node.
# https://leetcode-cn.com/problems/reaching-points/solution/python-cong-zhong-dian-wang-qi-dian-dao-zhao-tui-b/
class Solution:
    def reachingPoints(self, sx: int, sy: int, tx: int, ty: int) -> bool:
        while sx < tx and sy < ty:
            tx, ty = tx % ty, ty % tx
        return (
            sx == tx
            and sy <= ty
            and (ty - sy) % sx == 0
            or sy == ty
            and sx <= tx
            and (tx - sx) % sy == 0
        )

    # @lru_cache(typed=False, maxsize=128000000)
    # def reachingPoints(self, sx: int, sy: int, tx: int, ty: int) -> bool:
    #     if tx == sx and ty == sy:
    #         return True

    #     if tx < sx or ty < sy:
    #         return False

    #     if tx == sx:
    #         return (ty - sy) % sx == 0
    #     if ty == sy:
    #         return (tx - sx) % sy == 0

    #     return self.reachingPoints(sx, sy, tx, ty - tx) or self.reachingPoints(sx, sy, tx - ty, ty)


print(Solution().reachingPoints(sx=1, sy=1, tx=3, ty=5))
# 输出: True
# 解释:
# 可以通过以下一系列转换从起点转换到终点：
# (1, 1) -> (1, 2)
# (1, 2) -> (3, 2)
# (3, 2) -> (3, 5)
