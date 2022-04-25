from functools import lru_cache
from itertools import product
from typing import List, Set, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)
Shape = Tuple[str]
Matrix = List[List[str]]


def rotate(matrix: Matrix, times=1) -> List[Tuple[str, ...]]:
    """顺时针旋转矩阵90度`times`次
    [['0', '0', '0'], ['1', '1', '0'], ['0', '0', '0']]
    =>
    [('0', '1', '0'), ('0', '1', '0'), ('0', '0', '0')]
    """
    assert times >= 1
    res = [col[::-1] for col in zip(*matrix)]
    for _ in range(times - 1):
        res = [col[::-1] for col in zip(*res)]
    return res


# 记忆化生成形状
@lru_cache(None)
def getAllShape(shape: Shape) -> Set[Shape]:
    """输入：('000', '110', '000')"""
    matrix1 = [list(row) for row in shape]
    matrix2 = [list(row[::-1]) for row in shape]
    res = set()
    res.add(tuple(''.join(row) for row in matrix1))
    res.add(tuple(''.join(row) for row in matrix2))
    for times in (1, 2, 3):
        for matrix in (matrix1, matrix2):
            nextMatrix = rotate(matrix, times)
            res.add(tuple(''.join(row) for row in nextMatrix))
    return res


assert getAllShape(('000', '110', '000')) == {
    ('000', '010', '010'),
    ('010', '010', '000'),
    ('000', '011', '000'),
    ('000', '110', '000'),
}


class Solution:
    def composeCube(self, S: List[List[str]]) -> bool:
        """可左右对称、旋转方向"""

        def check(states: Tuple[Shape]) -> bool:
            """判断合法"""
            ...
            ...

        n = len(S[0])
        shapes = list(map(tuple, S))
        for states in product(*[getAllShape(shape) for shape in shapes]):
            if check(states):
                return True

        return False


print(
    Solution().composeCube(
        S=[
            ["000", "110", "000"],
            ["110", "011", "000"],
            ["110", "011", "110"],
            ["000", "010", "111"],
            ["011", "111", "011"],
            ["011", "010", "000"],
        ]
    )
)
