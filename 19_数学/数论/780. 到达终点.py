# 从点 (x, y) 可以转换到 (x, x+y)  或者 (x+y, y)。

# It looks like a tree, if you know one node, it's always easy to get the root because you have exact one way to get the parent node.
# https://leetcode-cn.com/problems/reaching-points/solution/python-cong-zhong-dian-wang-qi-dian-dao-zhao-tui-b/

# 正难则反
class Solution:
    def reachingPoints(self, sx: int, sy: int, tx: int, ty: int) -> bool:
        if sx > tx or sy > ty:
            return False
        # 从target反向到start 加法可以用模加速
        while tx > sx and ty > sy:
            if tx > ty:
                tx %= ty
            else:
                ty %= tx
        if sx == tx:
            return sy <= ty and (ty - sy) % sx == 0
        if sy == ty:
            return sx <= tx and (tx - sx) % sy == 0
        return False


print(Solution().reachingPoints(sx=1, sy=1, tx=3, ty=5))
# 输出: True
# 解释:
# 可以通过以下一系列转换从起点转换到终点：
# (1, 1) -> (1, 2)
# (1, 2) -> (3, 2)
# (3, 2) -> (3, 5)
