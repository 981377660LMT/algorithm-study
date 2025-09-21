from random import getrandbits


class FenwickTree:
    """
    Binary indexed tree (Fenwick tree)
    （外部的に）0-indexed, 関数は半開区間
    0からの区間和、1点加算
    """

    def __init__(self, n, init=None):
        self.size = n
        self.longest_interval = 1 << (n.bit_length() - 1)
        if init is None:
            self.data = [0] * n
        else:
            self.data = list(init)
            assert len(self.data) == n
            for i in range(n):
                i_above = i + ((i + 1) & -(i + 1))
                if i_above < n:
                    self.data[i_above] += self.data[i]

    def prefix_sum(self, r):
        """半閉区間 [0,r) 上の和 a[0] + ... + a[r-1] を返す"""
        s = 0
        while r > 0:
            s += self.data[r - 1]
            r -= r & -r
        return s

    def range_sum(self, l, r):
        """半閉区間 [l,r) 上の和 a[l] + ... + a[r-1] を返す"""
        return self.prefix_sum(r) - self.prefix_sum(l)

    def suffix_sum(self, l):  # a_l + ... (端まで)
        """l 以上の添字での和 a[l] + a[l+1] + ... + を返す"""
        return self.prefix_sum(self.size) - self.prefix_sum(l)

    def add(self, i, x):
        """a[i] += x"""
        i += 1
        while i <= self.size:
            self.data[i - 1] += x
            i += i & -i

    def bisect_left(self, w):
        """
        a[0]+ ... +a[idx] が w 以上になる最小の index (存在しない場合 self.size)
        つまり bit.prefix_sum(idx) < x なる最大の idx
        """
        if w <= 0:
            return 0
        x, k = 0, self.longest_interval
        while k:
            if x + k <= self.size and self.data[x + k - 1] < w:
                w -= self.data[x + k - 1]
                x += k
            k >>= 1
        return x


n, q = map(int, input().split())
bit = FenwickTree(n)
for _ in range(q):
    start, end = map(int, input().split())
    v = bit.range_sum(start, end)  # 如果区间内包含了已写区间的端点偶数次，总和会抵消成 0
    if v != 0:
        print("No")
    else:
        print("Yes")
        x = getrandbits(50)
        bit.add(start, x)
        bit.add(end, -x)
