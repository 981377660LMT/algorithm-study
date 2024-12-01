# 3373. 连接两棵树后最大目标节点数目 II (黑白染色)
# https://leetcode.cn/problems/maximize-the-number-of-target-nodes-after-connecting-trees-ii/description/
#
# !如果节点 u 和节点 v 之间路径的边数是偶数，那么我们称节点 u 是节点 v 的 目标节点 。
# 注意 ，一个节点一定是它自己的 目标节点 。
#
# !请你返回一个长度为 n 的整数数组 answer ，answer[i] 表示将第一棵树中的一个节点与第二棵树中的一个节点连接一条边后，
# !第一棵树中节点 i 的 目标节点 数目的 最大值 。
#
# 注意 ，每个查询相互独立。意味着进行下一次查询之前，你需要先把刚添加的边给删掉。
#
# 2 <= n,m <= 1e5


from typing import List, Tuple


class Solution:
    def maxTargetNodes(self, edges1: List[List[int]], edges2: List[List[int]]) -> List[int]:
        n, m = len(edges1) + 1, len(edges2) + 1
        tree1, tree2 = [[] for _ in range(n)], [[] for _ in range(m)]
        for u, v in edges1:
            tree1[u].append(v)
            tree1[v].append(u)
        for u, v in edges2:
            tree2[u].append(v)
            tree2[v].append(u)

        def calc(tree: List[List[int]]) -> Tuple[List[int], int, int]:
            """
            返回根节点到各个节点的距离，距离根节点为偶数的节点数，距离根节点为奇数的节点数.
            """
            dist = [-1] * len(tree)
            even, odd = 0, 0
            stack = [0]
            dist[0] = 0
            while stack:
                cur = stack.pop()
                if dist[cur] % 2 == 0:
                    even += 1
                else:
                    odd += 1
                for next in tree[cur]:
                    if dist[next] == -1:
                        dist[next] = dist[cur] + 1
                        stack.append(next)
            return dist, even, odd

        dist1, c1Even, c1Odd = calc(tree1)
        _, c2Even, c2Odd = calc(tree2)
        c2Max = max(c2Even, c2Odd)
        res = []
        for d in dist1:
            target = c1Even if d % 2 == 0 else c1Odd
            res.append(target + c2Max)
        return res
