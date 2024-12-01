# 滑动窗口字符串哈希


from collections import deque
from random import randint
from typing import Optional, Tuple


MOD61 = (1 << 61) - 1


class RollingHashDeque:
    @staticmethod
    def _exgcd(a: int, b: int) -> Tuple[int, int, int]:
        if b == 0:
            return a, 1, 0
        gcd_, x, y = RollingHashDeque._exgcd(b, a % b)
        return gcd_, y, x - a // b * y

    @staticmethod
    def _modInv(a: int, mod: int) -> int:
        gcd_, x, _ = RollingHashDeque._exgcd(a, mod)
        if gcd_ != 1:
            raise ValueError(f"no inverse element for {a}")
        return x % mod

    __slots__ = ("_a", "_base", "_ibase", "_powTable", "_s", "_data")

    def __init__(self, base: Optional[int] = None):
        if base is None:
            base = randint(1, MOD61 - 1)
        self._base = base
        self._ibase = RollingHashDeque._modInv(base, MOD61)
        self._powTable = [1, base]
        self._s = deque()
        self._data = deque()
        self._data.append(0)

    def popleft(self) -> int:
        self._data.popleft()
        ch = self._s.popleft()
        self._a = (self._a * self._base - ch) % MOD61
        return ch

    def pop(self) -> int:
        self._data.pop()
        return self._s.pop()

    def appendleft(self, ch: int) -> None:
        self._s.appendleft(ch)
        self._a = (self._a + ch) * self._ibase % MOD61
        self._data.appendleft(-self._a % MOD61)

    def append(self, ch: int) -> None:
        self._s.append(ch)
        self._data.append((self._data[-1] * self._base + ch) % MOD61)

    def get(self, l: int, r: int) -> int:
        if l < 0:
            l = 0
        if r > len(self._data):
            r = len(self._data)
        if l >= r:
            return 0
        return (self._data[r] - self._data[l] * self._pow(r - l)) % MOD61

    def _pow(self, i: int) -> int:
        base = self._base
        while i >= len(self._powTable):
            self._powTable.append(self._powTable[-1] * base % MOD61)
        return self._powTable[i]

    def __len__(self) -> int:
        return len(self._s)


if __name__ == "__main__":
    import unittest

    class TestRollingHashDeque(unittest.TestCase):
        def setUp(self):
            self.deque = RollingHashDeque(base=31)

        def test_append_and_pop(self):
            self.deque.append(1)
            self.deque.append(2)
            self.deque.append(3)
            self.assertEqual(len(self.deque), 3)
            self.assertEqual(self.deque.pop(), 3)
            self.assertEqual(len(self.deque), 2)

        def test_appendleft_and_popleft(self):
            self.deque.appendleft(1)
            self.deque.appendleft(2)
            self.deque.appendleft(3)
            self.assertEqual(len(self.deque), 3)
            self.assertEqual(self.deque.popleft(), 3)
            self.assertEqual(len(self.deque), 2)

        def test_get(self):
            for i in range(5):
                self.deque.append(i)
            expected_0_3 = (0 + 1 * 31 + 2 * 31**2) % MOD61
            expected_2_5 = (2 + 3 * 31 + 4 * 31**2) % MOD61
            self.assertEqual(self.deque.get(0, 3), expected_0_3)
            self.assertEqual(self.deque.get(2, 5), expected_2_5)

        def test_no_inverse(self):
            with self.assertRaises(ValueError):
                RollingHashDeque._modInv(2, 4)  # type: ignore # 2 和 4 不是互素

        def test_len(self):
            self.assertEqual(len(self.deque), 0)
            self.deque.append(10)
            self.assertEqual(len(self.deque), 1)
            self.deque.popleft()
            self.assertEqual(len(self.deque), 0)

        # 可以添加更多测试方法以覆盖其他功能

        unittest.main()
