"""字典序遍历十叉树"""

from typing import Generator, List


class Solution:
    def lexicalOrder(self, n: int) -> List[int]:
        """
        返回 [1..n] 的字典序列表。
        时间复杂度 O(n)，额外空间 O(1)（忽略返回结果占用）。
        """
        res: List[int] = [0] * n
        cur = 1
        for i in range(n):
            res[i] = cur
            # 尝试深入到下一层：在当前前缀后面加一个 '0'
            if cur * 10 <= n:
                cur *= 10
            else:
                # 向上回溯，直到能 +1 并且末尾不为 9
                while cur % 10 == 9 or cur + 1 > n:
                    cur //= 10
                cur += 1
        return res

    def lexicalOrder2(self, n: int) -> List[int]:
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

        return list(dfs(0))  # 0表示虚拟根节点


print(Solution().lexicalOrder(13))
