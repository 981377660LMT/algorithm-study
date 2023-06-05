# https://nyaannyaan.github.io/library/game/surreal-number.hpp
# SurrealNumber-超现实数


from typing import Tuple


class SurrealNumber:
    """超现实数,保持p/2^q."""

    @staticmethod
    def normalize(s: "SurrealNumber") -> "SurrealNumber":
        if s.p != 0:
            while s.p & 1 == 0 and s.q != 0:
                s.p >>= 1
                s.q -= 1
        else:
            s.q = 0
        return SurrealNumber(s.p, s.q)

    @staticmethod
    def reduce(left: "SurrealNumber", right: "SurrealNumber") -> "SurrealNumber":
        # assert left < right
        root = SurrealNumber(0, 0)
        while left >= root or root >= right:
            lr, rr = root.child()
            root = lr if right <= root else rr
        return root

    __slots__ = ("p", "q")

    def __init__(self, p=0, q=0):
        self.p = p
        self.q = q

    def child(self) -> Tuple["SurrealNumber", "SurrealNumber"]:
        p = self.p
        if p == 0:
            return SurrealNumber(-1, 1), SurrealNumber(1, 1)
        q = self.q
        if q == 0 and p > 0:
            return (SurrealNumber(p * 2 - 1, 1), SurrealNumber(p + 1, 1))
        if q == 0 and p < 0:
            return (SurrealNumber(p - 1, 1), SurrealNumber(p * 2 + 1, 1))
        return (self - SurrealNumber(1, q + 1), self + SurrealNumber(1, q + 1))

    def larger(self) -> "SurrealNumber":
        root = SurrealNumber(0, 0)
        while self >= root:
            root = root.child()[1]
        return root

    def smaller(self) -> "SurrealNumber":
        root = SurrealNumber(0, 0)
        while self <= root:
            root = root.child()[0]
        return root

    def __add__(self, other: "SurrealNumber") -> "SurrealNumber":
        q1, q2 = self.q, other.q
        cq = q1 if q1 > q2 else q2
        cp = (self.p << (cq - q1)) + (other.p << (cq - q2))
        return self.normalize(SurrealNumber(cp, cq))

    def __iadd__(self, other: "SurrealNumber") -> "SurrealNumber":
        return self + other

    def __sub__(self, other: "SurrealNumber") -> "SurrealNumber":
        q1, q2 = self.q, other.q
        cq = q1 if q1 > q2 else q2
        cp = (self.p << (cq - q1)) - (other.p << (cq - q2))
        return self.normalize(SurrealNumber(cp, cq))

    def __isub__(self, other: "SurrealNumber") -> "SurrealNumber":
        return self - other

    def __neg__(self) -> "SurrealNumber":
        return SurrealNumber(-self.p, self.q)

    def __pos__(self) -> "SurrealNumber":
        return SurrealNumber(self.p, self.q)

    def __lt__(self, other: "SurrealNumber") -> bool:
        return (other - self).p > 0

    def __le__(self, other: "SurrealNumber") -> bool:
        return (other - self).p >= 0

    def __gt__(self, other: "SurrealNumber") -> bool:
        return (self - other).p > 0

    def __ge__(self, other: "SurrealNumber") -> bool:
        return (self - other).p >= 0

    def __eq__(self, other: object) -> bool:
        return isinstance(other, SurrealNumber) and (self - other).p == 0

    def __ne__(self, other: object) -> bool:
        return not isinstance(other, SurrealNumber) or (self - other).p != 0

    def __repr__(self) -> str:
        return f"SurrealNumber({self.p}, {self.q})"

    def __hash__(self) -> int:
        return hash((self.p, self.q))


if __name__ == "__main__":
    s1 = SurrealNumber(1, 3)
    print(s1.child())
