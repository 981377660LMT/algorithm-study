# それぞれのテレポーターは 2 つの町を双方向に結んでおり、テレポーターを使用する事によってその 2 つの町の間を 1 分で移動することができます
# !假设所有缺少的边都与i相连 对这n个不同的i值 求 0 到 n-1 的最短路 不存在输出-1

# n,m<=3e5


# !从起点终点出发 枚举中间的传送点 求最短路
# 传送们虚拟节点统一作为-1 每次求解相当于-1与i连了权值为0的边

# 三种情况:
# 1.不走传送门
# 2.起点走传送门出口-1,终点走传送门入口i
# 3.终点走传送门-1,,起点走传送门入口i

# !虚拟点的情况需要枚举使用虚拟点的情况

from collections import defaultdict, deque
import sys
from typing import DefaultDict

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    def bfs(start: int) -> DefaultDict[int, int]:
        dist = defaultdict(lambda: int(1e18), {start: 0})
        queue = deque([start])
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                if dist[next] > dist[cur] + 1:
                    dist[next] = dist[cur] + 1
                    queue.append(next)
        return dist

    n, m = map(int, input().split())
    adjMap = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1  # !传送点出口统一作为-1
        adjMap[u].add(v)
        adjMap[v].add(u)

    dist1, dist2 = bfs(0), bfs(n - 1)
    res = []
    # !把i号点作为传送门入口
    for i in range(n):
        cand1 = dist1[n - 1]
        cand2 = dist1[i] + dist2[-1]
        cand3 = dist2[i] + dist1[-1]
        cur = min(cand1, cand2, cand3)
        res.append(cur if cur < int(1e18) else -1)

    print(*res)


while True:
    try:
        main()
    except (EOFError, ValueError):
        break
