# Frequency Table of Tree Distance
# 树上距离的频数表
# n<=2e5


# https://judge.yosupo.jp/submission/24182
class CentroidDecomposition:
    def __init__(self, n):
        self.n = n
        self.tree = [[] for _ in range(n)]
        self.par = [None] * n
        self.dep = [None] * n
        self.size = [1] * n
        self.used = [0] * n
        self.root = None

    def add_edge(self, u, v, idx):
        self.tree[u - idx].append(v - idx)
        self.tree[v - idx].append(u - idx)

    def set_root(self, r, update=True):
        self.root = r
        self.par[r] = -1
        self.dep[r] = 0
        self.size[r] = 1
        self.ord = [r]
        stack = [r]
        while stack:
            v = stack.pop()
            for adj in self.tree[v]:
                if self.par[v] == adj or self.used[adj]:
                    continue
                self.size[adj] = 1
                self.par[adj] = v
                self.dep[adj] = self.dep[v] + 1
                self.ord.append(adj)
                stack.append(adj)
        if update:
            for v in self.ord[::-1]:
                for adj in self.tree[v]:
                    if self.par[v] == adj or self.used[adj]:
                        continue
                    self.size[v] += self.size[adj]

    def centroid(self, r):
        v = r
        while True:
            for adj in self.tree[v]:
                if self.par[v] == adj or self.used[adj]:
                    continue
                if self.size[adj] > self.size[r] // 2:
                    v = adj
                    break
            else:
                return v

    def centroid_decomposition(self, f, g):
        stack = [0]
        ft = True
        while stack:
            v = stack.pop()
            self.set_root(v, True)
            if not ft:
                g(self, v)
            ft = False
            if self.size[v] == 1:
                self.used[v] = True
                continue
            # ---
            if self.size[v] == 2:
                ch = self.tree[v][0]
                self.used[v] = True
                self.used[ch] = True
                distCounter[1] += 2
                continue
            # ---
            c = self.centroid(v)
            self.used[c] = True
            self.set_root(c, False)
            f(self, v)
            for adj in self.tree[c]:
                if self.par[c] == adj or self.used[adj]:
                    continue
                stack.append(adj)


# Convolution_998244353

MOD1 = 167772161
MOD2 = 469762049

sum_e1 = [
    65249968,
    137365239,
    35921276,
    103665800,
    89728614,
    164955302,
    108901219,
    163950188,
    113252399,
    166581688,
    59783366,
    95476790,
    130818126,
    39440948,
    65800545,
    14559656,
    3285286,
    36462062,
    164082627,
    9320421,
    66343657,
    69024390,
    38289678,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
]
sum_ie1 = [
    102522193,
    71493608,
    26998229,
    133555027,
    128975965,
    16363816,
    145463520,
    130828795,
    26375299,
    18078794,
    87407453,
    28151929,
    49401241,
    112914531,
    118959329,
    68815302,
    71865958,
    21459372,
    44393528,
    43709352,
    30681399,
    153195333,
    141748999,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
]
sum_e2 = [
    450151958,
    26623616,
    25192837,
    305390008,
    399060560,
    78724413,
    312251397,
    151088193,
    437503217,
    339869829,
    197503427,
    460844482,
    64795813,
    392699793,
    323591778,
    435162849,
    324666788,
    397071166,
    191521520,
    39442863,
    102932772,
    52822010,
    231589706,
    155147527,
    0,
    0,
    0,
    0,
    0,
    0,
]
sum_ie2 = [
    19610091,
    129701348,
    104677229,
    445839763,
    375500824,
    451642859,
    145445927,
    77724141,
    367250623,
    54456563,
    257713867,
    444918711,
    335270416,
    371371281,
    307213086,
    452878044,
    243328637,
    152011944,
    315423951,
    456185089,
    218081060,
    136058803,
    203260256,
    412215962,
    0,
    0,
    0,
    0,
    0,
    0,
]


def butterfly1(arr):
    n = len(arr)
    h = (n - 1).bit_length()
    for ph in range(1, h + 1):
        w = 1 << (ph - 1)
        p = 1 << (h - ph)
        now = 1
        for s in range(w):
            offset = s << (h - ph + 1)
            for i in range(p):
                l = arr[i + offset]
                r = arr[i + offset + p] * now
                arr[i + offset] = (l + r) % MOD1
                arr[i + offset + p] = (l - r) % MOD1
            now *= sum_e1[(~s & -~s).bit_length() - 1]
            now %= MOD1


def butterfly_inv1(arr):
    n = len(arr)
    h = (n - 1).bit_length()
    for ph in range(1, h + 1)[::-1]:
        w = 1 << (ph - 1)
        p = 1 << (h - ph)
        inow = 1
        for s in range(w):
            offset = s << (h - ph + 1)
            for i in range(p):
                l = arr[i + offset]
                r = arr[i + offset + p]
                arr[i + offset] = (l + r) % MOD1
                arr[i + offset + p] = (MOD1 + l - r) * inow % MOD1
            inow *= sum_ie1[(~s & -~s).bit_length() - 1]
            inow %= MOD1


def convolution1(a, b):
    n = len(a)
    m = len(b)
    if not n or not m:
        return []
    if min(n, m) <= 100:
        if n < m:
            n, m = m, n
            a, b = b, a
        res = [0] * (n + m - 1)
        for i in range(n):
            for j in range(m):
                res[i + j] += a[i] * b[j]
                res[i + j] %= MOD1
        return res
    z = 1 << (n + m - 2).bit_length()
    a += [0] * (z - n)
    b += [0] * (z - m)
    butterfly1(a)
    butterfly1(b)
    for i in range(z):
        a[i] *= b[i]
        a[i] %= MOD1
    butterfly_inv1(a)
    a = a[: n + m - 1]
    iz = pow(z, MOD1 - 2, MOD1)
    for i in range(n + m - 1):
        a[i] *= iz
        a[i] %= MOD1
    return a


def autocorrelation1(a):
    n = len(a)
    if not n:
        return []
    if n <= 100:
        res = [0] * (2 * n - 1)
        for i in range(n):
            for j in range(n):
                res[i + j] += a[i] * a[j]
                res[i + j] %= MOD1
        return res
    z = 1 << (2 * n - 2).bit_length()
    a += [0] * (z - n)
    butterfly1(a)
    for i in range(z):
        a[i] *= a[i]
        a[i] %= MOD1
    butterfly_inv1(a)
    a = a[: 2 * n - 1]
    iz = pow(z, MOD1 - 2, MOD1)
    for i in range(2 * n - 1):
        a[i] *= iz
        a[i] %= MOD1
    return a


def butterfly2(arr):
    n = len(arr)
    h = (n - 1).bit_length()
    for ph in range(1, h + 1):
        w = 1 << (ph - 1)
        p = 1 << (h - ph)
        now = 1
        for s in range(w):
            offset = s << (h - ph + 1)
            for i in range(p):
                l = arr[i + offset]
                r = arr[i + offset + p] * now
                arr[i + offset] = (l + r) % MOD2
                arr[i + offset + p] = (l - r) % MOD2
            now *= sum_e2[(~s & -~s).bit_length() - 1]
            now %= MOD2


def butterfly_inv2(arr):
    n = len(arr)
    h = (n - 1).bit_length()
    for ph in range(1, h + 1)[::-1]:
        w = 1 << (ph - 1)
        p = 1 << (h - ph)
        inow = 1
        for s in range(w):
            offset = s << (h - ph + 1)
            for i in range(p):
                l = arr[i + offset]
                r = arr[i + offset + p]
                arr[i + offset] = (l + r) % MOD2
                arr[i + offset + p] = (MOD2 + l - r) * inow % MOD2
            inow *= sum_ie2[(~s & -~s).bit_length() - 1]
            inow %= MOD2


def convolution2(a, b):
    n = len(a)
    m = len(b)
    if not n or not m:
        return []
    if min(n, m) <= 100:
        if n < m:
            n, m = m, n
            a, b = b, a
        res = [0] * (n + m - 1)
        for i in range(n):
            for j in range(m):
                res[i + j] += a[i] * b[j]
                res[i + j] %= MOD2
        return res
    z = 1 << (n + m - 2).bit_length()
    a += [0] * (z - n)
    b += [0] * (z - m)
    butterfly2(a)
    butterfly2(b)
    for i in range(z):
        a[i] *= b[i]
        a[i] %= MOD2
    butterfly_inv2(a)
    a = a[: n + m - 1]
    iz = pow(z, MOD2 - 2, MOD2)
    for i in range(n + m - 1):
        a[i] *= iz
        a[i] %= MOD2
    return a


def autocorrelation2(a):
    n = len(a)
    if not n:
        return []
    if n <= 100:
        res = [0] * (2 * n - 1)
        for i in range(n):
            for j in range(n):
                res[i + j] += a[i] * a[j]
                res[i + j] %= MOD2
        return res
    z = 1 << (2 * n - 2).bit_length()
    a += [0] * (z - n)
    butterfly2(a)
    for i in range(z):
        a[i] *= a[i]
        a[i] %= MOD2
    butterfly_inv2(a)
    a = a[: 2 * n - 1]
    iz = pow(z, MOD2 - 2, MOD2)
    for i in range(2 * n - 1):
        a[i] *= iz
        a[i] %= MOD2
    return a


def inv_gcd(a, b):
    a %= b
    if a == 0:
        return b, 0
    s = b
    t = a
    m0 = 0
    m1 = 1
    while t:
        u = s // t
        s -= t * u
        m0 -= m1 * u
        s, t = t, s
        m0, m1 = m1, m0
    if m0 < 0:
        m0 += b // s
    return s, m0


def crt(r, m):
    assert len(r) == len(m)
    n = len(r)
    r0 = 0
    m0 = 1
    for i in range(n):
        assert 1 <= m[i]
        r1 = r[i] % m[i]
        m1 = m[i]
        if m0 < m1:
            r0, r1 = r1, r0
            m0, m1 = m1, m0
        if m0 % m1 == 0:
            if r0 % m1 != r1:
                return 0, 0
            continue
        g, im = inv_gcd(m0, m1)
        u1 = m1 // g
        if (r1 - r0) % g:
            return 0, 0
        x = (r1 - r0) // g * im % u1
        r0 += x * m0
        m0 *= u1
        if r0 < 0:
            r0 += m0
    return r0, m0


def autocorrelation(a):
    n = len(a)
    a1 = autocorrelation1(a.copy())
    a2 = autocorrelation2(a.copy())
    res = [0] * (2 * n - 1)
    for i in range(2 * n - 1):
        res[i] = crt([a1[i], a2[i]], [MOD1, MOD2])[0]
    return res


def f(tree, v):
    deps = []
    maxdep = 0
    for v in tree.ord:
        deps.append(tree.dep[v])
        maxdep = max(maxdep, tree.dep[v])
    cnt = [0] * (maxdep + 1)
    for d in deps:
        cnt[d] += 1
    a = autocorrelation(cnt)[: tree.n]
    for i, v in enumerate(a):
        distCounter[i] += v


def g(tree, v):
    if tree.size[v] == 1:
        distCounter[2] -= 1
        return
    if tree.size[v] == 2:
        distCounter[2] -= 1
        distCounter[3] -= 2
        distCounter[4] -= 1
        return
    deps = []
    maxdep = 0
    for v in tree.ord:
        deps.append(tree.dep[v])
        maxdep = max(maxdep, tree.dep[v])
    cnt = [0] * (maxdep + 1)
    for d in deps:
        cnt[d] += 1
    a = autocorrelation(cnt)[: tree.n - 2]
    for i, v in enumerate(a):
        distCounter[i + 2] -= v


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    cd = CentroidDecomposition(n)
    distCounter = [0] * n  # dist(u,v)=i 的(u,v)对数 (u<v)
    for _ in range(n - 1):
        u, v = map(int, input().split())
        cd.add_edge(u, v, 0)
    cd.centroid_decomposition(f, g)
    distCounter = [distCounter[i] // 2 for i in range(1, n)]
    print(*distCounter)
