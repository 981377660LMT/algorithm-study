# 矩阵旋转90


from typing import List, TypeVar


T = TypeVar("T")


def rotate(matrix: List[List[T]], times: int) -> List[List[T]]:
    """顺时针旋转矩阵90度`times`次"""
    if times == 0:
        return [list(row) for row in matrix]
    res = [list(col[::-1]) for col in zip(*matrix)]
    for _ in range(times - 1):
        res = [list(col[::-1]) for col in zip(*res)]
    return res


# 转置口诀:zip星负一
class Solution:
    def rotate(self, matrix: List[List[int]]):
        """
        Do not return anything, modify matrix in-place instead.
        """
        # matrix[:] = zip(*matrix[::-1])
        # zip(*matrix)表示转置矩阵
        return [list(col[::-1]) for col in zip(*matrix)]


a = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
]
print(*a[::-1])
print(Solution().rotate(a))
