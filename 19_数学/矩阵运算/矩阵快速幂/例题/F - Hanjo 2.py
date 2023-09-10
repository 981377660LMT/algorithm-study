# 铺瓷砖2
# ROW*COL的二维平面上，用1 * 1或者2 * 1的地砖来铺，要求铺满，求出方案数。
# 数据范围ROW <= 6, COL <= 1e12

# !dfs预处理状态转移+矩阵快速幂'
# dp[col][state] 表示当前在第col列，状态为state的方案数
import sys
from matqpow import matqpow1, matmul

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def hanjo2(ROW: int, COL: int) -> int:
    def dfs(col1: int, col2: int, start: int) -> None:
        if col1 == (1 << ROW) - 1:
            trans[start][col2] += 1
            return
        for i in range(ROW):
            if col1 & (1 << i):
                continue
            # 1*1瓷砖
            dfs(col1 | (1 << i), col2, start)
            # 2*1瓷砖竖铺
            if i + 1 < ROW and not col1 & (1 << (i + 1)):
                dfs(col1 | (1 << i) | (1 << (i + 1)), col2, start)
            # 2*1瓷砖横铺
            dfs(col1 | (1 << i), col2 | (1 << i), start)
            break  # !注意这里 每次铺从最上面开始铺

    res = [[0] * (1 << ROW) for _ in range(1 << ROW)]
    trans = [[0] * (1 << ROW) for _ in range(1 << ROW)]
    for i in range(1 << ROW):
        dfs(i, 0, i)
        res[i][i] = 1

    trans = matqpow1(trans, COL, MOD)
    res = matmul(trans, res, MOD)
    return res[0][0]


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    print(hanjo2(ROW, COL))
