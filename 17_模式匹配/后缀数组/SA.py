from typing import List, Sequence, Tuple


def useSA(ords: Sequence[int]) -> Tuple[List[int], List[int], List[int]]:
    """返回 sa, rank, height 数组

    Args:
        ords: 可比较的整数序列

    Returns:
        sa: 每个排名对应的后缀
        rank: 每个后缀对应的排名
        height: 第 i 名的后缀与它前一名的后缀的 `最长公共前缀`的长度
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


# https://leetcode.cn/u/freeyourmind/
def getSA(ords: Sequence[int]) -> List[int]:
    """返回sa数组 即每个排名对应的后缀"""

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
    LMS1 = [i for i in range(n) if isLMS[i]]
    if len(LMS1) > 1:
        SA = inducedSort(LMS1)
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
        LMS1 = [LMS1[i] for i in getSA([SA[i] for i in LMS1])]

    return inducedSort(LMS1)


if __name__ == "__main__":
    print(useSA(list(map(ord, "abcda"))))
    print(getSA([1, 2, 3, 1]))
