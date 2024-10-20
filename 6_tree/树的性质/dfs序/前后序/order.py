#       0
#     / | \
#   1   2  3
#  / \     |
# 4   5    6
#
# !preOrder:
#  0 => [0, 7)
#  1 => [1, 4)
#  4 => [2, 3)
#  5 => [3, 4)
#  2 => [4, 5)
#  3 => [5, 7)
#  6 => [6, 7)
#
# !postOrder:
#  4 => [0, 1)
#  5 => [1, 2)
#  1 => [0, 3)
#  2 => [3, 4)
#  6 => [4, 5)
#  3 => [4, 6)
#  0 => [0, 7)

from typing import List, Tuple


def dfsPreOrder(tree: List[List[int]], root=0) -> Tuple[List[int], List[int]]:
    """前序遍历dfs序.

    # !data[lid[i]] = values[i]
    """
    n = len(tree)
    lid, rid = [0] * n, [0] * n
    dfn = 0

    def dfs(cur: int, pre: int) -> None:
        nonlocal dfn
        lid[cur] = dfn
        dfn += 1
        for next_ in tree[cur]:
            if next_ != pre:
                dfs(next_, cur)
        rid[cur] = dfn

    dfs(root, -1)
    return lid, rid


def dfsPostOrder(tree: List[List[int]], root=0) -> Tuple[List[int], List[int]]:
    """后序遍历dfs序.

    # !data[rid[i]-1] = values[i]
    """
    n = len(tree)
    lid, rid = [0] * n, [0] * n
    dfn = 0

    def dfs(cur: int, pre: int) -> None:
        nonlocal dfn
        lid[cur] = dfn
        for next_ in tree[cur]:
            if next_ != pre:
                dfs(next_, cur)
        dfn += 1
        rid[cur] = dfn

    dfs(root, -1)
    return lid, rid


if __name__ == "__main__":
    n = 7
    edges = [(0, 1), (0, 2), (0, 3), (1, 4), (1, 5), (3, 6)]
    tree = [[] for _ in range(n)]
    for u, v in edges:
        tree[u].append(v)
        tree[v].append(u)

    assert dfsPreOrder(tree) == ([0, 1, 4, 5, 2, 3, 6], [7, 4, 5, 7, 3, 4, 7])
    assert dfsPostOrder(tree) == ([0, 0, 3, 4, 0, 1, 4], [7, 3, 4, 6, 1, 2, 5])
