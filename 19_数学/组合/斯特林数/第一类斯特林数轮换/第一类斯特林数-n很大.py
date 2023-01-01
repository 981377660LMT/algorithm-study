# n<=5e5

# Convolution_998244353

MOD = 998244353
IMAG = 911660635
IIMAG = 86583718
rate2 = (
    911660635,
    509520358,
    369330050,
    332049552,
    983190778,
    123842337,
    238493703,
    975955924,
    603855026,
    856644456,
    131300601,
    842657263,
    730768835,
    942482514,
    806263778,
    151565301,
    510815449,
    503497456,
    743006876,
    741047443,
    56250497,
    867605899,
)
irate2 = (
    86583718,
    372528824,
    373294451,
    645684063,
    112220581,
    692852209,
    155456985,
    797128860,
    90816748,
    860285882,
    927414960,
    354738543,
    109331171,
    293255632,
    535113200,
    308540755,
    121186627,
    608385704,
    438932459,
    359477183,
    824071951,
    103369235,
)
rate3 = (
    372528824,
    337190230,
    454590761,
    816400692,
    578227951,
    180142363,
    83780245,
    6597683,
    70046822,
    623238099,
    183021267,
    402682409,
    631680428,
    344509872,
    689220186,
    365017329,
    774342554,
    729444058,
    102986190,
    128751033,
    395565204,
)
irate3 = (
    509520358,
    929031873,
    170256584,
    839780419,
    282974284,
    395914482,
    444904435,
    72135471,
    638914820,
    66769500,
    771127074,
    985925487,
    262319669,
    262341272,
    625870173,
    768022760,
    859816005,
    914661783,
    430819711,
    272774365,
    530924681,
)


def _butterfly_(a):
    n = len(a)
    h = (n - 1).bit_length()
    len_ = 0
    while len_ < h:
        if h - len_ == 1:
            p = 1 << (h - len_ - 1)
            rot = 1
            for s in range(1 << len_):
                offset = s << (h - len_)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p] * rot % MOD
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) % MOD
                if s + 1 != 1 << len_:
                    rot *= rate2[(~s & -~s).bit_length() - 1]
                    rot %= MOD
            len_ += 1
        else:
            p = 1 << (h - len_ - 2)
            rot = 1
            for s in range(1 << len_):
                rot2 = rot * rot % MOD
                rot3 = rot2 * rot % MOD
                offset = s << (h - len_)
                for i in range(p):
                    a0 = a[i + offset]
                    a1 = a[i + offset + p] * rot
                    a2 = a[i + offset + p * 2] * rot2
                    a3 = a[i + offset + p * 3] * rot3
                    a1na3imag = (a1 - a3) % MOD * IMAG
                    a[i + offset] = (a0 + a2 + a1 + a3) % MOD
                    a[i + offset + p] = (a0 + a2 - a1 - a3) % MOD
                    a[i + offset + p * 2] = (a0 - a2 + a1na3imag) % MOD
                    a[i + offset + p * 3] = (a0 - a2 - a1na3imag) % MOD
                if s + 1 != 1 << len_:
                    rot *= rate3[(~s & -~s).bit_length() - 1]
                    rot %= MOD
            len_ += 2


def _butterfly_inv_(a):
    n = len(a)
    h = (n - 1).bit_length()
    len_ = h
    while len_:
        if len_ == 1:
            p = 1 << (h - len_)
            irot = 1
            for s in range(1 << (len_ - 1)):
                offset = s << (h - len_ + 1)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p]
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) * irot % MOD
                if s + 1 != (1 << (len_ - 1)):
                    irot *= irate2[(~s & -~s).bit_length() - 1]
                    irot %= MOD
            len_ -= 1
        else:
            p = 1 << (h - len_)
            irot = 1
            for s in range(1 << (len_ - 2)):
                irot2 = irot * irot % MOD
                irot3 = irot2 * irot % MOD
                offset = s << (h - len_ + 2)
                for i in range(p):
                    a0 = a[i + offset]
                    a1 = a[i + offset + p]
                    a2 = a[i + offset + p * 2]
                    a3 = a[i + offset + p * 3]
                    a2na3iimag = (a2 - a3) * IIMAG % MOD
                    a[i + offset] = (a0 + a1 + a2 + a3) % MOD
                    a[i + offset + p] = (a0 - a1 + a2na3iimag) * irot % MOD
                    a[i + offset + p * 2] = (a0 + a1 - a2 - a3) * irot2 % MOD
                    a[i + offset + p * 3] = (a0 - a1 - a2na3iimag) * irot3 % MOD
                if s + 1 != (1 << (len_ - 2)):
                    irot *= irate3[(~s & -~s).bit_length() - 1]
                    irot %= MOD
            len_ -= 2
    inv = pow(n, MOD - 2, MOD)
    for i in range(n):
        a[i] *= inv
        a[i] %= MOD


def build_exp(n, b):
    exp = [0] * (n + 1)
    exp[0] = 1
    for i in range(n):
        exp[i + 1] = exp[i] * b % MOD
    return exp


def build_factorial(n):
    fct = [0] * (n + 1)
    inv = [0] * (n + 1)
    fct[0] = inv[0] = 1
    for i in range(n):
        fct[i + 1] = fct[i] * (i + 1) % MOD
    inv[n] = pow(fct[n], MOD - 2, MOD)
    for i in range(n)[::-1]:
        inv[i] = inv[i + 1] * (i + 1) % MOD
    return fct, inv


def sqrt_mod(n):
    if n == 0:
        return 0
    if n == 1:
        return 1
    h = (MOD - 1) // 2
    if pow(n, h, MOD) != 1:
        return -1
    q, s = MOD - 1, 0
    while not q & 1:
        q >>= 1
        s += 1
    z = 1
    while pow(z, h, MOD) != MOD - 1:
        z += 1
    m, c, t, r = s, pow(z, q, MOD), pow(n, q, MOD), pow(n, (q + 1) // 2, MOD)
    while t != 1:
        k = 1
        while pow(t, 1 << k, MOD) != 1:
            k += 1
        x = pow(c, pow(2, m - k - 1, MOD - 1), MOD)
        m = k
        c = (x * x) % MOD
        t = (t * c) % MOD
        r = (r * x) % MOD
    if r * r % MOD != n:
        return -1
    return r


class FormalPowerSeries:
    def __init__(self, arr=None):
        if arr is None:
            arr = []
        self.arr = [v % MOD for v in arr]

    def __len__(self):
        return len(self.arr)

    def __getitem__(self, key):
        if isinstance(key, slice):
            return FormalPowerSeries(self.arr[key])
        else:
            assert key >= 0
            if key >= len(self):
                return 0
            else:
                return self.arr[key]

    def __setitem__(self, key, val):
        assert key >= 0
        if key >= len(self):
            self.arr += [0] * (key - len(self) + 1)
        self.arr[key] = val % MOD

    def __str__(self):
        return " ".join(map(str, self.arr))

    def resize(self, sz):
        assert sz >= 0
        if len(self) >= sz:
            return self[:sz]
        else:
            return FormalPowerSeries(self.arr + [0] * (sz - len(self)))

    def shrink(self):
        while self.arr and not self.arr[-1]:
            self.arr.pop()

    def times(self, k):
        if k:
            return FormalPowerSeries([v * k for v in self.arr])
        else:
            return FormalPowerSeries([])

    def __pos__(self):
        return self

    def __neg__(self):
        return self.times(-1)

    def __add__(self, other):
        if other.__class__ == FormalPowerSeries:
            n = len(self)
            m = len(other)
            arr = [self[i] + other[i] for i in range(min(n, m))]
            if n >= m:
                arr += self.arr[m:]
            else:
                arr += other.arr[n:]
            return FormalPowerSeries(arr)
        else:
            return self + FormalPowerSeries([other])

    def __iadd__(self, other):
        if other.__class__ == FormalPowerSeries:
            n = len(self)
            m = len(other)
            for i in range(min(n, m)):
                self.arr[i] += other[i]
                self.arr[i] %= MOD
            if n < m:
                self.arr += other.arr[n:]
        else:
            self.arr[0] += other
            self.arr[0] %= MOD
        return self

    def __radd__(self, other):
        return self + other

    def __sub__(self, other):
        return self + (-other)

    def __isub__(self, other):
        self += -other
        return self

    def __rsub__(self, other):
        return (-self) + other

    def __mul__(self, other):
        if other.__class__ == FormalPowerSeries:
            f = self.arr.copy()
            g = other.arr.copy()
            n = len(f)
            m = len(g)
            if not n or not m:
                return FormalPowerSeries()
            if min(n, m) <= 50:
                if n < m:
                    f, n, g, m = g, m, f, n
                arr = [0] * (n + m - 1)
                for i in range(n):
                    for j in range(m):
                        arr[i + j] += f[i] * g[j]
                        arr[i + j] %= MOD
                return FormalPowerSeries(arr)
            z = 1 << (n + m - 2).bit_length()
            f += [0] * (z - n)
            g += [0] * (z - m)
            _butterfly_(f)
            _butterfly_(g)
            for i in range(z):
                f[i] *= g[i]
                f[i] %= MOD
            _butterfly_inv_(f)
            f = f[: n + m - 1]
            return FormalPowerSeries(f)
        else:
            return self.times(other)

    def __imul__(self, other):
        if other.__class__ == FormalPowerSeries:
            f = self.arr.copy()
            g = other.arr.copy()
            n = len(f)
            m = len(g)
            if not n or not m:
                return FormalPowerSeries()
            if min(n, m) <= 50:
                if n < m:
                    f, n, g, m = g, m, f, n
                arr = [0] * (n + m - 1)
                for i in range(n):
                    for j in range(m):
                        arr[i + j] += f[i] * g[j]
                        arr[i + j] %= MOD
                self.arr = arr
                return self
            z = 1 << (n + m - 2).bit_length()
            f += [0] * (z - n)
            g += [0] * (z - m)
            _butterfly_(f)
            _butterfly_(g)
            for i in range(z):
                f[i] *= g[i]
                f[i] %= MOD
            _butterfly_inv_(f)
            self.arr = f[: n + m - 1]
            return self
        else:
            n = len(self)
            for i in range(n):
                self.arr[i] *= other
                self.arr[i] %= MOD
            return self

    def __rmul__(self, other):
        return self.times(other)

    def __pow__(self, k):
        n = len(self)
        for d, p in enumerate(self.arr):
            if p != 0:
                break
        else:
            return FormalPowerSeries([0] * n)
        res = FormalPowerSeries([0] * n)
        g = ((self[d:].resize(n) / p).log() * k).exp()
        for i in range(max(n - d * k, 0)):
            res[d * k + i] = g[i]
        res *= pow(p, k, MOD)
        return res

    def square(self):
        f = self.arr.copy()
        n = len(f)
        if not n:
            return FormalPowerSeries()
        if n <= 50:
            arr = [0] * (2 * n - 1)
            for i in range(n):
                for j in range(n):
                    arr[i + j] += f[i] * f[j]
                    arr[i + j] %= MOD
            return FormalPowerSeries(arr)
        z = 1 << (2 * n - 2).bit_length()
        f += [0] * (z - n)
        _butterfly_(f)
        for i in range(z):
            f[i] *= f[i]
            f[i] %= MOD
        _butterfly_inv_(f)
        f = f[: 2 * n - 1]
        return FormalPowerSeries(f)

    def __lshift__(self, key):
        assert key >= 0
        return FormalPowerSeries([0] * key + self.arr)

    def __rshift__(self, key):
        assert key >= 0
        return self[key:]

    def __invert__(self):
        assert self[0] != 0
        n = len(self)
        r = pow(self[0], MOD - 2, MOD)
        m = 1
        res = FormalPowerSeries([r])
        while m < n:
            f = [0] * (2 * m)
            g = [0] * (2 * m)
            for i in range(2 * m):
                f[i] = self[i]
            for i in range(m):
                g[i] = res[i]
            _butterfly_(f)
            _butterfly_(g)
            for i in range(2 * m):
                f[i] *= g[i]
                f[i] %= MOD
            _butterfly_inv_(f)
            for i in range(m):
                f[i] = 0
            _butterfly_(f)
            for i in range(2 * m):
                f[i] *= g[i]
                f[i] %= MOD
            _butterfly_inv_(f)
            for i in range(m, 2 * m):
                res[i] -= f[i]
            m <<= 1
        return res.resize(n)

    def __truediv__(self, other):
        if other.__class__ == FormalPowerSeries:
            n = max(len(self), len(other))
            return (self * ~other).resize(n)
        else:
            return self * pow(other, MOD - 2, MOD)

    def __rtruediv__(self, other):
        return other * ~self

    def differentiate(self):
        n = len(self)
        arr = [0] * n
        for i in range(1, n):
            arr[i - 1] = self[i] * i % MOD
        return FormalPowerSeries(arr)

    def integrate(self):
        n = len(self)
        arr = [0] * n
        inv = [1] * n
        for i in range(2, n):
            inv[i] = MOD - inv[MOD % i] * (MOD // i) % MOD
        for i in range(n - 1):
            arr[i + 1] = self[i] * inv[i + 1] % MOD
        return FormalPowerSeries(arr)

    def log(self):
        assert self[0] == 1
        n = len(self)
        return (self.differentiate() / self).integrate()

    def exp(self):
        assert self[0] == 0
        n = len(self)
        res = FormalPowerSeries([1])
        g = FormalPowerSeries([1])
        q = self.differentiate()
        m = 1
        while m < n:
            g = g * 2 - res * g.square().resize(m)
            res = res.resize(2 * m)
            m *= 2
            w = q.resize(m) + (g * (res.differentiate() - (res * q.resize(m)).resize(m))).resize(m)
            res = res + (res * (self.resize(m) - w.integrate())).resize(m)
        return res.resize(n)

    def __floordiv__(self, other):
        if other.__class__ == FormalPowerSeries:
            n = len(self)
            m = len(other)
            if n < m:
                return FormalPowerSeries([])
            l = n - m + 1
            if m <= 100:
                arr = [0] * l
                inv = pow(other[m - 1], MOD - 2, MOD)
                tmp = self[::-1]
                for i in range(l):
                    arr[i] = tmp[i] * inv % MOD
                    for j in range(m):
                        tmp[i + j] -= other[m - j - 1] * arr[i]
                        tmp[i + j] %= MOD
                return FormalPowerSeries(arr[::-1])
            res = (self[~l:][::-1] * ~(other[::-1].resize(l))).resize(l)[::-1]
            return res
        else:
            return self * pow(other, MOD - 2, MOD)

    def __rfloordiv__(self, other):
        return other * ~self

    def __mod__(self, other):
        if other.__class__ == FormalPowerSeries:
            n = len(self)
            m = len(other)
            if n < m:
                return FormalPowerSeries(self.arr)
            res = self[: m - 1] - ((self // other) * other)[: m - 1]
            res.shrink()
            return res
        else:
            return 0

    def divmod(self, other):
        if other.__class__ == FormalPowerSeries:
            div = self // other
            n = len(self)
            m = len(other)
            if n < m:
                mod = FormalPowerSeries(self.arr)
            else:
                mod = self[: m - 1] - ((self // other) * other)[: m - 1]
                mod.shrink()
        else:
            div = self // other
            mod = 0
        return div, mod

    def __matmul__(self, other):
        assert self.__class__ == other.__class__ == FormalPowerSeries
        assert other[0] == 0
        assert len(self) == len(other)
        n = len(self)
        # fをkブロックに分割する。dはブロック内の要素数。k >= dになるように。
        k = int((n - 1) ** 0.5 + 1)
        d = (n + k - 1) // k
        powg = [FormalPowerSeries([1])]
        for i in range(k):
            powg.append((powg[i] * other).resize(n))
        fi = [FormalPowerSeries([0] * n) for _ in range(k)]
        for i in range(k):
            for j in range(d):
                if i * d + j >= n:
                    break
                for t in range(n):
                    if t >= len(powg[j]):
                        break
                    fi[i][t] += powg[j][t] * self[i * d + j]
                    fi[i][t] %= MOD
        res = FormalPowerSeries([0] * n)
        gd = FormalPowerSeries([1])
        for i in range(k):
            fi[i] *= gd
            fi[i] = fi[i].resize(n)
            res += fi[i]
            gd *= powg[d]
            gd = gd.resize(n)
        return res

    def _sqrt_(self):
        assert self[0] != 0
        n = len(self)
        s = sqrt_mod(self[0])
        if s == -1:
            return
        m = 1
        res = FormalPowerSeries([s])
        minv2 = pow(2, MOD - 2, MOD)
        while m < n:
            res = res.resize(2 * m)
            m *= 2
            res = (res + self.resize(m) / res) * minv2
        return res.resize(n)

    def sqrt(self):
        n = len(self)
        for d, p in enumerate(self.arr):
            if p != 0:
                break
        else:
            return FormalPowerSeries([0] * n)
        if d % 2 == 1:
            return -1
        s = sqrt_mod(p)
        if s == -1:
            return -1
        res = FormalPowerSeries([0] * n)
        g = (self[d:].resize(n) / p)._sqrt_()
        if g is None:
            return -1
        g *= s
        for i in range(n - d + 1):
            if d // 2 + i >= n:
                break
            res[d // 2 + i] = g[i]
        return res

    def multipoint_evaluation(self, xs):
        n = len(xs)
        sz = 1 << (n - 1).bit_length()
        g = [FormalPowerSeries([1]) for _ in range(2 * sz)]
        for i in range(n):
            g[i + sz] = FormalPowerSeries([-xs[i], 1])
        for i in range(1, sz)[::-1]:
            g[i] = g[2 * i] * g[2 * i + 1]
        g[1] = self % g[1]
        for i in range(2, 2 * sz):
            g[i] = g[i >> 1] % g[i]
        res = [g[i + sz][0] for i in range(n)]
        return res

    def taylor_shift(self, c):
        n = len(self)
        fct, inv = build_factorial(n)
        d = build_exp(n, c)
        p = FormalPowerSeries([self[i] * fct[i] % MOD for i in range(n)])
        q = FormalPowerSeries([d[n - 1 - i] * inv[n - 1 - i] % MOD for i in range(n)])
        r = p * q
        res = FormalPowerSeries([r[n - 1 + i] * inv[i] % MOD for i in range(n)])
        return res


def polynomial_interpolation(xs, ys):
    assert len(xs) == len(ys)
    n = len(xs)
    sz = 1 << (n - 1).bit_length()
    f = [FormalPowerSeries([1]) for _ in range(2 * sz)]
    for i in range(n):
        f[i + sz] = FormalPowerSeries([-xs[i], 1])
    for i in range(1, sz)[::-1]:
        f[i] = f[2 * i] * f[2 * i + 1]
    g = [FormalPowerSeries([0])] * (2 * sz)
    g[1] = f[1].differentiate() % f[1]
    for i in range(2, n + sz):
        g[i] = g[i >> 1] % f[i]
    for i in range(n):
        g[i + sz] = FormalPowerSeries([ys[i] * pow(g[i + sz][0], MOD - 2, MOD) % MOD])
    for i in range(1, sz)[::-1]:
        g[i] = g[2 * i] * f[2 * i + 1] + g[2 * i + 1] * f[2 * i]
    return g[1][:n]


def berlekamp_massey(arr):
    if arr.__class__ == FormalPowerSeries:
        arr = arr.arr
    n = len(arr)
    b = [1]
    c = [1]
    l, m, p = 0, 0, 1
    for i in range(n):
        m += 1
        d = arr[i]
        for j in range(1, l + 1):
            d += c[j] * arr[i - j]
            d %= MOD
        if d == 0:
            continue
        t = c.copy()
        q = d * pow(p, MOD - 2, MOD) % MOD
        if len(c) < len(b) + m:
            c += [0] * (len(b) + m - len(c))
        for j in range(len(b)):
            c[j + m] -= q * b[j]
            c[j + m] %= MOD
        if 2 * l <= i:
            b = t
            l, m, p = i + 1 - l, 0, d
    return c


def linear_recurrence(arr, coeff, k):
    if arr.__class__ == FormalPowerSeries:
        arr = arr.arr
    d = len(arr)
    f = FormalPowerSeries(arr)
    q = FormalPowerSeries(coeff)
    p = (f * q).resize(d)
    while k:
        r = [-q[i] if i & 1 else q[i] for i in range(len(q))] + [0] * (d + 1 - len(q))
        r = FormalPowerSeries(r)
        p *= r
        q *= r
        p = p[(k & 1) :: 2]
        q = q[::2]
        k >>= 1
    return p[0] % MOD


def stirling_first(n):
    res = FormalPowerSeries([1])
    t = 0
    for i in range(n.bit_length())[::-1]:
        a = n >> i
        res *= res.taylor_shift(-t)
        t <<= 1
        if a & 1:
            res *= FormalPowerSeries([-t, 1])
            t += 1
    return res


n = int(input())

print(stirling_first(n))  # S(n,0) , S(n,1) , S(n,2) , ... , S(n,n)
