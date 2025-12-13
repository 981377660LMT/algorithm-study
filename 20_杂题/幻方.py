# 魔方阵 (Magic Square)
# !魔方阵（也称幻方）是一个 $N \times N$ 的方阵，其中填充了从 $1$ 到 $N^2$ 的互不相同的整数，使得：
# 每一行
# 每一列
# 两条对角线
# 上的数字之和都相等。这个相等的和被称为幻和 (Magic Constant)。


def magic_square(n: int) -> list[list[int]]:
    """Siamese Method (暹罗法) 生成奇数阶幻方."""
    assert n & 1, "n must be an odd integer."

    res = [[0] * n for _ in range(n)]
    r, c = 0, (n - 1) // 2
    res[r][c] = 1
    for k in range(2, n * n + 1):
        nr, nc = (r - 1) % n, (c + 1) % n  # 右上
        if res[nr][nc] == 0:
            r, c = nr, nc
        else:
            r = (r + 1) % n  # 正下
        res[r][c] = k
    return res
