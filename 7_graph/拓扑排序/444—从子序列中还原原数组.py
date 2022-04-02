# 从子序列中还原原数组，原数组每个元素唯一
# 拓扑排序
from collections import defaultdict, deque


class Solution:
    def solve(self, sequences):
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

        queue = deque([v for v in vertex if indegree[v] == 0])
        res = []
        while queue:
            cur = queue.popleft()
            res.append(cur)
            for next in adjMap[cur]:
                indegree[next] -= 1
                if indegree[next] == 0:
                    queue.append(next)

        return res


# 重建序列
print(Solution().solve(sequences=[[1, 3], [2, 3], [1, 2]]))
