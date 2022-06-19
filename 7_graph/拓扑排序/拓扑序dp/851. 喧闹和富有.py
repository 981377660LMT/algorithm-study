from typing import List
from collections import defaultdict
from collections import deque

# 1 <= quiet.length = N <= 500
# 如果能够肯定 person x 比 person y 更有钱的话，我们会说 richer[i] = [x, y]

# 返回答案 answer ，其中 answer[x] = y 的前提是，
# 在所有拥有的钱不少于 person x 的人中，
# person y 是最安静的人（也就是安静值 quiet[y] 最小的人）。

# 从上往下，使用父节点更新子节点的更安静的人
# 即对每个下面的结点 求上面的最小值


class Solution:
    # 拓扑排序不易错
    # 从上到下 一层一层地更新传递值(dp)
    def loudAndRich(self, richer: List[List[int]], quiet: List[int]) -> List[int]:
        adjMap = defaultdict(set)
        indeg = [0] * len(quiet)
        for u, v in richer:
            adjMap[u].add(v)
            indeg[v] += 1

        res = [i for i in range(len(quiet))]

        queue = deque([i for i, v in enumerate(indeg) if v == 0])

        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                if quiet[res[cur]] < quiet[res[next]]:
                    res[next] = res[cur]
                indeg[next] -= 1
                if indeg[next] == 0:
                    queue.append(next)

        return res


# print(
#     Solution().loudAndRich(
#         richer=[[1, 0], [2, 1], [3, 1], [3, 7], [4, 3], [5, 3], [6, 3]],
#         quiet=[3, 2, 5, 4, 6, 1, 7, 0],
#     )
# )
# 输出：[5,5,2,5,4,5,6,7]
print(Solution().loudAndRich(richer=[[0, 1], [1, 2]], quiet=[1, 0, 2],))
