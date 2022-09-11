"""内向基环树森林中 求到两个给定结点最近的节点"""
# 每个节点 至多 有一条出边(内向基环树)。
# 请你返回一个从 node1 和 node2 都能到达节点的编号，
# 使节点 node1 和节点 node2 到这个节点的距离 较大值最小化。
# 如果有多个答案，请返回 最小 的节点编号。如果答案不存在，返回 -1 。

# !启发：以后还是要用INF来约束 不要在里面写int(1e20) 还是不要急 要认真看题
# !如果是一群query 需要预处理结点是在链上还是环上，以及在链上的位置和环上的位置


from collections import defaultdict, deque
from typing import Hashable, List, TypeVar


class Solution:
    def closestMeetingNode(self, edges: List[int], node1: int, node2: int) -> int:
        """如果[node1, node2]是一群query的话 需要预处理结点是在链上还是环上(求出环分组)，以及在链上的位置和环上的位置"""
        # 两个点在同一颗树里，不经过环 => LCA
        # 两个点经过环 => 预处理深度
        ...

    def closestMeetingNode1(self, edges: List[int], node1: int, node2: int) -> int:
        """两个结点直接bfs找最近结点"""
        INF = int(1e20)
        T = TypeVar("T", bound=Hashable)

        def bfs(adjMap: defaultdict[T, set[T]], start: T) -> defaultdict[T, int]:
            """时间复杂度O(V+E)"""
            dist = defaultdict(lambda: INF, {start: 0})
            queue: deque[tuple[int, T]] = deque([(0, start)])

            while queue:
                _, cur = queue.popleft()
                for next in adjMap[cur]:
                    if dist[next] > dist[cur] + 1:
                        dist[next] = dist[cur] + 1
                        queue.append((dist[next], next))

            return dist

        n = len(edges)
        adjMap = defaultdict(set)
        for u, v in enumerate(edges):
            if v == -1:
                continue
            adjMap[u].add(v)

        dist1, dist2 = bfs(adjMap, node1), bfs(adjMap, node2)
        res, resMax = -1, INF
        for i in range(n):
            max_ = max(dist1[i], dist2[i])
            if max_ < resMax:
                res, resMax = i, max_
        return res


print(Solution().closestMeetingNode1(edges=[3, 0, 5, -1, 3, 4], node1=2, node2=0))
print(Solution().closestMeetingNode1(edges=[5, 4, 5, 4, 3, 6, -1], node1=0, node2=1))
