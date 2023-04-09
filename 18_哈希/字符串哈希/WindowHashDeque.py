# 滑动窗口哈希值


from collections import deque
from time import time
from typing import Optional, Sequence, Tuple


class WindowHashDeque:
    __slots__ = ("_mod", "_base", "_hash", "_deque", "_inv", "_power")

    def __init__(self, ords: Optional[Sequence[int]] = None, mod=10**9 + 7, base=13331) -> None:
        """Hash of sliding window.

        Args:
            ords (Optional[Sequence[int]], optional): ords of string. Must be non-negative. Defaults to None.
            mod (_type_, optional): should be prime. Defaults to 10**11+7.
            base (int, optional): power of base. Defaults to 1313131.

        Note:

        The hash function is:

        ```
        hash(s[left:right]) = ((s[left]-offset)*base^(right-left-1) +
        (s[left+1]-offset)*base^(right-left-2) + ... + (s[right-1]-offset)) % mod
        ```
        """
        self._mod = mod
        self._base = base
        self._hash = 0
        self._inv = self._modInv(base, mod)
        self._deque = deque()
        self._power = [1]
        for ord in ords or []:
            self.append(ord)

    def query(self) -> int:
        return self._hash

    def append(self, ord: int) -> None:
        self._hash = (self._hash * self._base + ord) % self._mod
        self._deque.append(ord)

    def appendleft(self, ord: int) -> None:
        self._expand(len(self._deque))
        pow_ = self._power[len(self._deque)]
        self._hash = (self._hash + ord * pow_) % self._mod
        self._deque.appendleft(ord)

    def pop(self) -> None:
        popped = self._deque.pop()
        self._hash = (self._hash - popped) * self._inv % self._mod

    def popleft(self) -> None:
        popped = self._deque.popleft()
        self._expand(len(self._deque))
        pow_ = self._power[len(self._deque)]
        self._hash = (self._hash - popped * pow_) % self._mod

    def _expand(self, size: int) -> None:
        if len(self._power) < size + 1:
            preSz = len(self._power)
            power, base, mod = self._power, self._base, self._mod
            power += [0] * (size + 1 - preSz)
            for i in range(preSz - 1, size):
                power[i + 1] = power[i] * base % mod

    def _modInv(self, a: int, mod: int) -> int:
        """Calculate pow(a, -1, mod) using exgcd."""

        def exgcd(a: int, b: int) -> Tuple[int, int, int]:
            """Extended Euclidean algorithm."""
            if b == 0:
                return a, 1, 0
            gcd_, x, y = exgcd(b, a % b)
            return gcd_, y, x - a // b * y

        gcd_, x, _ = exgcd(a, mod)
        assert gcd_ == 1, f"{self._base} and {self._mod} must be coprime"
        return x % mod

    def __len__(self) -> int:
        return len(self._deque)

    def __repr__(self) -> str:
        return f"WindowHash({list(self._deque)})"


if __name__ == "__main__":
    s = "a"
    ords = [ord(c) for c in s]
    windowHash = WindowHashDeque(ords)
    print(windowHash.query(), windowHash)
    windowHash.append(ord("b"))
    print(windowHash.query(), windowHash)
    windowHash.popleft()
    print(windowHash.query(), windowHash)

    time1 = time()
    for i in range(int(1e5)):
        windowHash.append(ord("b"))
        windowHash.popleft()
        windowHash.appendleft(ord("a"))
        windowHash.pop()
        windowHash.append(ord("b"))
    time2 = time()
    print(time2 - time1)
