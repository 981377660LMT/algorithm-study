# (())() → 得到(n//2)+1 个结点,其中0号结点是虚拟根结点
# graph: [[1 3] [2] [] []] (有向邻接表)
# leftRight: [[0 6) [0 4) [1 3) [4 6)] (每个顶点的括号范围,左闭右开)
#
#           0 [0,6)
#          / \
#   [0,4) 1   3 [4,6)
#        /
# [1,3) 2

# 括号树

from typing import List, Tuple


def buildBracketTree(s: str) -> Tuple[List[List[int]], List[Tuple[int, int]]]:
    """从有效的括号序列构建括号树.

    Args:
        s (str): 有效的括号序列,由'('和')'组成.

    Returns:
        Tuple[List[List[int]], List[Tuple[int, int]]]:
        `树的有向邻接表`和`每个结点对应的括号范围`.
        长为2*n的括号序列对应0~n这n+1个结点.0号结点是虚拟根结点.
    """
    n = len(s) // 2
    tree = [[] for _ in range(n + 1)]
    leftRight: List[List[int]] = [[-1, -1] for _ in range(n + 1)]
    leftRight[0] = [0, len(s)]
    parent = [-1] * (n + 1)
    cur, next_ = 0, 1

    for i, c in enumerate(s):
        if c == "(":
            tree[cur].append(next_)
            parent[next_] = cur
            leftRight[next_][0] = i
            cur, next_ = next_, next_ + 1
        else:
            leftRight[cur][1] = i + 1
            cur = parent[cur]

    return tree, [tuple(x) for x in leftRight]


if __name__ == "__main__":
    s = "(())()"
    tree, leftRight = buildBracketTree(s)
    assert tree == [[1, 3], [2], [], []]
    assert leftRight == [(0, 6), (0, 4), (1, 3), (4, 6)]
