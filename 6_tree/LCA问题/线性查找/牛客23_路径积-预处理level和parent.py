# 给定一棵n个节点的无根树(n个结点，n-1条边的无环连通图)，每个节点有一个权值ai
# 一共有m次查询，每次查询xi到yi的最短路径上所有点权的乘积。
# 为了防止答案过大，答案对1e9+7取模。

from collections import defaultdict
from typing import List


MOD = int(1e9 + 7)


class Solution:
    def solve(
        self,
        n: int,
        m: int,
        weights: List[int],
        u: List[int],
        v: List[int],
        queryStarts: List[int],
        queryEnds: List[int],
    ) -> List[int]:
        # dfs 处理深度level、parent
        # 暴力O(n)：每次选择深度较大的一个点往上跳，直到跳到跟另外一个点深度相同。接着两个点同时向上跳，直到两个点相遇，即找到最近公共祖先。把跳的过程中遇到的点的权值乘到答案里即可
        # 还可以树上倍增优化到log(n)
        def dfs(cur: int, parent: int, level: int) -> None:
            parentMap[cur] = parent
            levelMap[cur] = level
            for next in adjMap[cur]:
                if next == parent:
                    continue
                dfs(next, cur, level + 1)

        adjMap = defaultdict(list)
        for cur, next in zip(u, v):
            adjMap[cur].append(next)
            adjMap[next].append(cur)

        levelMap, parentMap = defaultdict(lambda: -1), defaultdict(lambda: -1)
        dfs(1, -1, 0)

        res = []
        for r1, r2 in zip(queryStarts, queryEnds):
            if levelMap[r1] < levelMap[r2]:
                r1, r2 = r2, r1
            cur = 1
            while levelMap[r1] > levelMap[r2]:
                cur *= weights[r1 - 1]
                cur %= MOD
                r1 = parentMap[r1]
            while r1 != r2:
                cur *= weights[r1 - 1]
                cur %= MOD
                r1 = parentMap[r1]
                cur *= weights[r2 - 1]
                cur %= MOD
                r2 = parentMap[r2]
            cur *= weights[r1 - 1]
            res.append(cur % MOD)

        return res


print(Solution().solve(2, 2, [1001, 3357], [1], [2], [2, 1], [1, 2]))
