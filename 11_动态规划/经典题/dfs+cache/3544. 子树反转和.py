# 3544. 子树反转和
# https://leetcode.cn/problems/subtree-inversion-sum/description/
#
# O(n) 解法 https://leetcode.cn/problems/subtree-inversion-sum/solutions/3673852/shu-xing-dppythonjavacgo-by-endlesscheng-pjwg/


from collections import deque
from functools import lru_cache
from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:

    def subtreeInversionSum(self, edges: List[List[int]], nums: List[int], k: int) -> int:
        """O(nk)解法."""

        def toRootedTree(tree: List[List[int]], root=0) -> List[List[int]]:
            n = len(tree)
            res = [[] for _ in range(n)]
            visited = [False] * n
            visited[root] = True
            queue = deque([root])
            while queue:
                cur = queue.popleft()
                for next in tree[cur]:
                    if not visited[next]:
                        visited[next] = True
                        queue.append(next)
                        res[cur].append(next)
            return res

        n = len(nums)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
        adjList = toRootedTree(adjList)

        @lru_cache(None)
        def dfs(cur: int, cd: int, mul: int) -> int:
            # 不反转
            res = nums[cur] * mul
            for next in adjList[cur]:
                res += dfs(next, cd - 1 if cd else 0, mul)

            # 反转
            if cd == 0:
                mul *= -1
                res2 = nums[cur] * mul
                for next in adjList[cur]:
                    res2 += dfs(next, k - 1, mul)
                res = max2(res, res2)

            return res

        res = dfs(0, 0, 1)
        dfs.cache_clear()
        return res
