# 拓扑排序性质
# 拓扑排序后点不足n个 => 有环
# 拓扑排序中队列中的点不恒为1 => 多种拓扑序
# 拓扑排序中队列中的点恒为1 => 唯一确定拓扑序
from collections import defaultdict, deque
from typing import List


class Solution:
    def sequenceReconstruction(self, nums: List[int], sequences: List[List[int]]) -> bool:
        adjMap = defaultdict(set)
        indegree = defaultdict(int)
        vertex = set()
        visitedPair = set()
        for seq in sequences:
            vertex |= set(seq)
            for pre, next in zip(seq, seq[1:]):
                if (pre, next) not in visitedPair:
                    visitedPair.add((pre, next))
                    adjMap[pre].add(next)
                    indegree[next] += 1

        if set(nums) != vertex:
            return False

        queue = deque([v for v in vertex if indegree[v] == 0])
        res = []
        count = 0
        while queue:
            if len(queue) > 1:
                return False
            cur = queue.popleft()
            res.append(cur)
            count += 1
            for next in adjMap[cur]:
                indegree[next] -= 1
                if indegree[next] == 0:
                    queue.append(next)

        return count == len(vertex)


print(
    Solution().sequenceReconstruction(
        nums=[4, 1, 5, 2, 6, 3], sequences=[[5, 2, 6, 3], [4, 1, 5, 2]]
    )
)

print(Solution().sequenceReconstruction(nums=[1], sequences=[[1]]))
