from collections import defaultdict, deque
from typing import List

# 有向图，可重复经过点，返回信息从小 A (编号 0 ) 经过 k 轮传递到编号为 n-1 的小伙伴处的方案数；若不能到达，返回 0。


class Solution:
    def numWays(self, n: int, relation: List[List[int]], k: int) -> int:
        adjMap = defaultdict(set)
        for cur, next in relation:
            adjMap[cur].add(next)
        queue = deque([0])

        res = 0
        while queue and k:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjMap[cur]:
                    if next == n - 1 and k == 1:
                        res += 1
                    queue.append(next)
            k -= 1
        return res


print(
    Solution().numWays(n=5, relation=[[0, 2], [2, 1], [3, 4], [2, 3], [1, 4], [2, 0], [0, 4]], k=3)
)
