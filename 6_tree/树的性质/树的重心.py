from collections import defaultdict
from typing import Iterable, List, Mapping, Sequence, Union

AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Tree = Union[AdjList, AdjMap]


def getCenter(n: int, tree: "Tree", root=0) -> List[int]:
    """求重心"""

    def dfs(cur: int, pre: int) -> None:
        subsize[cur] = 1
        for next in tree[cur]:
            if next == pre:
                continue
            dfs(next, cur)
            subsize[cur] += subsize[next]
            weight[cur] = max(weight[cur], subsize[next])
        weight[cur] = max(weight[cur], n - subsize[cur])
        if weight[cur] <= n / 2:
            res.append(cur)

    res = []
    weight = [0] * n  # 节点的`重量`，即以该节点为根的子树的最大节点数
    subsize = [0] * n  # 子树大小
    dfs(root, -1)
    return res


if __name__ == "__main__":
    adjMap = defaultdict(set)
    edges = [[1, 0], [1, 2], [1, 3]]
    for u, v in edges:
        adjMap[u].add(v)
        adjMap[v].add(u)
    print(getCenter(4, adjMap))
