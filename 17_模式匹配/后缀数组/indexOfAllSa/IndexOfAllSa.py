# IndexOfAllSa/IndexOfAllMultiString/Lookup/LookupAll


from typing import Any, List, Sequence


def lookupAll(longer: Sequence[Any], longerSa: List[int], shorter: Sequence[Any]) -> List[int]:
    """
    返回s在原串中所有匹配的位置(无序).
    O(len(s)*log(n))+len(result).
    """

    if len(longer) == 0 or len(shorter) == 0 or len(longer) < len(shorter):
        return []

    def check1(mid: int) -> bool:
        return longer[longerSa[mid] :] >= shorter  # type: ignore

    def check2(mid: int) -> bool:
        ptr1 = longerSa[mid + left1]
        n1, n2 = len(longer) - ptr1, len(shorter)
        if n1 < n2:
            return True
        return longer[ptr1 : ptr1 + n2] != shorter

    n = len(longerSa)
    left1, right1 = 0, n - 1
    while left1 <= right1:
        mid = (left1 + right1) // 2
        if check1(mid):
            right1 = mid - 1
        else:
            left1 = mid + 1

    left2, right2 = 0, n - left1 - 1
    while left2 <= right2:
        mid = (left2 + right2) // 2
        if check2(mid):
            right2 = mid - 1
        else:
            left2 = mid + 1

    return [longerSa[k] for k in range(left1, left1 + left2)]


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
            return [sorted(lookupAll(big, sa, small)) for small in smalls]

    longer = [1, 2, 31, 1, 2]
    shorter = [1, 2]
    longerSa = getSA(longer)
    print(lookupAll(longer, longerSa, shorter))  # [3, 0]
