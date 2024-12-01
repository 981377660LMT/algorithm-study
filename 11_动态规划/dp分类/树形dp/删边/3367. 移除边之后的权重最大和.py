# 3367. 移除边之后的权重最大和
# https://leetcode.cn/problems/maximize-sum-of-weights-after-edge-removals/
# 存在一棵具有 n 个节点的无向树，节点编号为 0 到 n - 1。
# 给你一个长度为 n - 1 的二维整数数组 edges，其中 edges[i] = [ui, vi, wi] 表示在树中节点 ui 和 vi 之间有一条权重为 wi 的边。
#
#
# 你的任务是移除零条或多条边，使得：
# 每个节点与至多 k 个其他节点有边直接相连，其中 k 是给定的输入。
# 剩余边的权重之和 最大化 。
# 返回在进行必要的移除后，剩余边的权重的 最大 可能和。


from typing import List

from selectEdges import selectEdges


class Solution:
    def maximizeSumOfWeights(self, edges: List[List[int]], k: int) -> int:
        n = len(edges) + 1
        limits = [k] * n
        return selectEdges(n, edges, limits)  # type: ignore
