"""
从叶子结点往上dp
后序dfs和从叶子往上拓扑是一样的

# !叶子结点:度数为0(只有一个结点时)或者1的结点
"""
from collections import deque
from typing import Iterable, List, Mapping, Sequence, Union

AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Tree = Union[AdjList, AdjMap]


def treeDp(n: int, tree: "Tree", deg: List[int]) -> List[int]:
    """从叶子结点往上dp

    Args:
        n (int): 结点数
        tree (Tree): 树(无向图)
        deg (List[int]): 每个结点的度数
    Returns:
        List[int]: 每个结点的dp值
    """
    queue, visited = deque(), [False] * n
    for i in range(n):
        if deg[i] <= 1:  # !叶子结点(包括孤立的结点)
            queue.append(i)
            visited[i] = True

    dp = [1] * n  # 示例:求每个子树的结点个数
    while queue:
        cur = queue.popleft()
        for next in tree[cur]:
            if not visited[next]:
                dp[next] += dp[cur]  # !dp转移逻辑
                deg[next] -= 1
                if deg[next] == 1:
                    visited[next] = True
                    queue.append(next)
    return dp


if __name__ == "__main__":
    n, edges = 3, [[0, 1], [1, 2]]
    adjList = [[] for _ in range(n)]
    deg = [0] * n
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
        deg[u] += 1
        deg[v] += 1

    print(treeDp(n, adjList, deg))
