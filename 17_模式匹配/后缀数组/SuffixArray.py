# -支持sa/rank/lcp
# -比较任意两个子串的字典序
# -求出任意两个子串的最长公共前缀(lcp)

# sa : 排第几的后缀是谁.
# rank : 每个后缀排第几.
# lcp : 排名相邻的两个后缀的最长公共前缀.
# lcp[0] = 0
# lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
#
# !"banana" -> sa: [5 3 1 0 4 2], rank: [3 2 5 1 4 0], lcp: [0 1 3 0 0 2]

from typing import Any, List, Sequence, Tuple, Union


class SuffixArray:
    """后缀数组."""

    __slots__ = ("sa", "rank", "height", "_n", "_st")

    def __init__(self, sOrOrds: Union[str, Sequence[int]]) -> None:
        """
        Args:
            sOrOrds (Union[str, Sequence[int]]): 字符串或者字符的`非负数`序列.
            !当ord很大时(>1e7),需要对数组进行离散化,减少内存占用.
        """
        ords = [ord(c) for c in sOrOrds] if isinstance(sOrOrds, str) else sOrOrds
        self._n = len(ords)

        sa = self._saIs(ords)
        rank, height = self._rankLcp(ords, sa)
        self.sa = sa
        self.rank = rank
        self.height = height
        self._st = None

    def lcp(self, a: int, b: int, c: int, d: int) -> int:
        """
        求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
        0 <= a < b <= n, 0 <= c < d <= n.
        """
        if a >= b or c >= d:
            return 0
        res = self._lcp(a, c)
        diff1, diff2 = b - a, d - c
        res = diff1 if res > diff1 else res
        res = diff2 if res > diff2 else res
        return res

    def compareSubstr(self, a: int, b: int, c: int, d: int) -> int:
        """
        比较任意两个子串s[a,b)和s[c,d)的字典序.
        s[a,b) < s[c,d) => -1
        s[a,b) = s[c,d) => 0
        s[a,b) > s[c,d) => 1
        """
        len1, len2 = b - a, d - c
        lcp = self._lcp(a, c)
        if len1 == len2 and lcp >= len1:
            return 0
        if lcp >= len1 or lcp >= len2:
            return -1 if len1 < len2 else 1
        if self.rank[a] < self.rank[c]:
            return -1
        return 1

    def _lcp(self, i: int, j: int) -> int:
        """求两个后缀s[i:]和s[j:]的最长公共前缀(lcp)."""
        if self._st is None:
            self._st = MinSparseTable(self.height)
        if i == j:
            return self._n - i
        r1, r2 = self.rank[i], self.rank[j]
        if r1 > r2:
            r1, r2 = r2, r1
        return self._st.query(r1 + 1, r2 + 1)

    @staticmethod
    def _saIs(ords: Sequence[int]) -> List[int]:
        """SA-IS, linear-time suffix array construction

        Args:
            s (Sequence[int]): Sequence of integers in [0, upper]
            upper (int): Upper bound of the integers in s

        Returns:
            List[int]: Suffix array
        """

        def inducedSort(lms: List[int]) -> List[int]:
            sa = [-1] * n
            sa.append(n)
            endpoint = buckets[1:]
            for j in lms[::-1]:
                endpoint[ords[j]] -= 1
                sa[endpoint[ords[j]]] = j
            startpoint = buckets[:-1]
            for i in range(-1, n):
                j = sa[i] - 1
                if j >= 0 and isL[j]:
                    sa[startpoint[ords[j]]] = j
                    startpoint[ords[j]] += 1
            sa.pop()
            endpoint = buckets[1:]
            for i in range(n - 1, -1, -1):
                j = sa[i] - 1
                if j >= 0 and not isL[j]:
                    endpoint[ords[j]] -= 1
                    sa[endpoint[ords[j]]] = j
            return sa

        n = len(ords)
        buckets = [0] * (max(ords) + 2)
        for a in ords:
            buckets[a + 1] += 1
        for b in range(1, len(buckets)):
            buckets[b] += buckets[b - 1]
        isL = [1] * n
        for i in range(n - 2, -1, -1):
            isL[i] = +(ords[i] > ords[i + 1]) if ords[i] != ords[i + 1] else isL[i + 1]

        isLMS = [(i and isL[i - 1] and not isL[i]) for i in range(n)]
        isLMS.append(True)
        lms1 = [i for i in range(n) if isLMS[i]]
        if len(lms1) > 1:
            sa = inducedSort(lms1)
            lms2 = [i for i in sa if isLMS[i]]
            pre = -1
            j = 0
            for i in lms2:
                i1 = pre
                i2 = i
                while pre >= 0 and ords[i1] == ords[i2]:
                    i1 += 1
                    i2 += 1
                    if isLMS[i1] or isLMS[i2]:
                        j -= isLMS[i1] and isLMS[i2]
                        break
                j += 1
                pre = i
                sa[i] = j
            lms1 = [lms1[i] for i in SuffixArray._saIs([sa[i] for i in lms1])]

        return inducedSort(lms1)

    @staticmethod
    def _rankLcp(ords: Sequence[int], sa: List[int]) -> Tuple[List[int], List[int]]:
        """Rank and LCP array construction

        Args:
            s (Sequence[int]): Sequence of integers in [0, upper]
            sa (List[int]): Suffix array

        Returns:
            Tuple[List[int], List[int]]: Rank array and LCP array

        example:
        ```
        ords = [1, 2, 3, 1, 2, 3]
        sa = _saIs(ords, max(ords))
        rank, lcp = _rankLcp(ords, sa)
        print(rank, lcp)  # [1, 3, 5, 0, 2, 4] [0, 3, 0, 2, 0, 1]
        ```
        """
        n = len(ords)
        rank = [0] * n
        for i, saIndex in enumerate(sa):
            rank[saIndex] = i
        lcp = [0] * n
        h = 0
        for i in range(n):
            if h > 0:
                h -= 1
            if rank[i] == 0:
                continue
            j = sa[rank[i] - 1]
            while j + h < n and i + h < n:
                if ords[j + h] != ords[i + h]:
                    break
                h += 1
            lcp[rank[i]] = h
        return rank, lcp


class MinSparseTable:
    """求区间最小值的ST表"""

    __slots__ = "_n", "_h", "_dp"

    def __init__(self, arr: List[int]):
        n = len(arr)
        h = n.bit_length()
        dp = [[0] * n for _ in range(h)]
        dp[0] = [a for a in arr]
        for k in range(1, h):
            t, p = dp[k], dp[k - 1]
            step = 1 << (k - 1)
            for i in range(n - step * 2 + 1):
                t[i] = p[i] if p[i] < p[i + step] else p[i + step]
        self._n = n
        self._h = h
        self._dp = dp

    def query(self, start: int, end: int) -> int:
        """[start,end)区间的最小值."""
        k = (end - start).bit_length() - 1
        cand1, cand2 = self._dp[k][start], self._dp[k][end - (1 << k)]
        return cand1 if cand1 < cand2 else cand2


class SuffixArray2:
    """用于求解`两个字符串s和t`相关性质的后缀数组."""

    __slots__ = ("_sa", "_offset")

    def __init__(self, sOrOrds1: Union[str, Sequence[int]], sOrOrds2: Union[str, Sequence[int]]):
        ords1 = [ord(c) if isinstance(c, str) else c for c in sOrOrds1]
        ords2 = [ord(c) if isinstance(c, str) else c for c in sOrOrds2]
        ords = ords1 + ords2
        self._sa = SuffixArray(ords)
        self._offset = len(ords1)

    def lcp(self, a: int, b: int, c: int, d: int) -> int:
        """求任意两个子串s[a,b)和t[c,d)的最长公共前缀(lcp)."""
        return self._sa.lcp(a, b, c + self._offset, d + self._offset)

    def compareSubstr(self, a: int, b: int, c: int, d: int) -> int:
        """比较任意两个子串s[a,b)和t[c,d)的字典序.
        s[a,b) < t[c,d) 返回-1.
        s[a,b) = t[c,d) 返回0.
        s[a,b) > t[c,d) 返回1.
        """
        return self._sa.compareSubstr(a, b, c + self._offset, d + self._offset)


def longestCommonSubstring(arr1: Sequence[Any], arr2: Sequence[Any]) -> Tuple[int, int, int, int]:
    """两个序列的最长公共子串.元素的值很大时,需要对元素进行离散化."""
    n1 = len(arr1)
    n2 = len(arr2)
    if not n1 or not n2:
        return 0, 0, 0, 0

    if isinstance(arr1, str):
        arr1 = [ord(c) for c in arr1]
    if isinstance(arr2, str):
        arr2 = [ord(c) for c in arr2]

    dummy = max(max(arr1), max(arr2)) + 1
    sb = list(arr1) + [dummy] + list(arr2)
    S = SuffixArray(sb)
    sa = S.sa
    height = S.height
    maxSame = 0
    start1 = 0
    start2 = 0
    for i in range(1, len(sb)):
        if (sa[i - 1] < n1) == (sa[i] < n1) or height[i] <= maxSame:
            continue
        maxSame = height[i]
        i1 = sa[i - 1]
        i2 = sa[i]
        if i1 > i2:
            i1, i2 = i2, i1
        start1 = i1
        start2 = i2 - n1 - 1

    return start1, start1 + maxSame, start2, start2 + maxSame


if __name__ == "__main__":
    # https://leetcode.cn/problems/sum-of-scores-of-built-strings/
    class Solution:
        def sumScores(self, s: str) -> int:
            sa = SuffixArray(s)
            n = len(s)
            return sum(sa.lcp(0, n, i, n) for i in range(n))

    def lcpNaive(s, a: int, b: int, c: int, d: int) -> int:
        res = 0
        while a < b and c < d and s[a] == s[c]:
            res += 1
            a += 1
            c += 1
        return res

    def compareSubstr(s, a: int, b: int, c: int, d: int) -> int:
        while a < b and c < d and s[a] == s[c]:
            a += 1
            c += 1
        if a == b:  # s[a:b] 到头了
            return 0 if c == d else -1
        if c == d:  # s[c:d] 到头了
            return 1
        return -1 if s[a] < s[c] else 1

    import random

    n = 30
    ords = [random.randint(0, 100) for _ in range(n)]
    sa = SuffixArray(ords)
    for a in range(n):
        for b in range(a, n):
            for c in range(n):
                for d in range(c, n):
                    assert sa.lcp(a, b, c, d) == lcpNaive(ords, a, b, c, d)
                    assert sa.compareSubstr(a, b, c, d) == compareSubstr(ords, a, b, c, d)

    assert (longestCommonSubstring("abcde", "cdeab")) == (2, 5, 0, 3)

    sa2 = SuffixArray2("abcde", "cdeab")
    assert sa2.lcp(2, 5, 0, 3) == 3
    assert sa2.compareSubstr(2, 5, 0, 3) == 0

    print("pass")
