from typing import List
from collections import defaultdict, deque

# 请你返回给定图中有效路径里面的 最大颜色值 。如果图中含有环，请返回 -1 。


# 总结：
# 1.拓扑排序判断环
# 2.对每个结点cur的next更新处，都更新next的color的count (每个节点处都有26个颜色的counter)


class Solution:
    def largestPathValue(self, colors: str, edges: List[List[int]]) -> int:
        n = len(colors)
        adjMap, indeg = defaultdict(set), [0] * n
        for u, v in edges:
            indeg[v] += 1
            adjMap[u].add(v)

        colorCounters = [[0] * 26 for _ in range(n)]
        queue = deque([i for i in range(n) if indeg[i] == 0])
        count = 0
        while queue:
            cur = queue.popleft()
            count += 1
            colorCounters[cur][ord(colors[cur]) - ord('a')] += 1

            for next in adjMap[cur]:
                indeg[next] -= 1
                if indeg[next] == 0:
                    queue.append(next)

                # 注意这一步
                for color in range(26):
                    colorCounters[next][color] = max(
                        colorCounters[next][color], colorCounters[cur][color]
                    )

        if count != n:
            return -1

        return max((max(counter) for counter in colorCounters), default=0)


print(Solution().largestPathValue(colors="abaca", edges=[[0, 1], [0, 2], [2, 3], [3, 4]]))
# 输出：3
# 解释：路径 0 -> 2 -> 3 -> 4 含有 3 个颜色为 "a" 的节点（上图中的红色节点）。
