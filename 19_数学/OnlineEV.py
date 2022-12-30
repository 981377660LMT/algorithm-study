# 在线求均值方差
# https://github.com/hhy3/cp-library/blob/master/hy/online_ev.hpp
# E : expected value
# V : variance


from decimal import Decimal


class OnlineEV:
    __slots__ = ("_s1", "_s2", "_size")

    def __init__(self) -> None:
        zero = Decimal(0)
        self._s1 = zero
        self._s2 = zero
        self._size = zero

    def add(self, x: int) -> None:
        d = Decimal(x)
        self._s1 += d
        self._s2 += d * d
        self._size += 1

    def remove(self, x: int) -> None:
        d = Decimal(x)
        self._s1 -= d
        self._s2 -= d * d
        self._size -= 1

    def e(self) -> Decimal:
        if self._size == 0:
            raise ValueError("empty")
        return self._s1 / self._size

    def v(self) -> Decimal:
        if self._size == 0:
            raise ValueError("empty")
        return self._s2 / self._size - self._s1 * self._s1 / self._size / self._size


if __name__ == "__main__":
    ev = OnlineEV()
    ev.add(1)
    ev.add(2)
    assert ev.e() == 1.5
    assert ev.v() == 0.25
