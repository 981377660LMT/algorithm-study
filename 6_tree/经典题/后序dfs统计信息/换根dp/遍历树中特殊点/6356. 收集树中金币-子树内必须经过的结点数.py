# https://leetcode.cn/problems/collect-coins-in-a-tree/
# 给你一个长度为 n 的数组 coins ，其中 coins[i] 可能为 0 也可能为 1 ，1 表示节点 i 处有一个金币。
# 一开始，你需要选择树中任意一个节点出发。你可以执行下述操作任意次：
#  - 收集距离当前节点距离为 2 以内的所有金币，或者
#  - 移动到树中一个相邻节点。
# !你需要收集树中所有的金币，并且`回到出发节点`，请你返回最少经过的边数。
# 如果你多次经过一条边，每一次经过都会给答案加一。

# !注意回到出发点:2*边权和; 如果不回到出发点,那么就是 `2*边权和-距离coin最远的点的距离``
# !换根dp:求出每个点为根时子树内的必经结点个数,乘以二就是答案(边数)
# 1.每个点维护 (距离最远的一个coin的距离,子树内必须要抵达的点的个数) 元组;
# 2.转移贡献时，如果子树内最远距离 >=2 ，那么说明这个点是必经点，向上返回 (最远距离 + 1, 必经点 + 1)，
#   否则返回 (最远距离 + 1, 必经点)。
# 广义的`距离`


from collections import defaultdict
from typing import List, Tuple
from Rerooting import Rerooting


INF = int(1e18)


class Solution:
    def collectTheCoins(self, coins: List[int], edges: List[List[int]]) -> int:
        E = Tuple[int, int]  # (maxDist, weightSum)

        def e(root: int) -> E:
            return (-INF, 0) if coins[root] == 0 else (0, 0)

        def op(childRes1: E, childRes2: E) -> E:
            dist1, wSum1 = childRes1
            dist2, wSum2 = childRes2
            return (max(dist1, dist2), wSum1 + wSum2)

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            dist, wSum = fromRes
            w = weights[cur][parent] if direction == 0 else weights[parent][cur]
            # dist>=2 时这条边才开始算入必须经过的路径
            return (dist + w, wSum + 1) if dist >= 2 else (dist + w, wSum)

        n = len(coins)
        R = Rerooting(n)
        weights = [defaultdict(int) for _ in range(n)]
        for u, v in edges:
            R.addEdge(u, v)
            weights[u][v] = 1
            weights[v][u] = 1

        dp = R.rerooting(e=e, op=op, composition=composition, root=0)
        res = min(wSum * 2 for _, wSum in dp)  # !如果不回到出发点,那么就是 `2*边权和-距离coin最远的点的距离``
        return res


assert Solution().collectTheCoins([1, 0, 0, 0, 0, 1], [[0, 1], [1, 2], [2, 3], [3, 4], [4, 5]]) == 2
