from typing import List


MOD = 998244353


class Convolution:
    """fft求卷积
    https://atcoder.jp/contests/typical90/submissions/24203987

    时间复杂度 nlogn
    """

    def __init__(self, mod: int):
        self._mod = mod
        self._g = self._primitive_root(mod)
        self._first_butterfly = True
        self._first_butterfly_inv = True
        self._sum_e = [0] * 30
        self._sum_ie = [0] * 30

    def convolution(self, nums1: List[int], nums2: List[int]):
        MOD = self._mod
        n = len(nums1)
        m = len(nums2)
        if (n == 0) | (m == 0):
            return []
        if min(n, m) <= 60:
            if n < m:
                nums1, nums2 = nums2, nums1
                n, m = m, n
            ans = [0] * (n + m - 1)
            for i in range(n):
                for j in range(m):
                    ans[i + j] += nums1[i] * nums2[j]
                    ans[i + j] %= MOD
            return ans

        z = 1 << (n + m - 2).bit_length()
        nums1 += [0] * (z - n)
        nums2 += [0] * (z - m)
        self._butterfly(nums1)
        self._butterfly(nums2)
        for i in range(z):
            nums1[i] *= nums2[i]
            nums1[i] %= MOD
        self._butterfly_inv(nums1)
        nums1 = nums1[: (n + m - 1)]
        iz = pow(z, MOD - 2, MOD)
        for i in range(n + m - 1):
            nums1[i] *= iz
            nums1[i] %= MOD
        return nums1

    # 原始根の取得
    def _primitive_root(self, m: int):
        if m == 2:
            return 1
        if m == 167772161:
            return 3
        if m == 469762049:
            return 3
        if m == 754974721:
            return 11
        if m == 998244353:
            return 3
        divs = [0] * 20
        divs[0] = 2
        cnt = 1
        x = (m - 1) // 2
        while x % 2 == 0:
            x //= 2
        for i in range(3, x + 1, 2):
            if i**2 > x:
                break
            if x % i == 0:
                divs[cnt] = i
                cnt += 1
                while x % i == 0:
                    x //= i
        if x > 1:
            divs[cnt] = x
            cnt += 1
        g = 2
        while True:
            ok = True
            for i in range(cnt):
                if pow(g, (m - 1) // divs[i], m) == 1:
                    ok = False
                    break
            if ok:
                return g
            g += 1

    def _butterfly(self, a: List[int]):
        MOD = self._mod
        n = len(a)
        h = (n - 1).bit_length()
        if self._first_butterfly:
            self._first_butterfly = False
            es = [0] * 30
            ies = [0] * 30
            mod_m = MOD - 1
            cnt2 = (mod_m & -mod_m).bit_length() - 1
            e = pow(self._g, mod_m >> cnt2, MOD)
            ie = pow(e, MOD - 2, MOD)
            for i in range(cnt2 - 2, -1, -1):
                es[i] = e
                ies[i] = ie
                e *= e
                e %= MOD
                ie *= ie
                ie %= MOD
            now = 1
            for i in range(cnt2 - 1):
                self._sum_e[i] = (es[i] * now) % MOD
                now *= ies[i]
                now %= MOD
        for ph in range(1, h + 1):
            w = 1 << (ph - 1)
            p = 1 << (h - ph)
            now = 1
            for s in range(w):
                offset = s << (h - ph + 1)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p] * now
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) % MOD
                now *= self._sum_e[(~s & -~s).bit_length() - 1]
                now %= MOD

    def _butterfly_inv(self, a: List[int]):
        MOD = self._mod
        n = len(a)
        h = (n - 1).bit_length()
        if self._first_butterfly_inv:
            self._first_butterfly_inv = False
            es = [0] * 30
            ies = [0] * 30
            mod_m = MOD - 1
            cnt2 = (mod_m & -mod_m).bit_length() - 1
            e = pow(self._g, mod_m >> cnt2, MOD)
            ie = pow(e, MOD - 2, MOD)
            for i in range(cnt2 - 2, -1, -1):
                es[i] = e
                ies[i] = ie
                e *= e
                e %= MOD
                ie *= ie
                ie %= MOD
            now = 1
            for i in range(cnt2 - 1):
                self._sum_ie[i] = (ies[i] * now) % MOD
                now *= es[i]
                now %= MOD
        for ph in range(h, 0, -1):
            w = 1 << (ph - 1)
            p = 1 << (h - ph)
            inow = 1
            for s in range(w):
                offset = s << (h - ph + 1)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p]
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = ((l - r) * inow) % MOD
                inow *= self._sum_ie[(~s & -~s).bit_length() - 1]
                inow %= MOD


if __name__ == "__main__":
    # https://atcoder.jp/contests/typical90/submissions/26046157
    R, G, B, K = map(int, input().split())
    X, Y, Z = map(int, input().split())
    needR, needG, needB = K - Y, K - Z, K - X  # 不能少于

    factorial = [1, 1]
    inverse = [1, 1]
    invere_base = [0, 1]
    for i in range(2, 2 * 10**5 + 1):
        factorial.append((factorial[-1] * i) % MOD)
        invere_base.append((-invere_base[MOD % i] * (MOD // i)) % MOD)  # 逆元
        inverse.append((inverse[-1] * invere_base[-1]) % MOD)  # 階乗逆元

    def nCr(n, r):
        if r < 0 or r > n:
            return 0
        r = min(r, n - r)
        return factorial[n] * inverse[r] % MOD * inverse[n - r] % MOD

    rc = [nCr(R, i) if needR <= i <= R else 0 for i in range(K + 1)]
    gc = [nCr(G, i) if needG <= i <= G else 0 for i in range(K + 1)]
    bc = [nCr(B, i) if needB <= i <= B else 0 for i in range(K + 1)]

    conv = Convolution(MOD)

    rgc = conv.convolution(rc, gc)
    res = 0
    for i in range(K + 1):
        res += rgc[i] * bc[K - i] % MOD
        res %= MOD
    print(res)
