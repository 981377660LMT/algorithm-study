# 矩阵转置
# O(n+m)log(n+m)
# https://maspypy.github.io/library/linalg/xor/transpose.hpp


from typing import List


def transpose(row: int, col: int, matrix1D: List[int], inplace=False) -> List[int]:
    """矩阵转置

    matrix1D:每个元素是状压的数字
    inplace:是否修改原矩阵
    """
    assert row == len(matrix1D)
    m = matrix1D[:] if not inplace else matrix1D
    log = 0
    max_ = max(row, col)
    while (1 << log) < max_:
        log += 1
    if len(m) < (1 << log):
        m += [0] * ((1 << log) - len(m))
    w = 1 << log
    mask = 1
    for i in range(log):
        mask = mask | (mask << (1 << i))
    for t in range(log):
        w >>= 1
        mask = mask ^ (mask >> w)
        for i in range(1 << t):
            for j in range(w):
                m[w * 2 * i + j] = ((m[w * (2 * i + 1) + j] << w) & mask) ^ m[w * 2 * i + j]
                m[w * (2 * i + 1) + j] = ((m[w * 2 * i + j] & mask) >> w) ^ m[w * (2 * i + 1) + j]
                m[w * 2 * i + j] = ((m[w * (2 * i + 1) + j] << w) & mask) ^ m[w * 2 * i + j]
    return m[:col]


if __name__ == "__main__":
    A = [1, 2, 3]  # [[1,0,0],[0,1,0],[1,1,0]]
    res = transpose(3, 3, A)  # [[1,0,1],[0,1,1],[0,0,0]]
    print(res)  # [5, 6, 0]
