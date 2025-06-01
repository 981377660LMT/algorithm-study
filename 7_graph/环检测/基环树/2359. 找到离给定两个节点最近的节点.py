"""内向基环树森林中 求到两个给定结点最近的节点"""

# 2359. 找到离给定两个节点最近的节点
# https://leetcode.cn/problems/find-closest-node-to-given-two-nodes/solutions/1710829/ji-suan-dao-mei-ge-dian-de-ju-chi-python-gr2u/?envType=daily-question&envId=2025-05-30
#
# 每个节点 至多 有一条出边(内向基环树)。
# 请你返回一个从 node1 和 node2 都能到达节点的编号，
# 使节点 node1 和节点 node2 到这个节点的距离 较大值最小化。
# 如果有多个答案，请返回 最小 的节点编号。如果答案不存在，返回 -1 。
#
#
# 如果输入的不止两个节点 node 1 ​ 和 node 2 ​ ，而是一个很长的 nodes 列表，要怎么做呢？
# 如果输入的是 queries 询问数组，每个询问包含两个节点 node 1 ​ 和 node 2 ​ ，你需要快速计算 closestMeetingNode(edges, node1, node2)，要怎么做呢？
#
# 我吹过你吹过的晚风

from typing import List


class Solution:
    def closestMeetingNode(self, edges: List[int], node1: int, node2: int) -> int:
        """两个结点直接bfs找最近结点"""
        n = len(edges)

        def bfs(start: int) -> List[int]:
            dist = [n] * n
            d = 0
            while start >= 0 and dist[start] == n:
                dist[start] = d
                start = edges[start]
                d += 1
            return dist

        dist1 = bfs(node1)
        dist2 = bfs(node2)

        minDist, res = n, -1
        for i, (d1, d2) in enumerate(zip(dist1, dist2)):
            max_ = max(d1, d2)
            if max_ < minDist:
                minDist = max_
                res = i
        return res

    def closestMeetingNode2(self, edges: List[int], node1: int, node2: int) -> int:
        """如果[node1, node2]是一群query的话 需要预处理结点是在链上还是环上(求出环分组)，以及在链上的位置和环上的位置"""
        # 两个点在同一颗树里，不经过环 => LCA
        # 两个点经过环 => 预处理深度
        ...
