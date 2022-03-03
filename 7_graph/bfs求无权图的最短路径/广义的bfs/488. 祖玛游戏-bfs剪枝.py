# 1 <= board.length <= 16
# 1 <= hand.length <= 5
from collections import deque
from functools import lru_cache
import re


# 为什么使用广度优先搜索？
# 因为只需要找出需要回合数最少的方案，因此使用广度优先搜索可以得到可以消除桌面上所有球的方案时就直接返回结果，而不需要继续遍历更多需要回合数更多的方案。
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
            for i in range(len(b)):
                # 删除那个元素
                for j in range(len(h)):
                    # 最重要的剪枝是，当手上的球 h[j] 和插入位置 i 前后的球 b[i-1], b[i] 三个球各不相同时，插入是不必要的：
                    sequence = [b[i - 1], b[i], h[j]] if i else [b[i], h[j]]
                    if len(set(sequence)) < len(sequence):
                        nextB = clean(b[:i] + h[j] + b[i:])
                        nextH = h[:j] + h[j + 1 :]
                        if (nextB, nextH) not in visited:
                            visited.add((nextB, nextH))
                            queue.append((nextB, nextH, step + 1))
        return -1


print(Solution().findMinStep(board="WRRBBW", hand="RB"))
print(Solution().findMinStep(board="WWRRBBWW", hand="WRBRW"))

# re.subn返回一个元组
