# https://www.acwing.com/problem/content/1171/

# 幼儿园里有 N 个小朋友，老师现在想要给这些小朋友们分配糖果，
# !要求每个小朋友都要分到糖果。(xi>=1)
# 但是小朋友们也有嫉妒心，总是会提出一些要求，比如小明不希望小红分到的糖果比他的多，
# 于是在分配糖果的时候， 老师需要满足小朋友们的 K 个要求。
# 幼儿园的糖果总是有限的，老师想知道他至少需要准备多少个糖果，
# 才能使得每个小朋友都能够分到糖果，并且满足小朋友们所有的要求。


# 接下来 K 行，表示分配糖果时需要满足的关系，每行 3 个数字 X,A,B。

# 如果 X=1．表示第 A 个小朋友分到的糖果必须和第 B 个小朋友分到的糖果一样多。
# 如果 X=2，表示第 A 个小朋友分到的糖果必须少于第 B 个小朋友分到的糖果。
# 如果 X=3，表示第 A 个小朋友分到的糖果必须不少于第 B 个小朋友分到的糖果。
# 如果 X=4，表示第 A 个小朋友分到的糖果必须多于第 B 个小朋友分到的糖果。
# 如果 X=5，表示第 A 个小朋友分到的糖果必须不多于第 B 个小朋友分到的糖果。
# 小朋友编号从 1 到 N。

# !输出一行，表示老师至少需要准备的糖果数，如果不能满足小朋友们的所有要求，就输出 −1。

# n<=1e5

# 最小值=>求最长路
# 无解：存在正环

from collections import defaultdict, deque
from typing import List, Mapping, Tuple

INF = int(1e18)


def spfa(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> Tuple[bool, List[int]]:
    """spfa求虚拟节点为起点的单源最长路 并检测正环"""
    dist = [0] * n
    queue = deque(list(range(n)))
    count = [0] * n
    inQueue = [True] * n

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand > dist[next]:
                dist[next] = cand
                count[next] = count[cur] + 1
                if count[next] >= n:
                    return False, []
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] > dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return True, dist


n, k = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(lambda: -INF))
for _ in range(k):
    kind, u, v = map(int, input().split())
    u, v = u - 1, v - 1
    if kind == 1:  # u == v
        adjMap[u][v] = max(adjMap[u][v], 0)
        adjMap[v][u] = max(adjMap[v][u], 0)
    elif kind == 2:  # v - u >= 1
        adjMap[u][v] = max(adjMap[u][v], 1)
    elif kind == 3:  # u - v >= 0
        adjMap[v][u] = max(adjMap[v][u], 0)
    elif kind == 4:  # u - v >= 1
        adjMap[v][u] = max(adjMap[v][u], 1)
    elif kind == 5:  # v - u >= 0
        adjMap[u][v] = max(adjMap[u][v], 0)

# !所有xi>=1 引入虚拟源点n 即xi -xn >= 1
SUPER_NODE = n
for i in range(n):
    adjMap[SUPER_NODE][i] = max(adjMap[SUPER_NODE][i], 1)

ok, dist = spfa(n + 1, adjMap)
if not ok:
    print(-1)
else:
    print(sum(dist))
