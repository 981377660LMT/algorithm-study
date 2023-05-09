# https://github.com/maspypy/library/blob/5ff559bf9799694fc40011465b27c73a54c1563a/mod/barrett.hpp
# https://www.luogu.com.cn/blog/Sweetlemon/barrett-reduction


class Barrett:
    __slots__ = ("_mod", "_imod")

    def __init__(self, mod: int) -> None:
        self._mod = mod
        self._imod = ((1 << 64) - 1) // mod + 1

    def mod(self, num: int) -> int:
        x = (num * self._imod) >> 64
        y = x * self._mod
        return num - y + (self._mod if num < y else 0)

    def floor(self, num: int) -> int:
        x = (num * self._imod) >> 64
        y = x * self._mod
        return x - 1 if num < y else x

    def mul(self, a: int, b: int) -> int:
        return self.mod(a * b)


if __name__ == "__main__":
    from time import time

    B = Barrett(10)
    time1 = time()
    for _ in range(int(1e7)):
        (B.mod(1234567890123456789))
    time2 = time()
    print(time2 - time1)
    time1 = time()
    for _ in range(int(1e7)):
        (1234567890123456789 % 10)
    time2 = time()
    print(time2 - time1)

    # 内置取模还是快
