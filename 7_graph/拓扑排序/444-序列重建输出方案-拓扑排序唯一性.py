# 给定0-n-1这n个下标处元素的大小关系,如果能唯一确定顺序,则输出顺序,否则输出空数组

from collections import deque
from typing import List, Optional, Tuple


def getUniqueOrder(n: int, deps: List[Tuple[int, int]]) -> Optional[List[int]]:
    adjList = [[] for _ in range(n)]
    indeg = [0] * n
    for x, y in deps:
        adjList[x].append(y)
        indeg[y] += 1
    queue = deque([i for i in range(n) if indeg[i] == 0])
    order = []
    while queue:
        if len(queue) > 1:
            return
        cur = queue.popleft()
        order.append(cur)
        for next in adjList[cur]:
            indeg[next] -= 1
            if indeg[next] == 0:
                queue.append(next)
    if len(order) < n:
        return
    return order  # 拓扑序


if __name__ == "__main__":
    n, m = map(int, input().split())
    deps = []
    for _ in range(m):
        x, y = map(int, input().split())  # x<y
        x, y = x - 1, y - 1
        deps.append((x, y))
    order = getUniqueOrder(n, deps)
    if order is None:
        print("No")
    else:
        print("Yes")
        res = sorted(list(range(1, n + 1)), key=lambda i: order[i - 1])
        print(*res)
