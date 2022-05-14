from typing import List

# 0 <= xi, yi < 10^6
# 障碍的最大数量小于等于200。
# 范围很大，如果直接bfs。不剪枝就会TLE
# 每次移动，都可以走到网格中在四个方向上相邻的方格，只要该方格 不 在给出的封锁列表 blocked 上。同时，不允许走出网格。
# 只有在可以通过一系列的移动从源方格 source 到达目标方格 target 时才返回 true。


# 总结：
# 障碍的最大数量小于等于200 所以如果bfs层数大于障碍数200 则说明障碍没有包围这个点
# we can just use bfs to search from the source, and set a maximum step.
# after moving maximum step, if we can still move, then it must can reach the target point
# https://leetcode.com/problems/escape-a-large-maze/discuss/282870/python-solution-with-picture-show-my-thoughts
class Solution:
    def isEscapePossible(
        self, blocked: List[List[int]], source: List[int], target: List[int]
    ) -> bool:
        blockedSet = set(tuple(p) for p in blocked)

        def bfs(x, y, tx, ty):
            """Return True if (x, y) is not looped from (tx, ty)."""
            visited = {(x, y)}
            queue = [(x, y)]
            level = 0
            while queue:
                level += 1
                if level > len(blockedSet):  # 剪枝
                    return True
                nextQueue = []
                for x, y in queue:
                    if (x, y) == (tx, ty):
                        return True
                    for xx, yy in (x - 1, y), (x, y - 1), (x, y + 1), (x + 1, y):
                        if (
                            0 <= xx < 1e6
                            and 0 <= yy < 1e6
                            and (xx, yy) not in blockedSet
                            and (xx, yy) not in visited
                        ):
                            visited.add((xx, yy))
                            nextQueue.append((xx, yy))
                queue = nextQueue
            return False

        return bfs(*source, *target) and bfs(*target, *source)


print(Solution().isEscapePossible(blocked=[[0, 1], [1, 0]], source=[0, 0], target=[0, 2]))
# 输出：false
# 解释：
# 从源方格无法到达目标方格，因为我们无法在网格中移动。
# 无法向北或者向东移动是因为方格禁止通行。
# 无法向南或者向西移动是因为不能走出网格。

