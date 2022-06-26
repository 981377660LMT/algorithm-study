"""
删除树中两条 不同 的边以形成三个连通组件

分别获取三个组件 每个 组件中所有节点值的异或值。
最大 异或值和 最小 异或值的 差值 就是这一种删除边方案的分数。

3 <= n <= 1000
"""

from typing import List
from itertools import combinations
from collections import defaultdict
from DFSOrder import DFSOrder


class Solution:
    def minimumScore(self, nums: List[int], edges: List[List[int]]) -> int:
        """枚举删除的边的`下`端点
        
        parents写错了 不能只记录一个父节点 而是要记录祖先结点
        优化：判断祖先可以用 dfs 序 而不用dfs先序遍历记录祖先节点
        """

        def dfs(cur: int, pre: int) -> None:
            subXor[cur] = nums[cur]
            for next in adjMap[cur]:
                if next == pre:
                    continue
                ancestors[next] |= ancestors[cur] | {cur}  # !先序遍历传递祖先节点
                dfs(next, cur)
                subXor[cur] ^= subXor[next]  # !后序遍历统计异或值

        n = len(nums)
        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        subXor = [0] * n
        ancestors = [set() for _ in range(n)]  # 祖先节点
        dfs(0, -1)

        allXor = subXor[0]
        res = int(1e20)

        # !枚举删除的边的`下`端点 p1 可能是 p2 的祖先节点
        # !注意下端点不能是根节点 不选0
        for p1, p2 in combinations(range(1, n), 2):
            isP1Parent = p1 in ancestors[p2]
            isP2Parent = p2 in ancestors[p1]
            isParent = isP1Parent or isP2Parent
            if isP2Parent:
                p1, p2 = p2, p1

            xor1, xor2, xor3 = 0, 0, 0
            if not isParent:  # 不在同一子树
                xor2, xor1, xor3 = subXor[p2], subXor[p1], allXor ^ subXor[p1] ^ subXor[p2]
            else:  # 在同一子树
                xor2, xor1, xor3 = subXor[p2], subXor[p1] ^ subXor[p2], allXor ^ subXor[p1]

            a, _, c = sorted([xor1, xor2, xor3])
            if (cand := (c - a)) < res:
                res = cand

        return res

    def minimumScore2(self, nums: List[int], edges: List[List[int]]) -> int:
        """枚举删除的边的`下`端点
        优化：判断祖先可以用 dfs 序
        """

        def dfs(cur: int, pre: int) -> None:
            subXor[cur] = nums[cur]
            for next in adjMap[cur]:
                if next == pre:
                    continue
                dfs(next, cur)
                subXor[cur] ^= subXor[next]  # !后序遍历统计异或值

        n = len(nums)
        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        subXor = [0] * n
        dfs(0, -1)

        allXor = subXor[0]
        res = int(1e20)
        D = DFSOrder(n, adjMap)

        for p1, p2 in combinations(range(1, n), 2):
            isP1Parent = D.isAncestor(p1, p2)
            isP2Parent = D.isAncestor(p2, p1)
            isParent = isP1Parent or isP2Parent
            if isP2Parent:
                p1, p2 = p2, p1

            xor1, xor2, xor3 = 0, 0, 0
            if not isParent:  # 不在同一子树
                xor2, xor1, xor3 = subXor[p2], subXor[p1], allXor ^ subXor[p1] ^ subXor[p2]
            else:  # 在同一子树
                xor2, xor1, xor3 = subXor[p2], subXor[p1] ^ subXor[p2], allXor ^ subXor[p1]

            a, _, c = sorted([xor1, xor2, xor3])
            if (cand := (c - a)) < res:
                res = cand

        return res


print(Solution().minimumScore([9, 14, 2, 1], [[2, 3], [3, 0], [3, 1]]))
print(Solution().minimumScore2([9, 14, 2, 1], [[2, 3], [3, 0], [3, 1]]))
