# Run Enumerate
# n<=2e5
# 寻找字符串中的部分循环节(t,l,r)
# 其中t为循环节长度，l为循环节的起始位置，r为循环节的结束位置 (左闭右开)
# !s的第l个字符到第r-1个字符的最小周期为t且r-l>=2*t

from typing import List, Tuple


class SuffixArray:
    def __init__(self, string):
        self.n = len(string)
        self.str = [ord(c) for c in string]
        self.arr = self.build_sa()
        self.lcp, self.rnk = self.build_lcp()
        self.rmq = SparseTable(self.lcp)

    def build_sa(self):
        def sa_naive(s):
            n = len(s)
            sa = list(range(n))
            sa.sort(key=lambda x: s[x:])
            return sa

        def sa_doubling(s):
            n = len(s)
            sa = list(range(n))
            rnk = s
            tmp = [0] * n
            k = 1
            while k < n:
                sa.sort(key=lambda x: (rnk[x], rnk[x + k]) if x + k < n else (rnk[x], -1))
                tmp[sa[0]] = 0
                for i in range(1, n):
                    tmp[sa[i]] = tmp[sa[i - 1]]
                    x = (
                        (rnk[sa[i - 1]], rnk[sa[i - 1] + k])
                        if sa[i - 1] + k < n
                        else (rnk[sa[i - 1]], -1)
                    )
                    y = (rnk[sa[i]], rnk[sa[i] + k]) if sa[i] + k < n else (rnk[sa[i]], -1)
                    if x < y:
                        tmp[sa[i]] += 1
                k *= 2
                tmp, rnk = rnk, tmp
            return sa

        def sa_is(s, upper):  # recursive function
            n = len(s)
            if n == 0:
                return []
            if n == 1:
                return [0]
            if n == 2:
                if s[0] < s[1]:
                    return [0, 1]
                else:
                    return [1, 0]
            if n < 10:
                return sa_naive(s)
            if n < 100:
                return sa_doubling(s)
            ls = [0] * n
            for i in range(n - 1)[::-1]:
                ls[i] = ls[i + 1] if s[i] == s[i + 1] else s[i] < s[i + 1]
            sum_l = [0] * (upper + 1)
            sum_s = [0] * (upper + 1)
            for i in range(n):
                if ls[i]:
                    sum_l[s[i] + 1] += 1
                else:
                    sum_s[s[i]] += 1
            for i in range(upper):
                sum_s[i] += sum_l[i]
                if i < upper:
                    sum_l[i + 1] += sum_s[i]
            lms_map = [-1] * (n + 1)
            m = 0
            for i in range(1, n):
                if not ls[i - 1] and ls[i]:
                    lms_map[i] = m
                    m += 1
            lms = []
            for i in range(1, n):
                if not ls[i - 1] and ls[i]:
                    lms.append(i)
            sa = [-1] * n
            buf = sum_s.copy()
            for d in lms:
                if d == n:
                    continue
                sa[buf[s[d]]] = d
                buf[s[d]] += 1
            buf = sum_l.copy()
            sa[buf[s[n - 1]]] = n - 1
            buf[s[n - 1]] += 1
            for i in range(n):
                v = sa[i]
                if v >= 1 and not ls[v - 1]:
                    sa[buf[s[v - 1]]] = v - 1
                    buf[s[v - 1]] += 1
            buf = sum_l.copy()
            for i in range(n)[::-1]:
                v = sa[i]
                if v >= 1 and ls[v - 1]:
                    buf[s[v - 1] + 1] -= 1
                    sa[buf[s[v - 1] + 1]] = v - 1
            if m == 0:
                return sa
            sorted_lms = []
            for v in sa:
                if lms_map[v] != -1:
                    sorted_lms.append(v)
            rec_s = [0] * m
            rec_upper = 0
            rec_s[lms_map[sorted_lms[0]]] = 0
            for i in range(1, m):
                l = sorted_lms[i - 1]
                r = sorted_lms[i]
                end_l = lms[lms_map[l] + 1] if lms_map[l] + 1 < m else n
                end_r = lms[lms_map[r] + 1] if lms_map[r] + 1 < m else n
                same = True
                if end_l - l != end_r - r:
                    same = False
                else:
                    while l < end_l:
                        if s[l] != s[r]:
                            break
                        l += 1
                        r += 1
                    if l == n or s[l] != s[r]:
                        same = False
                if not same:
                    rec_upper += 1
                rec_s[lms_map[sorted_lms[i]]] = rec_upper
            rec_sa = sa_is(rec_s, rec_upper)  # recursive call
            for i in range(m):
                sorted_lms[i] = lms[rec_sa[i]]
            sa = [-1] * n
            buf = sum_s.copy()
            for d in sorted_lms:
                if d == n:
                    continue
                sa[buf[s[d]]] = d
                buf[s[d]] += 1
            buf = sum_l.copy()
            sa[buf[s[n - 1]]] = n - 1
            buf[s[n - 1]] += 1
            for i in range(n):
                v = sa[i]
                if v >= 1 and not ls[v - 1]:
                    sa[buf[s[v - 1]]] = v - 1
                    buf[s[v - 1]] += 1
            buf = sum_l.copy()
            for i in range(n)[::-1]:
                v = sa[i]
                if v >= 1 and ls[v - 1]:
                    buf[s[v - 1] + 1] -= 1
                    sa[buf[s[v - 1] + 1]] = v - 1
            return sa

        return sa_is(self.str.copy(), 255)

    def build_lcp(self):
        assert self.n >= 1
        rnk = [0] * self.n
        for i in range(self.n):
            rnk[self.arr[i]] = i
        lcp = [0] * (self.n - 1)
        h = 0
        for i in range(self.n):
            if h > 0:
                h -= 1
            if rnk[i] == 0:
                continue
            j = self.arr[rnk[i] - 1]
            while j + h < self.n and i + h < self.n:
                if self.str[j + h] != self.str[i + h]:
                    break
                h += 1
            lcp[rnk[i] - 1] = h
        return lcp, rnk

    def get_lcp(self, l, r):
        if max(l, r) >= self.n:
            return 0
        if l == r:
            return self.n - l
        l, r = self.rnk[l], self.rnk[r]
        if l > r:
            l, r = r, l
        return self.rmq.prod(l, r)

    def debug(self):
        return [self.str[i:] for i in self.arr]


class SparseTable:
    def __init__(self, arr, op=min):
        self.op = op
        self.n = len(arr)
        self.h = self.n.bit_length() - 1
        self.table = [[0] * self.n for _ in range(self.h + 1)]
        self.table[0] = [a for a in arr]
        for k in range(self.h):
            nxt, prv = self.table[k + 1], self.table[k]
            l = 1 << k
            for i in range(self.n - l * 2 + 1):
                nxt[i] = op(prv[i], prv[i + l])

    def prod(self, l, r):
        assert 0 <= l < r <= self.n
        k = (r - l).bit_length() - 1
        return self.op(self.table[k][l], self.table[k][r - (1 << k)])


# https://judge.yosupo.jp/submission/53029
def run_enum(string: str) -> List[Tuple[int, int, int]]:
    n = len(string)
    sa = SuffixArray(string)
    sa_rev = SuffixArray(string[::-1])
    runs = []
    vis = set()
    lst = -1
    for p in range(1, n // 2 + 1):
        for i in range(0, n - p + 1, p):
            l = i - sa_rev.get_lcp(n - i - p, n - i)
            r = i - p + sa.get_lcp(i, i + p)
            if l > r or l == lst:
                continue
            if (l, r + 2 * p) not in vis:
                vis.add((l, r + 2 * p))
                runs.append((p, l, r + 2 * p))
            lst = l
    return runs


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    s = input()
    res = run_enum(s)
    print(len(res))
    for t, l, r in res:
        print(t, l, r)
