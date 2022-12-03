"""
PAWN是国际象棋的一个棋子,每方各有八个,每步只可前行或者斜前向吃子,
除首次移动可选择1或2格外,每步只可向前移动一格；
在抵达对方底线后可选择升格为我方除王以外任何棋子。
"""

# 给出(2n+1)×(2n +1)的棋盘上m枚黑子的坐标
# 开始时一枚白子在(0, n)处,白子的移动策略如下:
# ·若(i+1,j)没有黑子,则可以移动到(i+1,j)
# ·若(i+1,j-1)有黑子,则可以移动到(i+1,j-1)
# ·若(i+1,j+1)有黑子,则可以移动到(i+1,j+1)
# 问白子可能到达第2n行的哪些坐标。

# !n<=1e9 m<=2e5


from collections import defaultdict
from typing import List, Tuple


def whitePawn(n: int, blackChess: List[Tuple[int, int]]) -> int:
    rowBlack = defaultdict(list)
    for r, c in blackChess:
        rowBlack[r].append(c)

    res = set([n])  # 当前可以到达的列
    for row in sorted(rowBlack):
        good, bad = set(), set()
        for col in rowBlack[row]:
            if (col - 1 >= 0 and col - 1 in res) or (col + 1 <= 2 * n and col + 1 in res):
                good.add(col)
            else:
                bad.add(col)
        res -= bad
        res |= good
    return len(res)


if __name__ == "__main__":
    n, m = map(int, input().split())
    blackChess = [tuple(map(int, input().split())) for _ in range(m)]
    print(whitePawn(n, blackChess))
