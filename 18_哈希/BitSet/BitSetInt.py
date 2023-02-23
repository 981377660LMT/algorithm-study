class BitSet:
    @staticmethod
    def fromlist(lst) -> "BitSet":
        res = BitSet(len(lst))
        for i, v in enumerate(lst):
            if int(v) == 1:
                res.add(i)
        return res

    __slots__ = "_n", "_bits"

    def __init__(self, n: int):
        self._n = n
        self._bits = 0

    def add(self, i: int) -> None:
        self._bits |= 1 << i

    def add_range(self, left: int, right: int) -> None:
        """Add bits in range [left, right]."""
        mask = (1 << (right - left + 1)) - 1
        self._bits |= mask << left

    def discard(self, i: int) -> None:
        self._bits &= ~(1 << i)

    def discard_range(self, left: int, right: int) -> None:
        """Discard bits in range [left, right]."""
        mask = (1 << (right - left + 1)) - 1
        self._bits &= ~(mask << left)

    def flip(self, i: int) -> None:
        self._bits ^= 1 << i

    def flip_range(self, left: int, right: int) -> None:
        """Flip bits in range [left, right]."""
        mask = (1 << (right - left + 1)) - 1
        self._bits ^= mask << left

    def bit_count(self) -> int:
        return self._bits.bit_count()

    def bit_length(self) -> int:
        return self._bits.bit_length()

    def __contains__(self, i: int) -> bool:
        return not not (self._bits & (1 << i))

    def __iter__(self):
        for i in range(self._n):
            yield self._bits >> i & 1

    def __repr__(self) -> str:
        return f"BitSet({list(self)})"

    def __and__(self, other: "BitSet") -> "BitSet":
        return BitSet(self._n).__iand__(other)

    def __iand__(self, other: "BitSet") -> "BitSet":
        self._bits &= other._bits
        return self

    def __or__(self, other: "BitSet") -> "BitSet":
        return BitSet(self._n).__ior__(other)

    def __ior__(self, other: "BitSet") -> "BitSet":
        self._bits |= other._bits
        return self

    def __xor__(self, other: "BitSet") -> "BitSet":
        return BitSet(self._n).__ixor__(other)

    def __ixor__(self, other: "BitSet") -> "BitSet":
        self._bits ^= other._bits
        return self

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, BitSet):
            return False
        return self._n == other._n and self._bits == other._bits

    def __hash__(self) -> int:
        return hash((self._bits, self._n))


if __name__ == "__main__":
    bs = BitSet.fromlist([1, 0, 1, 1, 0, 1, 0, 1])
    print(bs)
