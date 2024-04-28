# 二维仿射变换.
# api:
#  1. New() -> [3][3]T
#  2. Shift(mapping [3][3]T, shiftX, shiftY T) -> [3][3]T
#  3. Expand(mapping [3][3]T, ratioX, ratioY T) -> [3][3]T
#  4. Rotate90Clockwise(mapping [3][3]T) -> [3][3]T
#  5. Rotate90AntiClockwise(mapping [3][3]T) -> [3][3]T
#  6. RotateClockwise(mapping [3][3]T, degree T) -> [3][3]T
#  7. RotateAntiClockwise(mapping [3][3]T, degree T) -> [3][3]T
#  8. XSymmetricalMove(mapping [3][3]T, x T) -> [3][3]T
#  9. YSymmetricalMove(mapping [3][3]T, y T) -> [3][3]T
#  10. Get(mapping [3][3]T, x, y T) -> (T, T)

# 二维平面仿射变换.

from typing import List, Tuple, Union
from math import sin, cos, radians


class AffineMapping:
    @classmethod
    def new(cls) -> List[List[Union[int, float]]]:
        return [[1, 0, 0], [0, 1, 0], [0, 0, 1]]

    @classmethod
    def shift(
        cls,
        mapping: List[List[Union[int, float]]],
        shiftX: Union[int, float] = 0,
        shiftY: Union[int, float] = 0,
    ) -> List[List[Union[int, float]]]:
        b = [[1, 0, shiftX], [0, 1, shiftY], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def expand(
        cls,
        mapping: List[List[Union[int, float]]],
        ratioX: Union[int, float] = 1,
        ratioY: Union[int, float] = 1,
    ) -> List[List[Union[int, float]]]:
        b = [[ratioX, 0, 0], [0, ratioY, 0], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def rotate90Clockwise(
        cls, mapping: List[List[Union[int, float]]]
    ) -> List[List[Union[int, float]]]:
        b = [[0, 1, 0], [-1, 0, 0], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def rotate90AntiClockwise(
        cls, mapping: List[List[Union[int, float]]]
    ) -> List[List[Union[int, float]]]:
        b = [[0, -1, 0], [1, 0, 0], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def rotateClockwise(
        cls, mapping: List[List[Union[int, float]]], degree: Union[int, float] = 0
    ) -> List[List[Union[int, float]]]:
        cos_, sin_ = cos(degree), sin(degree)
        degree = radians(degree)
        b = [[cos_, sin_, 0], [-sin_, cos_, 0], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def rotateAntiClockwise(
        cls, mapping: List[List[Union[int, float]]], degree: Union[int, float] = 0
    ) -> List[List[Union[int, float]]]:
        cos_, sin_ = cos(degree), sin(degree)
        degree = radians(degree)
        b = [[cos_, -sin_, 0], [sin_, cos_, 0], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def xSymmetricalMove(
        cls, mapping: List[List[Union[int, float]]], x: Union[int, float]
    ) -> List[List[Union[int, float]]]:
        b = [[-1, 0, 2 * x], [0, 1, 0], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @classmethod
    def ySymmetricalMove(
        cls, mapping: List[List[Union[int, float]]], y: Union[int, float]
    ) -> List[List[Union[int, float]]]:
        b = [[1, 0, 0], [0, -1, 2 * y], [0, 0, 1]]
        return cls._matmul3(b, mapping)  # type: ignore

    @staticmethod
    def get(mapping: List[List[Union[int, float]]], x: float, y: float) -> Tuple[float, float]:
        a0, a1, _ = mapping
        x, y = a0[0] * x + a0[1] * y + a0[2], a1[0] * x + a1[1] * y + a1[2]
        return x, y

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


if __name__ == "__main__":
    m0 = AffineMapping.new()
    m1 = AffineMapping.shift(m0, 1, 2)
    m2 = AffineMapping.expand(m1, 2, 3)
    print(AffineMapping.get(m2, 1, 1))
    m3 = AffineMapping.rotate90Clockwise(m0)
    print(AffineMapping.get(m3, 1, 1))
