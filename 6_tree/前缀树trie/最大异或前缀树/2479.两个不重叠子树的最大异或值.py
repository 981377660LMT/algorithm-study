from typing import List
from XORTrieArray import useXORTrie

# !两个不重叠子树的`子树和`的最大异或值

# !如何保证子树不重叠?
# !前序遍历时查询,后序遍历时插入


class Solution:
    def maxXor(self, n: int, edges: List[List[int]], values: List[int]) -> int:
        def dfs1(cur: int, pre: int) -> int:
            """预处理子树和"""
            res = values[cur]
            for next in adjList[cur]:
                if next == pre:
                    continue
                res += dfs1(next, cur)
            subSum[cur] = res
            return res

        def dfs2(cur: int, pre: int) -> None:
            """查询两个不重叠子树的最大异或值"""
            nonlocal res
            res = max(res, xorTrie.search(subSum[cur]))
            for next in adjList[cur]:
                if next == pre:
                    continue
                dfs2(next, cur)
            xorTrie.insert(subSum[cur])

        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0
        subSum = [0] * n
        xorTrie = useXORTrie(int(1e14))
        dfs1(0, -1)
        dfs2(0, -1)
        return res


print(
    Solution().maxXor(
        n=6, edges=[[0, 1], [0, 2], [1, 3], [1, 4], [2, 5]], values=[2, 8, 3, 6, 2, 5]
    )
)
print(Solution().maxXor(n=3, edges=[[0, 1], [1, 2]], values=[4, 6, 1]))
