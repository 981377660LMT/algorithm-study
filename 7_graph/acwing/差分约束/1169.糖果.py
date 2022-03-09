# 幼儿园里有 N 个小朋友，老师现在想要给这些小朋友们分配糖果，要求每个小朋友都要分到糖果。
# 但是小朋友们也有嫉妒心，总是会提出一些要求，比如小明不希望小红分到的糖果比他的多，于是在分配糖果的时候， 老师需要满足小朋友们的 K 个要求。
# 幼儿园的糖果总是有限的，老师想知道他至少需要准备多少个糖果，才能使得每个小朋友都能够分到糖果，并且满足小朋友们所有的要求。


# 接下来 K 行，表示分配糖果时需要满足的关系，每行 3 个数字 X,A,B。

# 如果 X=1．表示第 A 个小朋友分到的糖果必须和第 B 个小朋友分到的糖果一样多。
# 如果 X=2，表示第 A 个小朋友分到的糖果必须少于第 B 个小朋友分到的糖果。
# 如果 X=3，表示第 A 个小朋友分到的糖果必须不少于第 B 个小朋友分到的糖果。
# 如果 X=4，表示第 A 个小朋友分到的糖果必须多于第 B 个小朋友分到的糖果。
# 如果 X=5，表示第 A 个小朋友分到的糖果必须不多于第 B 个小朋友分到的糖果。
# 小朋友编号从 1 到 N。

# 输出一行，表示老师至少需要准备的糖果数，如果不能满足小朋友们的所有要求，就输出 −1。

# n<=1e5

# 最小值=>求最长路
# 无解：存在正环

from collections import defaultdict, deque
from typing import List, Tuple


n, k = map(int, input().split())
adjMap = defaultdict(set)
for _ in range(k):

    kind, u, v = map(int, input().split())
    u, v = u - 1, v - 1
    # 所有的相对条件
    if kind == 1:
        adjMap[u].add((v, 0))
        adjMap[v].add((u, 0))
    elif kind == 2:
        # u+1<=v
        adjMap[u].add((v, 1))
    elif kind == 3:
        # u>=v
        adjMap[v].add((u, 0))
    elif kind == 4:
        # u>=v+1
        adjMap[v].add((u, 1))
    elif kind == 5:
        # u<=v
        adjMap[u].add((v, 0))

# 绝对条件 虚拟源点n 所有xi>=1 即xi>=xn+1 xn=0  n能到任意点就就能到任意边
SUPER_NODE = n
for i in range(n):
    adjMap[SUPER_NODE].add((i, 1))


def spfa(n: int, adjMap: defaultdict) -> Tuple[bool, List[int]]:
    """spfa求单源最长路，顺便判断正环"""
    # dist = [-int(1e20)] * (n + 1)
    dist = [0] * (n + 1)
    dist[SUPER_NODE] = 0

    queue = deque([SUPER_NODE])
    count = [0] * (n + 1)  # 边数
    isInqueue = [False] * (n + 1)  # 在队列里的点
    isInqueue[SUPER_NODE] = True

    while queue:
        cur = queue.popleft()
        isInqueue[cur] = False  # 点从队列出来了

        # 更新过谁，就拿谁去更新别人
        for next, weight in adjMap[cur]:
            if dist[cur] + weight > dist[next]:
                dist[next] = dist[cur] + weight
                count[next] = count[cur] + 1
                if count[next] >= n + 1:  # n+1 个点
                    return False, []
                if not isInqueue[next]:
                    isInqueue[next] = True
                    # 酸辣粉优化：队列不为空，且当前元素距离小于队头，则加入队头，否则加入队尾
                    if queue and dist[next] > dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return True, dist


ok, dist = spfa(n, adjMap)
if not ok:
    print(-1)
else:
    print(sum(dist))

