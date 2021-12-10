from typing import AsyncContextManager, List
from collections import defaultdict, deque

# 请返回计算机网络变为 空闲 状态的 最早秒数 。
# 编号为 0 的服务器是 主 服务器，其他服务器为 数据 服务器。每个数据服务器都要向主服务器发送信息，并等待回复。
# 如果还没收到任何回复信息，那么该服务器会周期性 重发 信息。数据服务器 i 每 patience[i] 秒都会重发一条信息
# 在 0 秒的开始，所有数据服务器都会发送各自需要处理的信息


# 1.普通bfs求出往返需要的最短路径长度(因为是无权图，所以不用pq)
# 2.根据dist与patience计算出最早空闲时间
INF = 0x7FFFFFFF


class Solution:
    def networkBecomesIdle(self, edges: List[List[int]], patience: List[int]) -> int:
        adjMap = defaultdict(list)
        for u, v in edges:
            adjMap[u].append(v)
            adjMap[v].append(u)

        dist = [INF] * len(adjMap)
        dist[0] = 0
        queue = deque([(0, 0)])
        while queue:
            cost, cur = queue.popleft()
            for next in adjMap[cur]:
                if cost + 2 < dist[next]:
                    dist[next] = cost + 2
                    queue.append((cost + 2, next))

        res = 0
        for dis, pat in zip(dist, patience):
            if pat == 0:
                continue
            div, mod = divmod(dis, pat)
            retry = div - int(mod == 0)
            res = max(res, dis + retry * pat)

        return res + 1


print(Solution().networkBecomesIdle(edges=[[0, 1], [1, 2]], patience=[0, 2, 1]))
print(Solution().networkBecomesIdle(edges=[[0, 1], [0, 2], [1, 2]], patience=[0, 10, 10]))
# 输出：3
# 解释：数据服务器 1 和 2 第 2 秒初收到回复信息。
# 从第 3 秒开始，网络变空闲。
