# 对角线遍历/遍历对角线


from typing import Generator, List, Tuple


def enumerateDiagnal(
    row: int, col: int, direction=0, upToDown=True
) -> Generator[List[Tuple[int, int]], None, None]:
    """
    对角线遍历二维矩阵.左上角为(0,0).

    Args:
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

    if direction == 0:
        for key in range(-col + 1, row) if upToDown else range(row - 1, -col, -1):
            r = 0 if key < 0 else key
            c = r - key
            cur = []
            while r < row and c < col:
                cur.append((r, c))
                r += 1
                c += 1
            if cur:
                yield cur

    elif direction == 1:
        for key in range(-col + 1, row) if upToDown else range(row - 1, -col, -1):
            r = row - 1 if key > row - col else key + col - 1
            c = r - key
            cur = []
            while r >= 0 and c >= 0:
                cur.append((r, c))
                r -= 1
                c -= 1
            if cur:
                yield cur

    elif direction == 2:
        for key in range(row + col - 1) if upToDown else range(row + col - 2, -1, -1):
            r = 0 if key < col else key - col + 1
            c = key - r
            cur = []
            while r < row and c >= 0:
                cur.append((r, c))
                r += 1
                c -= 1
            if cur:
                yield cur

    elif direction == 3:
        for key in range(row + col - 1) if upToDown else range(row + col - 2, -1, -1):
            r = row - 1 if key >= row else key
            c = key - r
            cur = []
            while r >= 0 and c < col:
                cur.append((r, c))
                r -= 1
                c += 1
            if cur:
                yield cur


if __name__ == "__main__":
    mat = [[1, 2, 3], [4, 5, 6]]
    row, col = len(mat), len(mat[0])
    assert list(enumerateDiagnal(row, col)) == [
        [(0, 2)],
        [(0, 1), (1, 2)],
        [(0, 0), (1, 1)],
        [(1, 0)],
    ]
    assert list(enumerateDiagnal(row, col, 1)) == [
        [(0, 2)],
        [(1, 2), (0, 1)],
        [(1, 1), (0, 0)],
        [(1, 0)],
    ]
    assert list(enumerateDiagnal(row, col, 2)) == [
        [(0, 0)],
        [(0, 1), (1, 0)],
        [(0, 2), (1, 1)],
        [(1, 2)],
    ]
    assert list(enumerateDiagnal(row, col, 3)) == [
        [(0, 0)],
        [(1, 0), (0, 1)],
        [(1, 1), (0, 2)],
        [(1, 2)],
    ]

    mat = [list(col[::-1]) for col in zip(*mat)]
    row, col = len(mat), len(mat[0])
    assert list(enumerateDiagnal(row, col)) == [
        [(0, 1)],
        [(0, 0), (1, 1)],
        [(1, 0), (2, 1)],
        [(2, 0)],
    ]
    assert list(enumerateDiagnal(row, col, 1)) == [
        [(0, 1)],
        [(1, 1), (0, 0)],
        [(2, 1), (1, 0)],
        [(2, 0)],
    ]
    assert list(enumerateDiagnal(row, col, 2)) == [
        [(0, 0)],
        [(0, 1), (1, 0)],
        [(1, 1), (2, 0)],
        [(2, 1)],
    ]
    assert list(enumerateDiagnal(row, col, 3)) == [
        [(0, 0)],
        [(1, 0), (0, 1)],
        [(2, 0), (1, 1)],
        [(2, 1)],
    ]

    print(*enumerateDiagnal(row, col, upToDown=False), sep="\n")
    print(*enumerateDiagnal(row, col, upToDown=False), sep="\n")
