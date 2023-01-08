# https://judge.yosupo.jp/submission/118808
# Deque Operate All Composite


from typing import Callable, Generic, TypeVar

E = TypeVar("E")


class SlidingWindowAggregationDeque(Generic[E]):
    __slots__ = ("_lval", "_rval", "_lsum", "_rsum", "_op", "_e")

    def __init__(self, e: Callable[[], E], op: Callable[[E, E], E]):
        self._e = e
        self._op = op
        self._lval = []
        self._rval = []
        self._lsum = []
        self._rsum = []

    def query(self) -> "E":
        """Aggregate all f(x) from left to right."""
        if not self:
            return self._e()
        if not self._rsum:
            return self._lsum[-1]
        if not self._lsum:
            return self._rsum[-1]
        return self._op(self._lsum[-1], self._rsum[-1])

    def appendleft(self, x: "E") -> None:
        if not self._lsum:
            self._lval.append(x)
            self._lsum.append(x)
        else:
            self._lval.append(x)
            self._lsum.append(self._op(x, self._lsum[-1]))

    def append(self, x: "E") -> None:
        if not self._rsum:
            self._rval.append(x)
            self._rsum.append(x)
        else:
            self._rval.append(x)
            self._rsum.append(self._op(self._rsum[-1], x))

    def popleft(self) -> None:
        # 暴力重构
        if not self._lsum:
            rn = len(self._rsum) // 2
            ln = len(self._rsum) - rn
            rv = []
            self._rsum.clear()
            for _ in range(rn):
                rv.append(self._rval.pop())
            for _ in range(ln):
                x = self._rval.pop()
                self._lval.append(x)
                if not self._lsum:
                    self._lsum.append(x)
                else:
                    self._lsum.append(self._op(x, self._lsum[-1]))
            for _ in range(rn):
                x = rv.pop()
                self._rval.append(x)
                if not self._rsum:
                    self._rsum.append(x)
                else:
                    self._rsum.append(self._op(self._rsum[-1], x))
        self._lval.pop()
        self._lsum.pop()

    def pop(self) -> None:
        # 暴力重构
        if not self._rsum:
            ln = len(self._lsum) // 2
            rn = len(self._lsum) - ln
            lv = []
            self._lsum.clear()
            for _ in range(ln):
                lv.append(self._lval.pop())
            for _ in range(rn):
                x = self._lval.pop()
                self._rval.append(x)
                if not self._rsum:
                    self._rsum.append(x)
                else:
                    self._rsum.append(self._op(self._rsum[-1], x))
            for _ in range(ln):
                x = lv.pop()
                self._lval.append(x)
                if not self._lsum:
                    self._lsum.append(x)
                else:
                    self._lsum.append(self._op(x, self._lsum[-1]))
        self._rval.pop()
        self._rsum.pop()

    def __bool__(self) -> bool:
        return (not not self._lval) or (not not self._rval)

    def __len__(self) -> int:
        return len(self._lval) + len(self._rval)


MOD = 998244353
MASK = (1 << 32) - 1


# 合并 (mul1*x+add1) 与 (mul2*x+add2) 操作
# f2(f1(x))
# f1(x) = mul1*x+add1, f2(x) = mul2*x+add2
# f2(f1(x)) = mul2*(mul1*x+add1)+add2 = mul2*mul1*x+mul2*add1+add2
def op(left: int, right: int) -> int:
    mul1, add1 = left >> 32, left & MASK
    mul2, add2 = right >> 32, right & MASK
    mul = mul1 * mul2 % MOD
    add = (mul2 * add1 + add2) % MOD
    return (mul << 32) | add


window = SlidingWindowAggregationDeque(lambda: 1 << 32, op)
res = []
q = int(input())
for _ in range(q):
    kind, *args = map(int, input().split())
    if kind == 0:  # f の先頭に一次関数 ax+b を追加する
        a, b = args
        window.appendleft((a << 32) | b)
    elif kind == 1:  # f の末尾に一次関数 ax+b を追加する
        a, b = args
        window.append((a << 32) | b)
    elif kind == 2:  # f の先頭の一次関数を削除する
        window.popleft()
    elif kind == 3:  # f の末尾の一次関数を削除する
        window.pop()
    elif kind == 4:  # f の全体を求める
        x = args[0]
        state = window.query()
        a, b = state >> 32, state & MASK
        res.append((a * x + b) % MOD)
print(*res, sep="\n")
