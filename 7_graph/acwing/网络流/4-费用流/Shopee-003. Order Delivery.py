# https://leetcode.cn/problems/0pXm7y/solution/by-981377660lmt-3uiw/

# 将商品X从各个`仓库`运送到各个`城市`，求出最小的总费用，题目保证存在运输方案
# 题目保证存在运输方案(即存在最小费用最大流)

# 1. 无权图最短路求出各个仓库到各个城市的最短距离，进而处理出运输费用 dist[i][j]
# 2. 建立虚拟源点和汇点，求解`仓库`=>`城市`网络的最小费用最大流

from MinCostMaxFlow import MinCostMaxFlow

from collections import defaultdict, deque

# 1. 输入
n, warehouse, m = map(int, input().split())  # 城市 仓库 道路 数量 点数<=20 边数<=200
adjMap = defaultdict(set)  # 城市道路的无向无权图(边权全部为1)
for _ in range(m):
    u, v = map(int, input().split())
    adjMap[u].add(v)
    adjMap[v].add(u)

stores = []  # (仓库的id 仓库存货数 仓库每公里的送货费)
for _ in range(warehouse):
    remain, fee, pos = map(int, input().split())
    stores.append((pos, remain, fee))

orderCount = int(input())  # 订单数量<=1e5
needs = defaultdict(int)  # (城市的id 物品在城市的需求量)
for _ in range(orderCount):
    need, pos = map(int, input().split())
    needs[pos] += need


# 2. bfs处理各个仓库到城市的最短距离
def bfs(start: int) -> None:
    queue = deque([(start, 0)])
    visited = set([start])
    while queue:
        cur, step = queue.popleft()
        dist[start][cur] = min(dist[start][cur], step)
        for next in adjMap[cur]:
            if next not in visited:
                visited.add(next)
                queue.append((next, step + 1))


dist = defaultdict(lambda: defaultdict(lambda: int(1e20)))
for start, *_ in stores:
    bfs(start)


# 3. 建图求最小费用最大流
START, END, OFFSET = 200, 201, 100
mcmf = MinCostMaxFlow(202, START, END)
for i, remain, _ in stores:
    mcmf.addEdge(START, i, remain, 0)  # 虚拟源点提货物
for j, need in needs.items():
    mcmf.addEdge(j + OFFSET, END, need, 0)  # 虚拟汇点接受货物
for i, remain, fee in stores:
    for j, need in needs.items():
        mcmf.addEdge(i, j + OFFSET, remain, fee * dist[i][j])  # 仓库转移虚拟源点的货物
print(mcmf.work()[1])
