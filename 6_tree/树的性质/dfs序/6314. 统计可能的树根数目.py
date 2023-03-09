from collections import defaultdict
from itertools import accumulate
from typing import List
from DFSOrder import DFSOrder

# Alice 有一棵 n 个节点的树，节点编号为 0 到 n - 1 。
# 树用一个长度为 n - 1 的二维整数数组 edges 表示，其中 edges[i] = [ai, bi] ，表示树中节点 ai 和 bi 之间有一条边。
# Alice 想要 Bob 找到这棵树的根。她允许 Bob 对这棵树进行若干次 猜测 。每一次猜测，Bob 做如下事情：
# 选择两个 不相等 的整数 u 和 v ，且树中必须存在边 [u, v] 。
# Bob 猜测树中 u 是 v 的 父节点 。
# !Bob 的猜测用二维整数数组 guesses 表示，其中 guesses[j] = [uj, vj] 表示 Bob 猜 uj 是 vj 的`父节点`。
# !Alice 非常懒，她不想逐个回答 Bob 的猜测，只告诉 Bob 这些猜测里面 至少 有 k 个猜测的结果为 true 。
# 给你二维整数数组 edges ，Bob 的所有猜测和整数 k ，请你返回可能成为树根的 节点数目 。如果没有这样的树，则返回 0。


class Solution:
    def rootCount(self, edges: List[List[int]], guesses: List[List[int]], k: int) -> int:
        """对每对查询,看能对哪些区间的根节点产生贡献"""
        n = len(edges) + 1
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
        D = DFSOrder(n, adjList, root=0)
        diff = [0] * (n + 10)
        for u, v in guesses:
            if D.isAncestor(u, v):
                start, end = D.querySub(v)  # [1,start-1] + [end+1,n] 可以作为根节点
                diff[1] += 1
                diff[start] -= 1
                diff[end + 1] += 1
            else:
                start, end = D.querySub(u)  # 子树[start,end]可以作为根节点
                diff[start] += 1
                diff[end + 1] -= 1

        diff = list(accumulate(diff))
        diff = diff[1 : n + 1]
        return sum(x >= k for x in diff)


# [[0,1],[1,2],[1,3],[4,2]]
# [[1,3],[0,1],[1,0],[2,4]]
# 3
print(Solution().rootCount([[0, 1], [1, 2], [1, 3], [4, 2]], [[1, 3], [0, 1], [1, 0], [2, 4]], 3))
