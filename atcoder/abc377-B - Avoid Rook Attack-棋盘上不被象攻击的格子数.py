# abc377-B - Avoid Rook Attack-棋盘上不被象攻击的格子数


from typing import List, Tuple


def avoidRookAttack(n: int, m: int, rooks: List[Tuple[int, int]]) -> int:
    rowStates, colStates = [False] * n, [False] * m
    for r, c in rooks:
        rowStates[r] = True
        colStates[c] = True
    c1, c2 = rowStates.count(False), colStates.count(False)
    return c1 * c2


if __name__ == "__main__":
    grid = [list(input()) for _ in range(8)]
    rooks = []
    for i in range(8):
        for j in range(8):
            if grid[i][j] == "#":
                rooks.append((i, j))
    res = avoidRookAttack(8, 8, rooks)  # type: ignore
    print(res)
