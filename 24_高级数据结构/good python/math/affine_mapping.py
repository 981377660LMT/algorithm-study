from typing import List, Tuple, Union
from math import sin, cos, radians


class AffineMapping:
    # 平面のアフィン変換クラス
    #
    # 使い方:
    # mat = AffineMapping.new()

    @classmethod
    def _matmul3(
        cls, a: List[List[Union[int, float]]], b: List[List[Union[int, float]]]
    ) -> List[List[Union[int, float]]]:
        res: List[List[Union[int, float]]] = [[0, 0, 0], [0, 0, 0], [0, 0, 0]]
        for i in range(3):
            for k in range(3):
                for j in range(3):
                    res[i][j] += b[k][j] * a[i][k]
        return res

    @classmethod
    def new(cls) -> List[List[Union[int, float]]]:
        return [[1, 0, 0], [0, 1, 0], [0, 0, 1]]

    @classmethod
    def shift(
        cls,
        a: List[List[Union[int, float]]],
        shift_x: Union[int, float] = 0,
        shift_y: Union[int, float] = 0,
    ) -> List[List[Union[int, float]]]:
        b = [[1, 0, shift_x], [0, 1, shift_y], [0, 0, 1]]
        return cls._matmul3(b, a)

    @classmethod
    def expand(
        cls,
        a: List[List[Union[int, float]]],
        ratio_x: Union[int, float] = 1,
        ratio_y: Union[int, float] = 1,
    ) -> List[List[Union[int, float]]]:
        b = [[ratio_x, 0, 0], [0, ratio_y, 0], [0, 0, 1]]
        return cls._matmul3(b, a)

    @classmethod
    def rotate(
        cls, a: List[List[Union[int, float]]], theta: Union[int, float] = 0
    ) -> List[List[Union[int, float]]]:
        if theta == 90:
            b = [[0, -1, 0], [1, 0, 0], [0, 0, 1]]
        elif theta == -90:
            b = [[0, 1, 0], [-1, 0, 0], [0, 0, 1]]
        else:
            theta = radians(theta)
            b: List[List[Union[int, float]]] = [
                [cos(theta), -sin(theta), 0],
                [sin(theta), cos(theta), 0],
                [0, 0, 1],
            ]
        return cls._matmul3(b, a)

    @classmethod
    def x_symmetrical_move(
        cls, a: List[List[Union[int, float]]], p: Union[int, float]
    ) -> List[List[Union[int, float]]]:
        b = [[-1, 0, 2 * p], [0, 1, 0], [0, 0, 1]]
        return cls._matmul3(b, a)

    @classmethod
    def y_symmetrical_move(
        cls, a: List[List[Union[int, float]]], p: Union[int, float]
    ) -> List[List[Union[int, float]]]:
        b = [[1, 0, 0], [0, -1, 2 * p], [0, 0, 1]]
        return cls._matmul3(b, a)

    @staticmethod
    def get(a: List[List[Union[int, float]]], x: float, y: float) -> Tuple[float, float]:
        a0, a1, _ = a
        x, y = a0[0] * x + a0[1] * y + a0[2], a1[0] * x + a1[1] * y + a1[2]
        return x, y
