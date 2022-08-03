"""求字典序最大的子串(某个后缀)
1 <= s.length <= 4 * 10^5

1.后缀数组
2.类似于求最大表示法
"""


class Solution:
    def lastSubstring2(self, s: str) -> str:
        """后缀数组解法"""
        nums = list(map(ord, s))
        sa = getSA(nums)
        last = sa[-1]
        return s[last:]

    def lastSubstring(self, s: str) -> str:
        """
        类似于最大表示法求字典序最大的子串
        """
        return find(s, isMin=False)


def find(s: str, isMin=True) -> str:
    """返回字典序最小的/最大的子串

    双指针,指针l记录字典序最大子串的首位下标,指针r向后扫描并与指针l进行比较
    注意这里不能循环位移
    """
    if len(s) <= 1:
        return s

    def compare(s1: str, s2: str) -> int:
        if s1 == s2:
            return 0
        if isMin:
            return 1 if s1 > s2 else -1
        else:
            return 1 if s1 < s2 else -1

    n = len(s)
    i1, i2, same = 0, 1, 0

    while i2 + same < n:  # 注意不能循环位移
        diff = compare(s[i1 + same], s[i2 + same])

        if diff == 0:
            same += 1
            continue
        elif diff > 0:
            i1 += same + 1
        else:
            i2 += same + 1

        if i1 == i2:
            i2 += 1

        same = 0

    return s[i1:]


from typing import List, Sequence

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


print(Solution().lastSubstring("abab"))
print(Solution().lastSubstring2("abab"))
# 输出："bab"
# 解释：我们可以找出 7 个子串 ["a", "ab", "aba", "abab", "b", "ba", "bab"]。按字典序排在最后的子串是 "bab"。
