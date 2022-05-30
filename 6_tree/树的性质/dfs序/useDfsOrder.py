from collections import defaultdict, namedtuple
from typing import DefaultDict, Set, Tuple


def useDfsOrder(n: int, tree: DefaultDict[int, Set[int]]):
    """dfs序

    Args:
        n (int): 树节点从0开始，根节点为0
        tree (DefaultDict[int, Set[int]]): 无向图邻接表
    """

    def dfs(cur: int, pre: int) -> None:
        nonlocal dfsId
        starts[cur] = dfsId
        for child in tree[cur]:
            if child == pre:
                continue
            dfs(child, cur)
        ends[cur] = dfsId
        dfsId += 1

    starts = [-1] * (n + 1)
    ends = [-1] * (n + 1)
    dfsId = 1
    dfs(0, -1)

    def queryRange(root: int) -> Tuple[int, int]:
        """求子树映射到的区间

        Args:
            root (int): 根节点

        Returns:
            Tuple[int, int]: [start, end] 1 <= start <= end <= n
        """
        assert 0 <= root < n
        return starts[root], ends[root]

    def queryId(root: int) -> int:
        """求root自身的dfsId

        Args:
            root (int): 根节点

        Returns:
            int: id  1 <= id <= n
        """
        assert 0 <= root < n
        return ends[root]

    return namedtuple('useDfsOrder', ['queryRange', 'queryId'])(queryRange, queryId)


if __name__ == '__main__':
    n = 4
    edges = [[0, 1], [1, 2], [2, 3]]
    tree = defaultdict(set)
    for u, v in edges:
        tree[u].add(v)
        tree[v].add(u)
    queryRange, queryId = useDfsOrder(n, tree)
    print(queryRange(1))
    print(queryRange(2))
    print(queryRange(3))
    print(queryId(1))
    print(queryId(2))
    print(queryId(3))
