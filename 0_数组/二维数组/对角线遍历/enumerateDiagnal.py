# 对角线遍历/遍历对角线


from typing import Any, Generator, List, Sequence, Tuple


def enumerateDiagnal(
    grid: Sequence[Sequence[Any]], direction=0, upToDown=True
) -> Generator[List[Tuple[int, int]], None, None]:
    """
    对角线遍历二维矩阵.左上角为(0,0).

    Args:
        - grid (Sequence[Sequence[Any]]): 二维数组
        - direction (int): 遍历方向.
          - `0: ↘`, 从左上到右下. 同一条对角线上 `r-c` 为定值.
          - `1: ↖`, 从右下到左上. 同一条对角线上 `r-c` 为定值.
          - `2: ↙`, 从右上到左下. 同一条对角线上 `r+c` 为定值.
          - `3: ↗`, 从左下到右上. 同一条对角线上 `r+c` 为定值.
        - upToDown (bool): 是否从上到下遍历. 默认为True.

    Returns:
        Generator[List[Tuple[int, int]], None, None]: 每个对角线上的元素坐标(r,c)的列表.
    """
    if direction not in (0, 1, 2, 3):
        raise ValueError("direction must be in (0, 1, 2, 3)")

    ROW, COL = len(grid), len(grid[0])

    if direction == 0:
        for key in range(-COL + 1, ROW) if upToDown else range(ROW - 1, -COL, -1):
            r = 0 if key < 0 else key
            c = r - key
            cur = []
            while r < ROW and c < COL:
                cur.append((r, c))
                r += 1
                c += 1
            if cur:
                yield cur

    elif direction == 1:
        for key in range(-COL + 1, ROW) if upToDown else range(ROW - 1, -COL, -1):
            r = ROW - 1 if key > ROW - COL else key + COL - 1
            c = r - key
            cur = []
            while r >= 0 and c >= 0:
                cur.append((r, c))
                r -= 1
                c -= 1
            if cur:
                yield cur

    elif direction == 2:
        for key in range(ROW + COL - 1) if upToDown else range(ROW + COL - 2, -1, -1):
            r = 0 if key < COL else key - COL + 1
            c = key - r
            cur = []
            while r < ROW and c >= 0:
                cur.append((r, c))
                r += 1
                c -= 1
            if cur:
                yield cur

    elif direction == 3:
        for key in range(ROW + COL - 1) if upToDown else range(ROW + COL - 2, -1, -1):
            r = ROW - 1 if key >= ROW else key
            c = key - r
            cur = []
            while r >= 0 and c < COL:
                cur.append((r, c))
                r -= 1
                c += 1
            if cur:
                yield cur


if __name__ == "__main__":
    mat = [[1, 2, 3], [4, 5, 6]]
    print(list(enumerateDiagnal(mat, 0)))
    assert list(enumerateDiagnal(mat, 0)) == [
        [(0, 2)],
        [(0, 1), (1, 2)],
        [(0, 0), (1, 1)],
        [(1, 0)],
    ]
    assert list(enumerateDiagnal(mat, 1)) == [
        [(0, 2)],
        [(1, 2), (0, 1)],
        [(1, 1), (0, 0)],
        [(1, 0)],
    ]
    assert list(enumerateDiagnal(mat, 2)) == [
        [(0, 0)],
        [(0, 1), (1, 0)],
        [(0, 2), (1, 1)],
        [(1, 2)],
    ]
    assert list(enumerateDiagnal(mat, 3)) == [
        [(0, 0)],
        [(1, 0), (0, 1)],
        [(1, 1), (0, 2)],
        [(1, 2)],
    ]

    mat = [list(col[::-1]) for col in zip(*mat)]

    assert list(enumerateDiagnal(mat, 0)) == [
        [(0, 1)],
        [(0, 0), (1, 1)],
        [(1, 0), (2, 1)],
        [(2, 0)],
    ]
    assert list(enumerateDiagnal(mat, 1)) == [
        [(0, 1)],
        [(1, 1), (0, 0)],
        [(2, 1), (1, 0)],
        [(2, 0)],
    ]
    assert list(enumerateDiagnal(mat, 2)) == [
        [(0, 0)],
        [(0, 1), (1, 0)],
        [(1, 1), (2, 0)],
        [(2, 1)],
    ]
    assert list(enumerateDiagnal(mat, 3)) == [
        [(0, 0)],
        [(1, 0), (0, 1)],
        [(2, 0), (1, 1)],
        [(2, 1)],
    ]

    print(*enumerateDiagnal(mat, 0, upToDown=False), sep="\n")
    print(*enumerateDiagnal(mat, 1, upToDown=False), sep="\n")
