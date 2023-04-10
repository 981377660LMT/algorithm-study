#  有一只跳蚤的家在数轴上的位置 x 处。请你帮助它从位置 0 出发，到达它的家。
#  它可以 往左 跳恰好 leftJump 个位置。
#  它可以 往右 跳恰好 rightJump 个位置。
#  它不能 连续 往左跳 2 次。
#  它不能跳到任何 forbidden 数组中的位置。
#  跳蚤可以往前跳 超过 它的家的位置，但是它 不能跳到负整数 的位置。
#  @param bad
#  @param rightJump  rightJump <= 2000
#  @param leftJump  leftJump <= 2000
#  !@param target  家 0 <= target <= 2000
#  @description 注意bfs范围剪枝
# !搜索范围不会超过 6000

from collections import deque
from typing import List


class Solution:
    def minimumJumps(self, forbidden: List[int], rightJump: int, leftJump: int, target: int) -> int:
        visited = set(forbidden) | {0}
        upper = 6000  # 上界,可以取大一点6000
        queue = deque([(0, 0, 0)])  # (pos, step, leftCount)
        while queue:
            cur, step, leftCount = queue.popleft()
            if cur == target:
                return step
            if leftCount == 0:
                next1 = cur - leftJump
                if next1 >= 0 and next1 not in visited:
                    visited.add(next1)
                    queue.append((next1, step + 1, leftCount + 1))
            next2 = cur + rightJump
            if next2 <= upper and next2 not in visited:  # !剪枝 :出右边界后没必要继续走
                visited.add(next2)
                queue.append((next2, step + 1, 0))
        return -1
