# 拓扑排序性质
# 拓扑排序后点不足n个 => 有环
# 拓扑排序中队列中的点不恒为1 => 多种拓扑序
# 拓扑排序中队列中的点恒为1 => 唯一确定拓扑序
from collections import deque
from typing import List


class Solution:
    def sequenceReconstruction(self, nums: List[int], sequences: List[List[int]]) -> bool:
        """
        注意拓扑排序可以有重边
        nums 是范围为 [1,n] 的整数的排列
        """
        n = len(nums)
        adjList = [[] for _ in range(n)]
        indeg = [0] * n
        for seq in sequences:
            for u, v in zip(seq, seq[1:]):
                u, v = u - 1, v - 1
                adjList[u].append(v)
                indeg[v] += 1

        queue = deque([i for i in range(n) if indeg[i] == 0])
        while queue:
            if len(queue) > 1:
                return False
            cur = queue.popleft()
            for next in adjList[cur]:
                indeg[next] -= 1
                if indeg[next] == 0:
                    queue.append(next)

        return True


print(
    Solution().sequenceReconstruction(
        nums=[4, 1, 5, 2, 6, 3], sequences=[[5, 2, 6, 3], [4, 1, 5, 2]]
    )
)

print(Solution().sequenceReconstruction(nums=[1], sequences=[[1]]))
