from typing import List, TypeVar


T = TypeVar("T")


def rotateRight(matrix: List[List[T]]) -> List[List[T]]:
    """顺时针旋转矩阵90度."""
    return [list(col[::-1]) for col in zip(*matrix)]


def rotateLeft(matrix: List[List[T]]) -> List[List[T]]:
    """逆时针旋转矩阵90度."""
    return [list(col) for col in zip(*matrix)][::-1]


def transpose(matrix: List[List[T]]) -> List[List[T]]:
    """矩阵转置."""
    return [list(row) for row in zip(*matrix)]


if __name__ == "__main__":
    # 原矩阵:     右旋转:     左旋转:
    # 1 2 3      7 4 1      3 6 9
    # 4 5 6  ->  8 5 2  or  2 5 8
    # 7 8 9      9 6 3      1 4 7
    grid = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
    assert rotateRight(grid) == [[7, 4, 1], [8, 5, 2], [9, 6, 3]]
    assert rotateLeft(grid) == [[3, 6, 9], [2, 5, 8], [1, 4, 7]]
    assert transpose(grid) == [[1, 4, 7], [2, 5, 8], [3, 6, 9]]
