# https://judge.yosupo.jp/submission/118808
# Deque Operate All Composite


class SlidingWindowAggregationDeque:
    __slots__ = ("lval", "rval", "lsum", "rsum", "op")

    def __init__(self, op):
        self.lval = []
        self.rval = []
        self.lsum = []
        self.rsum = []
        self.op = op

    def is_empty(self):
        return not (len(self.lsum) or len(self.rsum))

    def appendleft(self, x):
        if not self.lsum:
            self.lval.append(x)
            self.lsum.append(x)
        else:
            self.lval.append(x)
            self.lsum.append(self.op(x, self.lsum[-1]))

    def append(self, x):
        if not self.rsum:
            self.rval.append(x)
            self.rsum.append(x)
        else:
            self.rval.append(x)
            self.rsum.append(self.op(self.rsum[-1], x))

    def popleft(self):
        if not self.lsum:
            rn = len(self.rsum) // 2
            ln = len(self.rsum) - rn
            rv = []
            self.rsum.clear()
            for _ in range(rn):
                rv.append(self.rval.pop())
            for _ in range(ln):
                x = self.rval.pop()
                self.lval.append(x)
                if not self.lsum:
                    self.lsum.append(x)
                else:
                    self.lsum.append(self.op(x, self.lsum[-1]))
            for _ in range(rn):
                x = rv.pop()
                self.rval.append(x)
                if not self.rsum:
                    self.rsum.append(x)
                else:
                    self.rsum.append(self.op(self.rsum[-1], x))
        self.lval.pop()
        self.lsum.pop()

    def pop(self):
        if not self.rsum:
            ln = len(self.lsum) // 2
            rn = len(self.lsum) - ln
            lv = []
            self.lsum.clear()
            for _ in range(ln):
                lv.append(self.lval.pop())
            for _ in range(rn):
                x = self.lval.pop()
                self.rval.append(x)
                if not self.rsum:
                    self.rsum.append(x)
                else:
                    self.rsum.append(self.op(self.rsum[-1], x))
            for _ in range(ln):
                x = lv.pop()
                self.lval.append(x)
                if not self.lsum:
                    self.lsum.append(x)
                else:
                    self.lsum.append(self.op(x, self.lsum[-1]))
        self.rval.pop()
        self.rsum.pop()

    def query(self):
        if not self.rsum:
            return self.lsum[-1]
        if not self.lsum:
            return self.rsum[-1]
        return self.op(self.lsum[-1], self.rsum[-1])


MOD = 998244353
mask = (1 << 32) - 1

# 合并 (ax1+b1) 与 (ax2+b2) 操作
def op(f1: int, f2: int) -> int:
    mul1, add1 = f1 >> 32, f1 & mask
    mul2, add2 = f2 >> 32, f2 & mask
    mul = mul1 * mul2 % MOD
    add = (mul2 * add1 + add2) % MOD
    return (mul << 32) | add


window = SlidingWindowAggregationDeque(op)
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
        if window.is_empty():
            res.append(x)
        else:
            state = window.query()
            a, b = state >> 32, state & mask
            res.append((a * x + b) % MOD)
print(*res, sep="\n")
