# pop_count_depth / bit_count_depth


class PopCountDepth:
    __slots__ = ("_max", "_data")

    def __init__(self, max_: int):
        self._max = max_
        b = max_.bit_length()
        self._data = [0] * (b + 1)
        for i in range(2, b + 1):
            self._data[i] = self._data[i.bit_count()] + 1

    def get(self, n: int) -> int:
        if n <= 1:
            return 0
        return self._data[n.bit_count()] + 1


if __name__ == "__main__":

    def pop_count_depth(x: int) -> int:
        res = 0
        while x > 1:
            res += 1
            x = x.bit_count()
        return res

    P = PopCountDepth(int(1e18))
    for i in range(10000):
        assert pop_count_depth(i) == P.get(i), f"Error at {i}"
    print("All tests passed!")
