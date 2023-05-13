# https://leetcode.cn/problems/special-binary-string/
# 特殊的二进制序列是具有以下两个性质的二进制序列：


# 0 的数量与 1 的数量相等。
# 二进制序列的每一个前缀码中 1 的数量要大于等于 0 的数量。
# (括号树的最小表示法)

from typing import List, Tuple
from minLexEulerTour import minLexEulerTour01


class Solution:
    def makeLargestSpecial(self, s: str) -> str:
        bracket = ["(" if c == "1" else ")" for c in s]
        tree, _ = buildBracketTree("".join(bracket))
        res = minLexEulerTour01(tree)
        res = [v ^ 1 for v in res]
        return "".join(map(str, res[1:-1]))


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
