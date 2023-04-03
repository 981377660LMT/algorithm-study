# 在线bfs.
#   不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
#   setUsed(u)：将 u 标记为已访问。
#   findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `None`。
# https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/


from collections import deque
from typing import Callable, List, Optional

INF = int(1e18)


def onlineBfs(
    n: int, start: int, setUsed: Callable[[int], None], findUnused: Callable[[int], Optional[int]]
) -> List[int]:
    """在线bfs。不预先给出图, 而是通过两个函数 setUsed 和 findUnused 来在线寻找边。

    Args:
        n (int): 顶点数。
        start (int): 起点。
        setUsed (Callable[[int], None]): 将 u 标记为已访问。
        findUnused (Callable[[int], int]): 找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `None`。

    Returns:
        List[int]: 从起点到各个点的距离。
    """
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    setUsed(start)
    while queue:
        cur = queue.popleft()
        while True:
            next_ = findUnused(cur)
            if next_ is None:
                break
            dist[next_] = dist[cur] + 1
            queue.append(next_)
            setUsed(next_)
    return dist
