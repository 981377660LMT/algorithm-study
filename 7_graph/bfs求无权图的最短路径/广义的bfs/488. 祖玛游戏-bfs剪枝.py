# 1 <= board.length <= 16
# 1 <= hand.length <= 5
from collections import deque
from functools import lru_cache
import re

# 现在的代码已经TLE了
class Solution:
    def findMinStep(self, board: str, hand: str) -> int:
        """请你按上述操作步骤移除掉桌上所有球，计算并返回所需的 最少 球数。如果不能移除桌上所有的球，返回 -1 。"""

        @lru_cache(None)
        def clean(s: str) -> str:
            """碰到三个就删除整个"""
            count = 1
            while count:
                s, count = re.subn(r'(\w)\1{2,}', '', s)
            return s

        hand = ''.join(sorted(hand))
        queue = deque([(board, hand, 0)])
        visited = set([(board, hand)])

        while queue:
            b, h, step = queue.popleft()
            if not b:
                return step

            # 插入位置
            for i in range(len(b) + 1):
                # 删除那个元素
                for j in range(len(h)):
                    nextB = clean(b[:i] + h[j] + b[i:])
                    nextH = h[:j] + h[j + 1 :]
                    if (nextB, nextH) not in visited:
                        visited.add((nextB, nextH))
                        queue.append((nextB, nextH, step + 1))
        return -1


print(Solution().findMinStep(board="WRRBBW", hand="RB"))
print(Solution().findMinStep(board="WWRRBBWW", hand="WRBRW"))

# re.subn返回一个元组
