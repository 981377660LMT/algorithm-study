from collections import defaultdict
from typing import Iterable, List, Mapping, Sequence, Tuple, Union

AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Tree = Union[AdjList, AdjMap]


def getCenter(n: int, tree: "Tree") -> Tuple[int, List[int]]:
    """求重心"""

    def dfs(cur: int, pre: int) -> None:
        subsize[cur] = 1
        for next in tree[cur]:
            if next == pre:
                continue
            dfs(next, cur)
            subsize[cur] += subsize[next]
            maxsize[cur] = max(maxsize[cur], subsize[next])
        maxsize[cur] = max(maxsize[cur], n - subsize[cur])
        if maxsize[cur] <= n / 2:
            res.append(cur)

    res = []
    maxsize = [0] * n  # 最大连通块大小
    subsize = [0] * n  # 子树大小
    dfs(0, -1)
    return maxsize[res[0]], res


if __name__ == "__main__":
    adjMap = defaultdict(set)
    edges = [[1, 0], [1, 2], [1, 3]]
    for u, v in edges:
        adjMap[u].add(v)
        adjMap[v].add(u)
    print(getCenter(4, adjMap))
