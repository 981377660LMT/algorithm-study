# 半畳 hanjo
# 铺瓷砖1
# ROW*COL的二维平面上，用2 * 1或者1 * 1的地砖来铺，要求各用A和B块地砖，求出方案数。
# !ROW*COL<=16 2*A+B=ROW*COL

# 轮廓线dp
from functools import lru_cache
from typing import Tuple


def hanjo(ROW: int, COL: int, A: int, B: int) -> int:
    @lru_cache(None)
    def dfs(r: int, c: int, cState: Tuple[int, ...], remain11: int, remain21: int) -> int:
        if r == ROW:
            return 1 if remain11 == remain21 == 0 else 0
        if c == COL:
            return dfs(r + 1, 0, cState, remain11, remain21)

        # 不选择当前位置
        res = dfs(r, c + 1, cState[1:] + (0,), remain11, remain21)

        # 1*1
        if remain11:
            res += dfs(r, c + 1, cState[1:] + (1,), remain11 - 1, remain21)

        if remain21:
            # 规定横竖放的方向是右=>左，下=>上
            # 2*1竖放
            if r and cState[0] == 0:
                res += dfs(r, c + 1, cState[1:] + (1,), remain11, remain21 - 1)
            # 2*1横放
            if c and cState[-1] == 0:
                res += dfs(r, c + 1, cState[1:-1] + (1, 1), remain11, remain21 - 1)

        return res

    if COL > ROW:
        ROW, COL = COL, ROW
    res = dfs(0, 0, tuple([1] * COL), B, A)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    ROW, COL, A, B = map(int, input().split())
    print(hanjo(ROW, COL, A, B))
