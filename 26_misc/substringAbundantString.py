# substringAbundantString
# !长度为n的本质不同子串个数最多的01字符串.
#
# Maximal number of distinct nonempty substrings of any binary string of length n.
# https://oeis.org/A094913
# https://qoj.ac/contest/1096/submissions


from typing import List, Tuple


def max2(a: int, b: int) -> int:
    return a if a > b else b


def maxNumDistinctNonemptySubstrings(n: int) -> str:
    n0 = n
    n = 1
    while (1 << n) + (n - 1) < n0:
        n += 1

    def cal() -> str:
        if n == 1:
            return "01"
        if n == 2:
            return "00110"

        def shift(x: str, y: str) -> str:
            m = len(x)
            x += x
            for i in range(m, m + m):
                if x[i - len(y) : i] == y:
                    return x[i - m : i]
            return ""

        def oplus(x: str, y: str) -> str:
            m = len(y).bit_length() - 1
            return x + shift(y, x[len(x) - m :])

        def neg(x: str) -> str:
            return "".join("0" if s == "1" else "1" for s in x)

        def psi(x: str) -> str:
            res = []
            a = 0
            for v in x:
                a ^= int(v)
                res.append(str(a))
            return "".join(res)

        def nxt(x: str) -> str:
            x = psi(x)
            return oplus(x, neg(x))

        def otimes(x: str, y: str) -> str:
            t = ["0"] * (len(x).bit_length() - 1)
            x = shift(x, "".join(t))
            y = shift(y, "".join(t))
            x = x[-len(t) :] + x[: len(x) - len(t)]
            y = y[-len(t) :] + y[: len(y) - len(t)]

            x0, x1, y0, y1 = 0, 0, 0, 1
            for k, v in runLength(x):
                if k == "0":
                    x0 = max2(x0, v)
                if k == "1":
                    x1 = max2(x1, v)
            for k, v in runLength(y):
                if k == "0":
                    y0 = max2(y0, v)
                if k == "1":
                    y1 = max2(y1, v)
            sbX, sbY = [], []
            for k, v in runLength(x):
                if k == "0" and v < x0:
                    sbX.append(k * v)
                if k == "0" and v == x0:
                    sbX.append(k * (v - 1))
                if k == "1" and v < x1:
                    sbX.append(k * v)
                if k == "1" and v == x1:
                    sbX.append(k * (v + 1))
            for k, v in runLength(y):
                if k == "0" and v < y0:
                    sbY.append(k * v)
                if k == "0" and v == y0:
                    sbY.append(k * (v + 1))
                if k == "1" and v < y1:
                    sbY.append(k * v)
                if k == "1" and v == y1:
                    sbY.append(k * (v - 1))
            return "".join(sbX) + "".join(sbY)

        x, y = "0011", "0011"
        for i in range(2, n - 1):
            t = "1" * i
            x = shift(x, t)
            x = nxt(x)
            y = shift(y, t)
            y = neg(nxt(y))
        x = otimes(x, y)
        for i in range(n - 1):
            x += x[i]
        return x

    return cal()[:n0]


def runLength(s: str) -> List[Tuple[str, int]]:
    res = []
    for x in s:
        if not res or res[-1][0] != x:
            res.append([x, 0])
        res[-1][1] += 1
    return res


def countDistinct(s: str) -> int:
    """给定一个字符串 s, 返回 s 的不同子字符串的个数。"""

    def useSA(ords: List[int]) -> Tuple[List[int], List[int], List[int]]:
        """返回 sa, rank, height 数组.ord值很大时,需要先离散化.

        Args:
            ords: 可比较的整数序列

        Returns:
            sa: 每个排名对应的后缀
            rank: 每个后缀对应的排名
            height: 第 i 名的后缀与它前一名的后缀的 `最长公共前缀(LCP)`的长度
        """
        sa = getSA(ords)
        n, k = len(sa), 0
        rank, height = [0] * n, [0] * n
        for i, saIndex in enumerate(sa):
            rank[saIndex] = i

        for i in range(n):
            if k > 0:
                k -= 1
            while (
                i + k < n
                and rank[i] - 1 >= 0
                and sa[rank[i] - 1] + k < n
                and ords[i + k] == ords[sa[rank[i] - 1] + k]
            ):
                k += 1
            height[rank[i]] = k
        return sa, rank, height

    def getSA(ords: List[int]) -> List[int]:
        """
        返回sa数组 即每个排名对应的后缀.
        ord值很大时,需要先离散化.
        """

        def inducedSort(LMS: List[int]) -> List[int]:
            SA = [-1] * (n)
            SA.append(n)
            endpoint = buckets[1:]
            for j in reversed(LMS):
                endpoint[ords[j]] -= 1
                SA[endpoint[ords[j]]] = j
            startpoint = buckets[:-1]
            for i in range(-1, n):
                j = SA[i] - 1
                if j >= 0 and isL[j]:
                    SA[startpoint[ords[j]]] = j
                    startpoint[ords[j]] += 1
            SA.pop()
            endpoint = buckets[1:]
            for i in reversed(range(n)):
                j = SA[i] - 1
                if j >= 0 and not isL[j]:
                    endpoint[ords[j]] -= 1
                    SA[endpoint[ords[j]]] = j
            return SA

        n = len(ords)
        buckets = [0] * (max(ords) + 2)
        for a in ords:
            buckets[a + 1] += 1
        for b in range(1, len(buckets)):
            buckets[b] += buckets[b - 1]
        isL = [1] * n
        for i in reversed(range(n - 1)):
            isL[i] = +(ords[i] > ords[i + 1]) if ords[i] != ords[i + 1] else isL[i + 1]

        isLMS = [+(i and isL[i - 1] and not isL[i]) for i in range(n)]
        isLMS.append(1)
        lms1 = [i for i in range(n) if isLMS[i]]
        if len(lms1) > 1:
            SA = inducedSort(lms1)
            LMS2 = [i for i in SA if isLMS[i]]
            pre = -1
            j = 0
            for i in LMS2:
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
                SA[i] = j
            lms1 = [lms1[i] for i in getSA([SA[i] for i in lms1])]

        return inducedSort(lms1)

    if len(s) == 0:
        return 0
    *_, height = useSA(list(map(ord, s)))
    n = len(s)
    return n * (n + 1) // 2 - sum(height)


if __name__ == "__main__":
    res = [
        0,
        1,
        3,
        5,
        8,
        12,
        16,
        21,
        27,
        34,
        42,
        50,
        59,
        69,
        80,
        92,
        105,
        119,
        134,
        150,
        166,
        183,
        201,
        220,
        240,
        261,
        283,
        306,
        330,
    ]
    for i in range(len(res)):
        assert countDistinct(maxNumDistinctNonemptySubstrings(i)) == res[i]
