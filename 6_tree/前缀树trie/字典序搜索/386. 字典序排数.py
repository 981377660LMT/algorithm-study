"""字典序遍历十叉树"""

from collections import deque
from typing import Generator, List


class Solution:
    def lexicalOrder(self, n: int) -> List[int]:
        """
        给你一个整数 n ，按字典序返回范围 [1, n] 内所有整数。
        你必须设计一个时间复杂度为 O(n) 且使用 O(1) 额外空间的算法。
        """

        def dfs(cur: int) -> Generator[int, None, None]:
            for i in range(10):
                next = cur * 10 + i
                if next > n:
                    return
                if next == 0:
                    continue
                yield next
                yield from dfs(next)

        # 0表示虚拟根节点
        return list(dfs(0))


print(Solution().lexicalOrder(13))

#####################################################################
