# IndexOfAllSa/IndexOfAllMultiString/Lookup/LookupAll


from bisect import bisect_left
from typing import Any, List, Sequence, Tuple


def lookupAll(
    longer: Sequence[Any], longerSa: List[int], shorter: Sequence[Any]
) -> Tuple[int, int]:
    """
    返回s在原串中所有匹配的位置区间(无序).
    O(len(s)*log(n))+len(result).
    """

    n, m = len(longer), len(shorter)
    if n == 0 or m == 0 or n < m:
        return 0, 0
    start = bisect_left(longerSa, True, key=lambda x: longer[x : x + m] >= shorter)  # type: ignore
    end = bisect_left(longerSa, True, key=lambda x: longer[x : x + m] > shorter)  # type: ignore
    return start, end


def getSA(ords: Sequence[int]) -> List[int]:
    """
    返回sa数组 即每个排名对应的后缀.
    ord值很大时,需要先离散化.
    """
    if not ords:
        return []

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


if __name__ == "__main__":
    # https://leetcode.cn/problems/multi-search-lcci/
    # 面试题 17.17. 多次搜索
    class Solution:
        def multiSearch(self, big: str, smalls: List[str]) -> List[List[int]]:
            sa = getSA([ord(c) for c in big])
            res = []
            for small in smalls:
                start, end = lookupAll(big, sa, small)
                res.append(sorted(sa[start:end]))
            return res

    # P5357 【模板】AC 自动机（二次加强版）
    # https://www.luogu.com.cn/problem/P5357
    # G - Count Substring Query
    # https://atcoder.jp/contests/abc362/tasks/abc362_g
    # 分别求出每个模式串在文本串中出现的次数。

    def abc362_g() -> None:
        import sys

        input = lambda: sys.stdin.readline().rstrip("\r\n")
        s = input()
        n = int(input())
        words = [input() for _ in range(n)]
        sa = getSA([ord(c) for c in s])
        for word in words:
            start, end = lookupAll(s, sa, word)
            print(end - start)

    abc362_g()
